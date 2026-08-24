import React, { useEffect, useState } from "react";
import {
  Card,
  FormControl,
  FormControlLabel,
  Grid,
  Icon,
  MenuItem,
  Radio,
  RadioGroup,
  Select,
  TextField,
  Tooltip,
} from "@mui/material";
import MDBox from "components/MDBox";
import MDButton from "components/MDButton";
import MDInput from "components/MDInput";
import MDTypography from "components/MDTypography";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import { useDispatch, useSelector } from "react-redux";
import DataTable from "examples/Tables/DataTable";
import {
  wrapCell,
  formatDateTime,
  radioLabelSx,
} from "components/CommonComponent/CommonFunction";
import { DEFAULT_RECORDS_PER_PAGE } from "components/common/constant";
import {
  getUsersList,
  updateUserDetail,
  addUser,
  getUserDetailById,
  deleteUser,
} from "actions/users";
import { userList } from "constants/users";
import CenterPopup from "components/CommonComponent/CenterPopup";
import RequiredLabel from "components/CommonComponent/RequiredLabel";
import MDBadge from "components/MDBadge";
import { REGEX } from "components/common/constant";
import { useNotify } from "components/CommonComponent/NotificationProvider";
import { showAlert } from "components/CommonComponent/ShowAlert";
import { selectSx } from "components/CommonComponent/CommonFunction";
import { ArrowDropDown } from "@mui/icons-material";
import { STATUS } from "components/common/constant";

