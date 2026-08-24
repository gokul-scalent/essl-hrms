import {
  Card,
  FormControl,
  FormControlLabel,
  Grid,
  Icon,
  Radio,
  TextField,
  Select,
  RadioGroup,
  Tooltip,
  MenuItem,
} from "@mui/material";
import {
  getEmailLists,
  deleteEmailListById,
  addEmailList,
  getEmailListById,
  updateEmailListById,
} from "actions/email";
import {
  REGEX,
  PRIORITY,
  DEFAULT_RECORDS_PER_PAGE,
} from "components/common/constant";
import CenterPopup from "components/CommonComponent/CenterPopup";
import {
  radioLabelSx,
  formatDateTime,
  wrapCell,
  selectSx,
} from "components/CommonComponent/CommonFunction";
import { ArrowDropDown } from "@mui/icons-material";
import { useNotify } from "components/CommonComponent/NotificationProvider";
import RequiredLabel from "components/CommonComponent/RequiredLabel";
import { showAlert } from "components/CommonComponent/ShowAlert";
import MDBox from "components/MDBox";
import MDButton from "components/MDButton";
import MDInput from "components/MDInput";
import MDTypography from "components/MDTypography";
import { emailList } from "constants/email";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DataTable from "examples/Tables/DataTable";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import { urlEncodedId } from "components/CommonComponent/encodeDecodeFunction";

const initialFields = {
  id: "",
  name: {
    value: "",
    original: "",
    error: "",
  },
  priority: {
    value: "p1", // default value
    original: "",
    error: "",
  },
};

export const renderPriorityBadge = (value) => {
  const item = PRIORITY.find((p) => p.value === value);
  if (!item) return value;
  return (
    <span
      style={{
        background: item.color,
        color: "#fff",
        padding: "0.25rem 0.75rem",
        borderRadius: "0.75rem",
        fontSize: "0.75rem",
        fontWeight: 600,
      }}
    >
      {item.label}
    </span>
  );
};

