import {
  Card,
  FormControl,
  Grid,
  Icon,
  TextField,
  Select,
  MenuItem,
  Button,
  Menu,
} from "@mui/material";
import colors from "assets/theme/base/colors";
import { useNotify } from "components/CommonComponent/NotificationProvider";
import MDBox from "components/MDBox";
import MDTypography from "components/MDTypography";
import { leadList } from "constants/leads";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DataTable from "examples/Tables/DataTable";
import React, { useEffect, useRef, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  DEFAULT_RECORDS_PER_PAGE,
  LEADS_STATUS,
  STATUS_BADGE_COLORS,
  VERIFICATION_STATUS,
  DOWNLOAD_OPTIONS,
} from "components/common/constant";
import {
  formatDateTime,
  wrapCell,
  leadListUiStyles,
  selectSx,
} from "components/CommonComponent/CommonFunction";
import MDButton from "components/MDButton";
import { useLocation, useNavigate, useParams } from "react-router-dom";
import {
  getLeadsList,
  DownloadSafeCsv,
  deleteLeadByID,
  ReverifyLeadById,
} from "actions/leads";
import AddLeadModal from "./AddLeadModal";
import { ArrowDropDown } from "@mui/icons-material";
import Tooltip from "@mui/material/Tooltip";
import { showAlert } from "components/CommonComponent/ShowAlert";
import LeadDetailModal from "./LeadDetailModal";
import { urlDecodedId } from "components/CommonComponent/encodeDecodeFunction";

export const renderStatusBadge = (value) => {
  if (!value) return "-";

  const badgeColor =
    STATUS_BADGE_COLORS[value.toUpperCase()] || colors.badgeColors.secondary;

  return (
    <span
      style={{
        backgroundColor: badgeColor.background,
        color: badgeColor.text,
        padding: "0.3rem 0.75rem",
        borderRadius: "0.75rem",
        fontSize: "0.75rem",
        fontWeight: 600,
        display: "inline-block",
        minWidth: "90px",
        textAlign: "center",
        textTransform: "capitalize",
      }}
    >
      {value}
    </span>
  );
};