const initialFields = {
  id: "",
  email: {
    value: "",
    original: "",
    error: "",
  },
  status: {
    value: "ACTIVE",
    original: "ACTIVE",
    error: "",
  },
};
function UsersListing() {
  const dispatch = useDispatch();
  const userListData = useSelector((state) => state.users.userList);
  const { notify } = useNotify();
  const [listState, setListState] = useState({
    isLoading: true,
    error: "",
    noRecords: false,
    showLoaderOnClick: false,
  });
  const [filterState, setFilterState] = useState({
    searchString: "",
    status: null,
    filterParams: {
      filters: [],
      searchString: "",
    },
  });
  const [formData, setFormData] = useState(initialFields);
  const [panelMode, setPanelMode] = useState(null); //add /edit
  const [pageNum, setPageNum] = useState(1);

  useEffect(() => {
    getUsersList(
      dispatch,
      pageNum,
      filterState.filterParams.filters,
      filterState.filterParams.searchString,
    );
    return () => {
      dispatch({ type: userList, payload: {} });
    };
  }, [
    dispatch,
    pageNum,
    filterState.filterParams.filters,
    filterState.filterParams.searchString,
  ]);

  useEffect(() => {
    if (userListData?.code === 200) {
      setListState({
        isLoading: false,
        error: "",
        noRecords: userListData?.data?.user?.length === 0,
        showLoaderOnClick: false,
      });
    } else if (userListData?.code || userListData?.message) {
      setListState({
        isLoading: false,
        error:
          userListData?.message ||
          "Failed to load users list! Please try again.",
        noRecords: false,
        showLoaderOnClick: false,
      });
    }
  }, [userListData]);

  //table data columns
  const formattedTableData =
    userListData?.data?.user?.map((item, index) => ({
      ...item,
      id: item?.ID,
      srNo:
        (pageNum - 1) *
          (userListData?.data?.noOfRecordsPerPage || DEFAULT_RECORDS_PER_PAGE) +
        index +
        1,
      email: item?.email || "-",
      status: item?.status || "ACTIVE",
      lastLoginAt: item?.lastLoginAt || null,
    })) || [];

  const handleKeyDown = (e) => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleSearch();
    }
  };

  const renderStatusBadge = (status) => {
    if (!status) return null;
    return (
      <MDBadge
        badgeContent={status}
        color={status === "ACTIVE" ? "success" : "secondary"}
        variant="gradient"
        size="sm"
        circular
      />
    );
  };

  const columns = [
    {
      Header: "Sr No.",
      accessor: "srNo",
      align: "left",
      width: "85px",
    },
    {
      Header: "Email",
      accessor: "email",
      align: "left",
      width: "200px",
      Cell: (cell) => wrapCell(cell.value || "-", "200px"),
    },
    {
      Header: "Status",
      accessor: "status",
      align: "left",
      width: "100px",
      Cell: (cell) => renderStatusBadge(cell.value),
    },
    {
      Header: "Last Login At",
      accessor: "lastLoginAt",
      align: "left",
      width: "120px",
      Cell: (cell) => (cell.value ? formatDateTime(cell.value) : "-"),
    },
    {
      Header: "Action",
      accessor: "actions",
      align: "center",
      width: "120px",
      disableSortBy: true,
      Cell: (cell) => {
        const rowData = cell.row.original;

        return (
          <MDBox
            display="flex"
            alignItems="center"
            justifyContent="center"
            gap={1}
          >
            <Tooltip title="Edit" arrow>
              <Icon
                sx={({ palette }) => ({
                  cursor: "pointer",
                  color: palette.secondary.main,
                })}
                onClick={(e) => {
                  e.stopPropagation();
                  openEditModal(rowData);
                }}
              >
                <i className="fas fa-edit fa-sm" />
              </Icon>
            </Tooltip>

            <MDTypography variant="caption" mx={0.5}>
              |
            </MDTypography>

            <Tooltip title="Delete" arrow>
              <Icon
                color="error"
                sx={{
                  cursor: "pointer",
                  opacity: 0.5,
                  "&:hover": { opacity: 1 },
                }}
                onClick={() => handleDelete(rowData)}
              >
                <i className="fas fa-trash" />
              </Icon>
            </Tooltip>
          </MDBox>
        );
      },
    },
  ];

  //close the center popup modal
  const closeModal = () => {
    setPanelMode(null);
    setFormData({
      ...initialFields,
    });
  };

  //on change function
  const handleChange = (field, value) => {
    setFormData((prev) => ({
      ...prev,
      [field]: {
        ...prev[field],
        value,
        error: "",
      },
    }));
  };
  //on blur to get if any validation error exits
  const handleOnBlur = (name, value) => {
    const trimmedValue = typeof value === "string" ? value.trim() : value;
    const error = getValidationError(name, trimmedValue);
    setFormData((prev) => ({
      ...prev,
      [name]: {
        ...prev[name],
        value: trimmedValue,
        error,
      },
    }));
  };

  const getValidationError = (name, value) => {
    const trimmedValue = typeof value === "string" ? value.trim() : value;

    switch (name) {
      case "email":
        if (!trimmedValue) {
          return "Email is required.";
        }
        if (!REGEX.EMAIL_REGEX.test(trimmedValue)) {
          return "Please enter a valid email address.";
        }
        return "";

      default:
        return "";
    }
  };

  const handleSubmit = async () => {
    if (!validateForm()) return;

    try {
      setListState((prev) => ({
        ...prev,
        showLoaderOnClick: true,
      }));

      const payload = {};

      Object.keys(formData).forEach((key) => {
        if (key === "id") return;

        const field = formData[key];

        if (panelMode !== "edit" || field.value !== field.original) {
          payload[key] =
            typeof field.value === "string" ? field.value.trim() : field.value;
        }
      });

      // For edit, don't call API if nothing changed
      if (panelMode === "edit" && Object.keys(payload).length === 0) {
        notify("No changes detected", "info");
        return;
      }
      const res =
        panelMode === "edit"
          ? await updateUserDetail(formData.id, payload)
          : await addUser(payload);

      if (res?.code === 200) {
        notify(
          res?.message ||
            (panelMode === "edit"
              ? "User updated successfully"
              : "User added successfully"),
          "success",
        );

        closeModal();
        // Refresh users list
        getUsersList(
          dispatch,
          pageNum,
          filterState.filterParams.filters,
          filterState.filterParams.searchString,
        );
      } else {
        const message = Array.isArray(res?.message)
          ? res.message
              .map((err) => (err?.Msg ? err.Msg : JSON.stringify(err)))
              .join(", ")
          : res?.message || "Failed to save user";

        notify(message, "error");
      }
    } catch (error) {
      notify(error?.message || "Something went wrong", "error");
    } finally {
      setListState((prev) => ({
        ...prev,
        showLoaderOnClick: false,
      }));
    }
  };

  const validateForm = () => {
    const newData = { ...formData };

    Object.keys(newData).forEach((key) => {
      if (key === "id") return;
      newData[key] = {
        ...newData[key],
        error: getValidationError(key, newData[key].value),
      };
    });
    setFormData(newData);
    return Object.keys(newData)
      .filter((key) => key !== "id")
      .every((key) => newData[key].error === "");
  };

  const handleSearch = () => {
    const newfilterParams = buildFilterParams();

    if (
      JSON.stringify(filterState.filterParams) !==
      JSON.stringify(newfilterParams)
    ) {
      setListState((prev) => ({
        ...prev,
        isLoading: true,
        showLoaderOnClick: true,
      }));
      setFilterState((prev) => ({
        ...prev,
        filterParams: newfilterParams,
      }));

      if (pageNum !== 1) setPageNum(1);
    }
  };

  //filters functions
  const buildFilterParams = () => {
    const filterParams = [];

    if (filterState.status) {
      filterParams.push({
        field: "Status",
        condition: "eq",
        filterValues: [filterState.status],
      });
    }

    return {
      filters: filterParams.length ? JSON.stringify(filterParams) : [],
      searchString: filterState.searchString.trim(),
    };
  };

  const handleClear = () => {
    setFilterState({
      searchString: "",
      status: null,
      filterParams: {
        filters: [],
        searchString: "",
      },
    });
    setListState((prev) => ({
      ...prev,
      isLoading: true,
    }));
    if (pageNum !== 1) {
      setPageNum(1);
    }
  };

  const openAddModal = () => {
    setFormData({
      email: {
        value: "",
        original: "",
        error: "",
      },
      status: {
        value: "ACTIVE",
        original: "ACTIVE",
        error: "",
      },
    });

    setPanelMode("add");
  };

  const openEditModal = async (rowData) => {
    try {
      setPanelMode("edit");
      setListState((prev) => ({
        ...prev,
        isLoading: true,
      }));

      const res = await getUserDetailById(rowData.id);
      if (res?.code === 200) {
        const data = res?.data;

        setFormData({
          id: res.data.ID,
          email: {
            value: data?.email || "",
            original: data?.email || "",
            error: "",
          },
          status: {
            value: data?.status || "ACTIVE",
            original: data?.status || "ACTIVE",
            error: "",
          },
        });
      } else {
        const message = Array.isArray(res?.message)
          ? res.message.map((err) => err.Msg || JSON.stringify(err)).join(", ")
          : res?.message || "Failed to fetch user details";

        notify(message, "error");
        setPanelMode(null);
      }
    } catch (error) {
      notify(error?.message || "Something went wrong", "error");
      setPanelMode(null);
    } finally {
      setListState((prev) => ({
        ...prev,
        isLoading: false,
      }));
    }
  };

  //on delete
  const handleDelete = (row) => {
    showAlert({
      title: "Are you sure you want to delete this user?",
      message: `User "${row.email}" will be deleted permanently.`,
      type: "warning",
      showCancel: true,
      confirmText: "Yes, delete",
      cancelText: "Cancel",

      onConfirm: async () => {
        try {
          setListState((prev) => ({
            ...prev,
            isLoading: true,
          }));

          const res = await deleteUser(row.id);

          if (res?.code === 200) {
            notify(res?.message || "User deleted successfully", "success");
            // refresh list
            getUsersList(dispatch, pageNum);
          } else {
            notify(res?.message || "Failed to delete user", "error");
          }
        } catch (error) {
          notify(error?.message || "Something went wrong", "error");
        } finally {
          setListState((prev) => ({
            ...prev,
            isLoading: false,
          }));
        }
      },

      onCancel: () => {},
    });
  };

  return (
    <>
      <DashboardLayout>
        <MDBox width="100">
          <Grid container spacing={6}>
            <Grid item xs={12}>
              <Card sx={{ width: "100%" }}>
                <MDBox
                  px={3}
                  pt={3}
                  pb={2}
                  display="flex"
                  justifyContent="space-between"
                  alignItems="center"
                >
                  <MDTypography variant="h5" fontWeight="bold">
                    Users Lists
                  </MDTypography>

                  <MDButton
                    type="button"
                    variant="contained"
                    color="deepBlue"
                    sx={{
                      py: 1.2,
                      px: 3,
                      fontWeight: "bold",
                      fontSize: "0.9rem",
                      textTransform: "none",
                    }}
                    onClick={openAddModal}
                  >
                    Add User
                  </MDButton>
                </MDBox>

                {/* Filters */}
                <MDBox px={3} pb={2} onKeyDown={handleKeyDown}>
                  <Grid container spacing={2} alignItems="center">
                    <Grid item xs={12} sm={6} md={4} lg={2}>
                      <TextField
                        fullWidth
                        placeholder="Search by name "
                        size="small"
                        value={filterState.searchString}
                        onChange={(e) =>
                          setFilterState({
                            ...filterState,
                            searchString: e.target.value,
                          })
                        }
                      />
                    </Grid>

                    <Grid item xs={12} sm={6} md={4} lg={2}>
                      <FormControl fullWidth size="small" sx={selectSx("100%")}>
                        <Select
                          value={filterState.status || ""}
                          displayEmpty
                          onChange={(e) =>
                            setFilterState({
                              ...filterState,
                              status: e.target.value,
                            })
                          }
                          IconComponent={ArrowDropDown}
                          onKeyDown={(e) => {
                            if (e.key === "Enter") {
                              e.preventDefault();
                              handleSearch();
                            }
                          }}
                          renderValue={(selected) => {
                            if (!selected)
                              return (
                                <span style={{ color: "rgba(0, 0, 0, 0.38)" }}>
                                   Select Status
                                </span>
                              );

                            return renderStatusBadge(selected);
                          }}
                        >
                          {STATUS.map((status) => (
                            <MenuItem key={status.label} value={status.label}>
                              {renderStatusBadge(status.label)}
                            </MenuItem>
                          ))}
                        </Select>
                      </FormControl>
                    </Grid>

                    <Grid item xs={12} sm={6} md={4} lg={3} xl={2}>
                      <MDBox
                        display="flex"
                        gap={1}
                        flexWrap="nowrap"
                        width="100%"
                      >
                        <MDButton
                          variant="contained"
                          color="deepBlue"
                          sx={{
                            flex: 1,
                            minWidth: "90px",
                            maxWidth: "160px",
                            whiteSpace: "nowrap",
                          }}
                          onClick={handleSearch}
                          disabled={listState.showLoaderOnClick}
                        >
                          {listState.showLoaderOnClick ? (
                            <i className="fa fa-spinner fa-spin fa-lg text-white" />
                          ) : (
                            <>
                              <Icon sx={{ mr: 1 }}>search</Icon>
                              Search
                            </>
                          )}
                        </MDButton>

                        {(filterState.filterParams.filters?.length > 0 ||
                          filterState.filterParams.searchString) && (
                          <MDButton
                            variant="contained"
                            color="deepBlue"
                            sx={{
                              flex: 1,
                              minWidth: "90px",
                              maxWidth: "160px",
                              whiteSpace: "nowrap",
                            }}
                            onClick={handleClear}
                          >
                            <Icon sx={{ mr: 1 }}>close</Icon>
                            Clear
                          </MDButton>
                        )}
                      </MDBox>
                    </Grid>
                  </Grid>
                </MDBox>

                <MDBox
                  px={3}
                  py={1}
                  display="flex"
                  justifyContent="space-between"
                  alignItems="center"
                >
                  <MDTypography variant="body2" color="deepBlue" mb={1}>
                    Total Records: {userListData?.data?.totalRecords || 0}
                  </MDTypography>
                </MDBox>

                <MDBox sx={{ width: "100%", overflowX: "auto" }}>
                  {listState.isLoading ? (
                    <MDBox
                      display="flex"
                      justifyContent="center"
                      alignItems="center"
                      height="150px"
                    >
                      <i className="fa fa-spinner fa-spin fa-xl loader" />
                    </MDBox>
                  ) : listState.error ? (
                    <MDBox
                      display="flex"
                      justifyContent="center"
                      alignItems="center"
                      height="150px"
                    >
                      <MDTypography variant="body2" color="error">
                        {listState.error}
                      </MDTypography>
                    </MDBox>
                  ) : listState.noRecords ? (
                    <MDBox
                      display="flex"
                      justifyContent="center"
                      alignItems="center"
                      height="150px"
                    >
                      <MDTypography variant="body2" color="text">
                        No records found.
                      </MDTypography>
                    </MDBox>
                  ) : (
                    <DataTable
                      columns={columns}
                      data={formattedTableData || []}
                      isSorted={true}
                      pageNum={pageNum}
                      setPageNum={setPageNum}
                      totalPages={Math.ceil(
                        (userListData?.data?.totalRecords || 0) /
                          (userListData?.data?.noOfRecordsPerPage ||
                            DEFAULT_RECORDS_PER_PAGE),
                      )}
                      //   onRowClick={handleRowClick}
                    />
                  )}
                </MDBox>
              </Card>
            </Grid>
          </Grid>

          <CenterPopup isOpen={panelMode !== null} onClose={closeModal}>
            <MDBox sx={{ p: 2 }}>
              <MDTypography variant="h6" mb={2}>
                {panelMode === "edit" ? "Edit Email List" : "Add Email List"}
              </MDTypography>

              {listState.isLoading ? (
                <MDBox
                  display="flex"
                  justifyContent="center"
                  alignItems="center"
                  height="150px"
                >
                  <i className="fa fa-spinner fa-spin fa-xl loader" />
                </MDBox>
              ) : (
                <Grid container spacing={2}>
                  <Grid item xs={12}>
                    <RequiredLabel required>Email</RequiredLabel>
                    <MDInput
                      variant="outlined"
                      type="text"
                      name="email"
                      placeholder="Enter email"
                      value={formData.email.value}
                      onChange={(e) =>
                        handleChange(e.target.name, e.target.value)
                      }
                      onBlur={(e) =>
                        handleOnBlur(e.target.name, e.target.value)
                      }
                      error={!!formData.email.error}
                      fullWidth
                      inputProps={{ maxLength: 100 }}
                    />
                    {formData.email.error && (
                      <MDTypography variant="caption" color="error">
                        {formData.email.error}
                      </MDTypography>
                    )}
                  </Grid>

                  <Grid item xs={12}>
                    <RequiredLabel>Status</RequiredLabel>
                    <RadioGroup
                      row
                      name="status"
                      value={formData.status.value}
                      onChange={(e) =>
                        handleChange(e.target.name, e.target.value)
                      }
                    >
                      <FormControlLabel
                        value="ACTIVE"
                        control={<Radio size="small" />}
                        label="Active"
                        sx={radioLabelSx}
                      />
                      <FormControlLabel
                        value="INACTIVE"
                        control={<Radio size="small" />}
                        label="Inactive"
                        sx={radioLabelSx}
                      />
                    </RadioGroup>
                  </Grid>

                  {/* Save Button */}
                  <Grid item xs={12}>
                    <MDBox display="flex" justifyContent="flex-end">
                      <MDButton
                        variant="gradient"
                        color="light"
                        onClick={closeModal}
                        sx={{ mr: 2 }}
                      >
                        Cancel
                      </MDButton>

                      <MDButton
                        type="submit"
                        variant="contained"
                        color="deepBlue"
                        onClick={handleSubmit}
                        disabled={listState.showLoaderOnClick}
                      >
                        {listState.showLoaderOnClick ? (
                          <i className="fa fa-spinner fa-spin" />
                        ) : (
                          "Proceed"
                        )}
                      </MDButton>
                    </MDBox>
                  </Grid>
                </Grid>
              )}
            </MDBox>
          </CenterPopup>
        </MDBox>
      </DashboardLayout>
    </>
  );
}

export default UsersListing;