function EmailListing() {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const emailListData = useSelector((state) => state.email.emailList);
  const { notify } = useNotify();
  const [pageNum, setPageNum] = useState(1);
  const [panelMode, setPanelMode] = useState(null); //add /edit
  const [formData, setFormData] = useState(initialFields);
  const [listState, setListState] = useState({
    isLoading: true,
    error: "",
    noRecords: false,
    showLoaderOnClick: false,
  });
  const [filterState, setFilterState] = useState({
    searchString: "",
    priority: null,
    filterParams: {
      filters: [],
      searchString: "",
    },
  });

  useEffect(() => {
    getEmailLists(
      dispatch,
      pageNum,
      filterState.filterParams.filters,
      filterState.filterParams.searchString,
    );
    return () => {
      dispatch({ type: emailList, payload: {} });
    };
  }, [
    dispatch,
    pageNum,
    filterState.filterParams.filters,
    filterState.filterParams.searchString,
  ]);

  useEffect(() => {
    if (emailListData?.code === 200) {
      setListState({
        isLoading: false,
        error: "",
        noRecords: emailListData?.data?.emailList?.length === 0,
        showLoaderOnClick: false,
      });
    } else if (emailListData?.code || emailListData?.message) {
      setListState({
        isLoading: false,
        error:
          emailListData?.message ||
          "Failed to load email list! Please try again.",
        noRecords: false,
        showLoaderOnClick: false,
      });
    }
  }, [emailListData]);

  //table data columns
  const formattedTableData =
    emailListData?.data?.emailList?.map((item, index) => ({
      ...item,
      id: item?.ID,
      srNo:
        (pageNum - 1) *
          (emailListData?.data?.noOfRecordsPerPage ||
            DEFAULT_RECORDS_PER_PAGE) +
        index +
        1,
      name: item.name,
      priority: item.priority,
      totalRecords: item.totalRecords,
      processedRecords: item.processedRecords,
      safeCount: item.safeCount,
      riskyCount: item.riskyCount,
      invalidCount: item.invalIDCount,
      unknownCount: item.unknownCount,
      pendingCount: item.pendingCount,
      createdAt: item.createdAt,
    })) || [];

  const columns = [
    {
      Header: "Sr No.",
      accessor: "srNo",
      align: "left",
      width: "85px",
    },
    {
      Header: "Name",
      accessor: "name",
      align: "left",
      width: "170px",
      Cell: (cell) => wrapCell(cell.value || "-", "170px"),
    },
    {
      Header: "Priority",
      accessor: "priority",
      align: "left",
      width: "95px",
      Cell: (cell) => renderPriorityBadge(cell.value),
    },
    {
      Header: "Total Records",
      accessor: "totalRecords",
      align: "left",
      width: "130px",
      Cell: (cell) => wrapCell(cell.value || "-", "130px"),
    },
    {
      Header: "Processed Records",
      accessor: "processedRecords",
      align: "left",
      width: "150px",
      Cell: (cell) => wrapCell(cell.value || "-", "150px"),
    },
    {
      Header: "Remaining Records",
      accessor: "pendingCount",
      align: "left",
      width: "150px",
      Cell: (cell) => wrapCell(cell.value || "-", "150px"),
    },
    {
      Header: "Safe Records",
      accessor: "safeCount",
      align: "left",
      width: "130px",
      Cell: (cell) => wrapCell(cell.value || "-", "130px"),
    },
    {
      Header: "Created At",
      accessor: "createdAt",
      align: "left",
      width: "120px",
      Cell: (cell) => formatDateTime(cell.value),
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

            <Tooltip title="View" arrow>
              <Icon
                sx={({ palette }) => ({
                  cursor: "pointer",
                  color: palette.secondary.main,
                })}
                onClick={() =>
                  navigate(`/admin/lead-list/${rowData.id}`, {
                    state: {
                      emailListName: rowData.name,
                    },
                  })
                }
              >
                <i className="fas fa-eye fa-sm" />
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

  //on row click
  const handleRowClick = (row) => {
    const encodedId = urlEncodedId(row.ID)
    navigate(`/admin/lead-list/${encodedId}`, {
      state: {
        lead: row,
        emailListName: row.name,
      },
    });
  };

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
      [field]: { value, error: "" },
    }));
  };

  const getValidationError = (name, value) => {
    const trimmedValue = typeof value === "string" ? value.trim() : value;

    switch (name) {
      case "name":
        return !trimmedValue ? "Please enter name" : "";

      case "priority":
        return !trimmedValue ? "Please select priority" : "";

      default:
        return "";
    }
  };

  const validateForm = () => {
    let newData = { ...formData };
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

  //on submit
  const handleSubmit = async () => {
    if (!validateForm()) return;

    try {
      setListState((prev) => ({
        ...prev,
        showLoaderOnClick: true,
      }));

      let payload = {};
      Object.keys(formData).forEach((key) => {
        if (key === "id") return;

        const field = formData[key];

        if (panelMode !== "edit" || field.value !== field.original) {
          payload[key] =
            typeof field.value === "string" ? field.value.trim() : field.value;
        }
      });

      if (panelMode === "edit" && !Object.keys(payload).length) {
        notify("No changes detected", "info");
        return;
      }
      const res =
        panelMode === "edit"
          ? await updateEmailListById(formData.id, payload)
          : await addEmailList(payload);

      if (res?.code === 200) {
        notify(
          res?.message ||
            (panelMode === "edit"
              ? "Email list updated successfully"
              : "Email list added successfully"),
          "success",
        );

        closeModal();
        getEmailLists(dispatch, pageNum);
      } else {
        const message = Array.isArray(res?.message)
          ? res.message
              .map((err) => (err.Msg ? err.Msg : JSON.stringify(err)))
              .join(", ")
          : res?.message || "Failed to save email list";

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

  const openAddModal = (e) => {
    e.stopPropagation();
    setPanelMode("add");
    setFormData(initialFields);
  };

  //edit
  const openEditModal = async (rowData) => {
    try {
      setPanelMode("edit");

      setListState((prev) => ({
        ...prev,
        isLoading: true,
      }));

      const res = await getEmailListById(rowData.id);
      if (res?.code === 200) {
        const data = res?.data;

        setFormData({
          id: res.data.ID,
          name: {
            value: data?.name || "",
            original: data?.name || "",
            error: "",
          },
          priority: {
            value: data?.priority || "p1",
            original: data?.priority || "p1",
            error: "",
          },
        });
      } else {
        const message = Array.isArray(res?.message)
          ? res.message.map((err) => err.Msg || JSON.stringify(err)).join(", ")
          : res?.message || "Failed to fetch email list details";

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
      title: "Are you sure you want to delete this email list?",
      message: `Email list "${row.name}" will be deleted permanently.`,
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

          const res = await deleteEmailListById(row.id);

          if (res?.code === 200) {
            notify(
              res?.message || "Email list deleted successfully",
              "success",
            );

            // refresh list
            getEmailLists(dispatch, pageNum);
          } else {
            notify(res?.message || "Failed to delete email list", "error");
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

  //filters functions
  const buildFilterParams = () => {
    const filterParams = [];

    if (filterState.priority) {
      filterParams.push({
        field: "Priority",
        condition: "eq",
        filterValues: [filterState.priority],
      });
    }

    return {
      filters: filterParams.length ? JSON.stringify(filterParams) : [],
      searchString: filterState.searchString.trim(),
    };
  };

  //search function
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

  //clear function to reset the filters and search
  const handleClear = () => {
    setFilterState({
      searchString: "",
      priority: null,
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

  const handleKeyDown = (e) => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleSearch();
    }
  };

  return (
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
                  Email Lists
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
                  Add Email List
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
                        value={filterState.priority || ""}
                        displayEmpty
                        onChange={(e) =>
                          setFilterState({
                            ...filterState,
                            priority: e.target.value,
                          })
                        }
                        IconComponent={ArrowDropDown}
                        renderValue={(selected) => {
                          if (!selected)
                            return (
                              <span style={{ color: "rgba(0, 0, 0, 0.38)" }}>
                                 Select Priority
                              </span>
                            );

                          return renderPriorityBadge(selected);
                        }}
                      >
                        {PRIORITY.map((p) => (
                          <MenuItem key={p.value} value={p.value}>
                            {renderPriorityBadge(p.value)}
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
                  Total Records: {emailListData?.data?.totalRecords || 0}
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
                      (emailListData?.data?.totalRecords || 0) /
                        (emailListData?.data?.noOfRecordsPerPage ||
                          DEFAULT_RECORDS_PER_PAGE),
                    )}
                    onRowClick={handleRowClick}
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
                  <RequiredLabel required>Name</RequiredLabel>
                  <MDInput
                    variant="outlined"
                    type="text"
                    name="name"
                    placeholder="Enter name"
                    value={formData.name.value}
                    onChange={(e) =>
                      handleChange(e.target.name, e.target.value)
                    }
                    onBlur={(e) => handleOnBlur(e.target.name, e.target.value)}
                    error={!!formData.name.error}
                    fullWidth
                    inputProps={{ maxLength: 50 }}
                  />
                  {formData.name.error && (
                    <MDTypography variant="caption" color="error">
                      {formData.name.error}
                    </MDTypography>
                  )}
                </Grid>

                <Grid item xs={12}>
                  <RequiredLabel required>Priority</RequiredLabel>

                  <RadioGroup
                    row
                    name="priority"
                    value={formData.priority.value}
                    onChange={(e) => handleChange("priority", e.target.value)}
                    onBlur={(e) => handleOnBlur("priority", e.target.value)}
                  >
                    {PRIORITY.map((item) => (
                      <FormControlLabel
                        key={item.value}
                        value={item.value}
                        control={<Radio size="small" />}
                        label={renderPriorityBadge(item.value)}
                        sx={radioLabelSx}
                      />
                    ))}
                  </RadioGroup>

                  {formData.priority.error && (
                    <MDTypography variant="caption" color="error">
                      {formData.priority.error}
                    </MDTypography>
                  )}
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
  );
}

export default EmailListing;