function LeadListing() {
  const { id: emailListID } = useParams();
  const decodedId = urlDecodedId(emailListID);
  const location = useLocation();
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const leadListData = useSelector((state) => state.leads.leadList);
  const { notify } = useNotify();
  const [pageNum, setPageNum] = useState(1);
  const [panelMode, setPanelMode] = useState(null); //add /edit
  const [listState, setListState] = useState({
    isLoading: true,
    error: "",
    noRecords: false,
    showLoaderOnClick: false,
  });
  const [verificationState, setVerificationState] = useState({
    isChecking: false,
    pendingCount: 0,
  });
  const [filterState, setFilterState] = useState({
    searchString: "",
    status: null,
    verificationStatus: null,
    filterParams: {
      filters: [],
      searchString: "",
    },
  });
  const [selectedLeadId, setSelectedLeadId] = useState(null);
  const [anchorEl, setAnchorEl] = useState(null);
  const emailListName = location.state?.emailListName;

  const handleOpenDownload = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleCloseDownload = () => {
    setAnchorEl(null);
  };

  //to get the leads list data
  useEffect(() => {
    getLeadsList(
      dispatch,
      decodedId,
      pageNum,
      filterState.filterParams.filters,
      filterState.filterParams.searchString,
    );
  }, [
    dispatch,
    decodedId,
    pageNum,
    filterState.filterParams.filters,
    filterState.filterParams.searchString,
  ]);

  useEffect(() => {
    return () => {
      dispatch({ type: leadList, payload: {} });
    };
  }, [dispatch]);

  //to check the status or message code of get list data
  useEffect(() => {
    if (leadListData?.code !== 200) {
      return;
    }

    const leads = leadListData?.data?.lead || [];
    const pendingLeads = leads.filter(
      (item) =>
        item.verificationStatus === "PENDING" ||
        item.verificationStatus === "VERIFYING" ||
        (item.verificationStatus === "FAILED" &&
          item.finalStatus === "TIMEOUT" &&
          Number(item.verificationAttempt || 0) < 3),
    );

    setVerificationState({
      isChecking: pendingLeads.length > 0,
      pendingCount: pendingLeads.length,
    });

    setListState({
      isLoading: false,
      error: "",
      noRecords: leads.length === 0,
      showLoaderOnClick: false,
    });
  }, [leadListData]);

  useEffect(() => {
    if (!verificationState.isChecking) {
      return;
    }

    const interval = setInterval(() => {
      getLeadsList(
        dispatch,
        decodedId,
        pageNum,
        filterState.filterParams.filters,
        filterState.filterParams.searchString,
      );
    }, 60000); // 60 seconds

    return () => {
      clearInterval(interval);
    };
  }, [
    verificationState.isChecking,
    dispatch,
    decodedId,
    pageNum,
    filterState.filterParams.filters,
    filterState.filterParams.searchString,
  ]);

  //format the data coming from api for column mapping
  const formattedTableData =
    leadListData?.data?.lead?.map((item, index) => ({
      ...item,
      id: item.ID,
      srNo:
        (pageNum - 1) *
          (leadListData?.data?.noOfRecordsPerPage || DEFAULT_RECORDS_PER_PAGE) +
        index +
        1,

      email: item.email,
      emailProvider: item.emailProvIDer,
      priority: item.priority,
      firstName: item.firstName,
      lastName: item.lastName,
      company: item.company,
      isSafe: item.isSafe,
      verificationStatus: item.verificationStatus,
      verifiedOn: item.verifiedOn,
      verificationAttempt: item.verificationAttempt,
    })) || [];

  const columns = [
    {
      Header: "Sr No.",
      accessor: "srNo",
      align: "left",
      width: "50px",
    },
    {
      Header: "Email",
      accessor: "email",
      align: "left",
      width: "220px",
      Cell: ({ row }) => {
        const { email, emailProvider } = row.original;

        return (
          <div
            style={{
              lineHeight: "1.3",
              whiteSpace: "normal",
            }}
          >
            <div>{email || "-"}</div>
            <div style={{ marginTop: "1rem" }}>{emailProvider || "-"}</div>
          </div>
        );
      },
    },
    {
      Header: "Name",
      accessor: "fullName",
      align: "left",
      width: "140px",
      Cell: ({ row }) => {
        const { firstName, lastName } = row.original;
        const fullName = `${firstName || ""} ${lastName || ""}`.trim() || "-";
        return wrapCell(fullName, "140px");
      },
    },
    {
      Header: "Designation",
      accessor: "jobTitle",
      align: "left",
      width: "110px",
      Cell: (cell) => cell.value || "-",
    },
    {
      Header: "Company",
      accessor: "company",
      align: "left",
      width: "130px",
      Cell: (cell) => wrapCell(cell.value || "-", "130px"),
    },
    {
      Header: "Location",
      accessor: "location",
      align: "left",
      width: "150px",
      Cell: ({ row }) => {
        const { city, country } = row.original;

        if (city && country) return `${city}, ${country}`;
        return city || country || "-";
      },
    },
    {
      Header: "Status",
      accessor: "isSafe",
      align: "center",
      width: "120px",
      Cell: (cell) => renderStatusBadge(cell.value),
    },
    {
      Header: "Verification",
      accessor: "verificationStatus",
      align: "left",
      width: "150px",
      Cell: ({ row }) => {
        const { verificationStatus, verificationAttempt, verifiedOn, isSafe } =
          row.original;

        const isVerifying =
          verificationStatus === "PENDING" ||
          verificationStatus === "VERIFYING" ||
          (isSafe === "TIMEOUT" && verificationAttempt < 3);

        if (isVerifying) {
          return (
            <MDBox
              display="flex"
              alignItems="center"
              justifyContent="center"
              gap={1}
            >
              <i className="fa fa-spinner fa-spin" />

              <MDTypography variant="caption">Verifying...</MDTypography>
            </MDBox>
          );
        }

        return (
          <div
            style={{
              lineHeight: "1.3",
              whiteSpace: "normal",
            }}
          >
            <div>{verificationStatus || "-"}</div>

            <div style={{ marginTop: "0.5rem" }}>
              {verifiedOn ? formatDateTime(verifiedOn) : "-"}
            </div>
          </div>
        );
      },
    },
    {
      Header: "Action",
      accessor: "actions",
      align: "center",
      disableSortBy: true,

      Cell: ({ row }) => {
        const lead = row.original;
        const isPending =
          lead.verificationStatus === "PENDING" ||
          lead.verificationStatus === "VERIFYING";

        const isTimeoutRetrying =
          lead.isSafe === "TIMEOUT" &&
          Number(lead.verificationAttempt || 0) < 3;

        // Final verification-in-progress state
        const isVerifying = isPending || isTimeoutRetrying;

        // Reverify is allowed only when verification is NOT in progress
        const canReverify =
          !isVerifying &&
          (lead.verificationStatus === "COMPLETED" ||
            lead.verificationStatus === "FAILED" ||
            lead.verificationStatus === "NOT_VERIFIED");

        return (
          <MDBox
            display="flex"
            justifyContent="center"
            alignItems="center"
            gap={1}
          >
            {/* Reverify / Verification in progress */}
            <Tooltip
              title={isVerifying ? "Verification is in progress" : "Reverify"}
              arrow
            >
              <span>
                <Icon
                  sx={({ palette }) => ({
                    cursor: canReverify ? "pointer" : "not-allowed",
                    color: isVerifying
                      ? palette.warning.main
                      : canReverify
                      ? palette.secondary.main
                      : palette.grey[500],
                    opacity: isVerifying || canReverify ? 1 : 0.5,
                  })}
                  onClick={(e) => {
                    e.stopPropagation();
                    if (!canReverify) return;
                    handleReverify(lead);
                  }}
                >
                  {isVerifying ? (
                    <i className="fa fa-spinner fa-spin" />
                  ) : (
                    <i className="fas fa-sync-alt" />
                  )}
                </Icon>
              </span>
            </Tooltip>

            <MDTypography variant="caption">|</MDTypography>

            {/* Delete */}
            <Tooltip title="Delete" arrow>
              <Icon
                sx={{
                  cursor: "pointer",
                  color: "error.main",
                  opacity: 0.7,
                  "&:hover": {
                    opacity: 1,
                  },
                }}
                onClick={(e) => {
                  e.stopPropagation();
                  handleDelete(lead);
                }}
              >
                <i className="fas fa-trash" />
              </Icon>
            </Tooltip>
          </MDBox>
        );
      },
    },
  ];

  const openAddModal = (e) => {
    e.stopPropagation();
    setPanelMode("add");
  };

  const closeModal = () => {
    setPanelMode(null);
  };

  //download csv
  const handleDownloadSafeCSV = async (status = "") => {
    try {
      const res = await DownloadSafeCsv({
        emailListID: decodedId,
        status,
      });

      const url = window.URL.createObjectURL(res.data);
      const link = document.createElement("a");
      link.href = url;
      const fileName = status ? `${status}-leads.csv` : "all-leads.csv";

      link.setAttribute("download", fileName);

      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      notify(res?.message || "Download successfully", "success");
    } catch (err) {
      notify("Failed to download CSV", "error");
    }
  };

  //filters functions
  const buildFilterParams = () => {
    const filterParams = [];

    if (filterState.status) {
      filterParams.push({
        field: "IsSafe",
        condition: "eq",
        filterValues: [filterState.status],
      });
    }

    if (filterState.verificationStatus) {
      filterParams.push({
        field: "VerificationStatus",
        condition: "eq",
        filterValues: [filterState.verificationStatus],
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
      status: null,
      verificationStatus: null,
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

  //delete the leads by id
  const handleDelete = (rowData) => {
    showAlert({
      title: "Are you sure you want to delete this lead?",
      message: `Are you sure you want to permanently delete "${rowData.email}"?`,
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

          const res = await deleteLeadByID(rowData.id);

          if (res?.code === 200) {
            notify(res?.message || "Lead deleted successfully.", "success");

            await getLeadsList(
              dispatch,
              decodedId,
              pageNum,
              filterState.filterParams.filters,
              filterState.filterParams.searchString,
            );
          } else {
            const message = Array.isArray(res?.message)
              ? res.message
                  .map((err) => (err.Msg ? err.Msg : JSON.stringify(err)))
                  .join(", ")
              : res?.message || "Failed to delete lead.";

            notify(message, "error");
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
    });
  };

  const handleReverify = (lead) => {
    showAlert({
      title: "Reverify Lead",
      message: `Do you want to reverify "${lead.email}"?`,
      type: "warning",
      showCancel: true,
      confirmText: "Reverify",
      cancelText: "Cancel",

      onConfirm: async () => {
        try {
          const res = await ReverifyLeadById(lead.id);

          if (res?.code === 200) {
            notify(
              res?.message || "Lead queued for reverification.",
              "success",
            );

            await getLeadsList(
              dispatch,
              decodedId,
              pageNum,
              filterState.filterParams.filters,
              filterState.filterParams.searchString,
            );
          } else {
            const message = Array.isArray(res?.message)
              ? res.message
                  .map((err) => (err.Msg ? err.Msg : JSON.stringify(err)))
                  .join(", ")
              : res?.message || "Failed to reverify lead.";

            notify(message, "error");
          }
        } catch (error) {
          notify(error?.message, "failed to reverify lead");
        } finally {
          setListState((prev) => ({
            ...prev,
            isLoading: false,
          }));
        }
      },
    });
  };

  const statusSummary = [
    {
      label: "Safe",
      count: leadListData?.data?.safeCount || 0,
      color: STATUS_BADGE_COLORS.SAFE,
    },
    {
      label: "Risky",
      count: leadListData?.data?.riskyCount || 0,
      color: STATUS_BADGE_COLORS.RISKY,
    },
    {
      label: "Invalid",
      count: leadListData?.data?.invalidCount || 0,
      color: STATUS_BADGE_COLORS.INVALID,
    },
    {
      label: "Unknown",
      count: leadListData?.data?.unknownCount || 0,
      color: STATUS_BADGE_COLORS.UNKNOWN,
    },
    {
      label: "Pending",
      count: leadListData?.data?.pendingCount || 0,
      color: STATUS_BADGE_COLORS.PENDING,
    },
  ];

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
                <MDBox
                  display="flex"
                  justifyContent="space-between"
                  alignItems="center"
                >
                  <Icon
                    sx={{
                      cursor: "pointer",
                      mr: 1.5,
                      color: "text.secondary",
                    }}
                    onClick={() => navigate("/admin/email-lists")}
                  >
                    arrow_back
                  </Icon>
                  <MDTypography>{emailListName || "Leads"}</MDTypography>
                </MDBox>

                <MDBox
                  display="flex"
                  justifyContent="space-between"
                  alignItems="center"
                  gap={2}
                >
                  {leadListData?.data?.totalRecords > 0 && (
                    <>
                      <MDButton
                        variant="contained"
                        color="deepBlue"
                        startIcon={<Icon>file_download</Icon>}
                        endIcon={
                          <ArrowDropDown
                            sx={{
                              fontSize: "1.7rem",
                              transition: "0.2s",
                              transform: Boolean(anchorEl)
                                ? "rotate(180deg)"
                                : "rotate(0deg)",
                            }}
                          />
                        }
                        onClick={handleOpenDownload}
                        sx={{
                          ...leadListUiStyles.primaryButton,
                          py: 1,
                          px: 2,
                        }}
                      >
                        Download CSV
                      </MDButton>

                      <Menu
                        anchorEl={anchorEl}
                        open={Boolean(anchorEl)}
                        onClose={handleCloseDownload}
                        PaperProps={{
                          sx: {
                            width: anchorEl?.clientWidth,
                            mt: 0.5,
                          },
                        }}
                      >
                        {DOWNLOAD_OPTIONS.map((item) => (
                          <MenuItem
                            key={item.label}
                            onClick={() => {
                              handleCloseDownload();
                              handleDownloadSafeCSV(item.value);
                            }}
                          >
                            {item.label}
                          </MenuItem>
                        ))}
                      </Menu>
                    </>
                  )}
                  <MDButton
                    type="button"
                    variant="contained"
                    color="deepBlue"
                    sx={leadListUiStyles.primaryButton}
                    onClick={openAddModal}
                  >
                    Add Leads
                  </MDButton>
                </MDBox>
              </MDBox>

              {/* Filters */}
              <MDBox px={3} pb={2} onKeyDown={handleKeyDown}>
                <Grid container spacing={2} alignItems="center">
                  <Grid item xs={12} sm={6} md={4} lg={2}>
                    <TextField
                      fullWidth
                      placeholder="Search by email or name "
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
                        {LEADS_STATUS.map((status) => (
                          <MenuItem key={status.label} value={status.label}>
                            {renderStatusBadge(status.label)}
                          </MenuItem>
                        ))}
                      </Select>
                    </FormControl>
                  </Grid>

                  <Grid item xs={12} sm={6} md={4} lg={2}>
                    <FormControl fullWidth size="small" sx={selectSx("100%")}>
                      <Select
                        value={filterState.verificationStatus || ""}
                        displayEmpty
                        onChange={(e) =>
                          setFilterState({
                            ...filterState,
                            verificationStatus: e.target.value,
                          })
                        }
                        IconComponent={ArrowDropDown}
                        renderValue={(selected) => {
                          if (!selected)
                            return (
                              <span style={{ color: "rgba(0, 0, 0, 0.38)" }}>
                                 Select Verification status
                              </span>
                            );

                          return selected;
                        }}
                      >
                        {VERIFICATION_STATUS.map((p) => (
                          <MenuItem key={p.label} value={p.label}>
                            {p.label}
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
                mt={2}
                px={3}
                py={1}
                display="flex"
                flexWrap="wrap"
                alignItems="center"
                gap={1}
              >
                <MDTypography variant="body2" color="deepBlue" mr={2} mb={1}>
                  Total Records: {leadListData?.data?.totalRecords || 0}
                </MDTypography>

                {statusSummary.map((item) => (
                  <MDBox
                    key={item.label}
                    sx={{
                      backgroundColor: item.color.background,
                      color: item.color.text,
                      px: 1.5,
                      py: 0.5,
                      borderRadius: "20px",

                      fontSize: "0.9rem",
                      display: "flex",
                      alignItems: "center",
                      gap: 0.5,
                      mb: 1,
                    }}
                  >
                    <span>{item.label}</span>
                    <span>({item.count})</span>
                  </MDBox>
                ))}
              </MDBox>

              {verificationState.isChecking && (
                <MDBox
                  sx={{
                    mx: 1,
                    mb: 2,
                    p: 2,
                    borderRadius: "10px",
                    display: "flex",
                    alignItems: "center",
                    gap: 2,
                    backgroundColor: "#FFF8E1",
                    border: "1px solid #FFD54F",
                  }}
                >
                  <i
                    className="fa fa-spinner fa-spin"
                    style={{
                      fontSize: 20,
                      color: "#F57C00",
                    }}
                  />

                  <MDBox>
                    <MDTypography fontWeight="bold">
                      Email verification is in progress
                    </MDTypography>

                    <MDTypography variant="caption">
                      {verificationState.pendingCount} lead(s) are being
                      verified. This table updates automatically.
                    </MDTypography>
                  </MDBox>
                </MDBox>
              )}

              <MDBox sx={{ width: "100%", overflowX: "auto" }}>
                {listState.isLoading ? (
                  <MDBox sx={leadListUiStyles.tableState}>
                    <i className="fa fa-spinner fa-spin fa-xl loader" />
                  </MDBox>
                ) : listState.error ? (
                  <MDBox sx={leadListUiStyles.tableState}>
                    <MDTypography variant="body2" color="error">
                      {listState.error}
                    </MDTypography>
                  </MDBox>
                ) : listState.noRecords ? (
                  <MDBox sx={leadListUiStyles.tableState}>
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
                      (leadListData?.data?.totalRecords || 0) /
                        (leadListData?.data?.noOfRecordsPerPage ||
                          DEFAULT_RECORDS_PER_PAGE),
                    )}
                    onRowClick={(lead) => {
                      setSelectedLeadId(lead.id);
                    }}
                  />
                )}
              </MDBox>
            </Card>
          </Grid>
        </Grid>

        <AddLeadModal
          open={panelMode !== null}
          panelMode={panelMode}
          emailListID={decodedId}
          pageNum={pageNum}
          onClose={closeModal}
          setListState={setListState}
        />

        <LeadDetailModal
          open={Boolean(selectedLeadId)}
          leadId={selectedLeadId}
          onClose={() => setSelectedLeadId(null)}
          onUpdateSuccess={() => {
            getLeadsList(
              dispatch,
              decodedId,
              pageNum,
              filterState.filterParams.filters,
              filterState.filterParams.searchString,
            );
          }}
        />
      </MDBox>
    </DashboardLayout>
  );
}
export default LeadListing;
