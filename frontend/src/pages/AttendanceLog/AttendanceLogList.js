import React, { useEffect, useState } from "react";
import {
  Card,
  FormControl,
  Grid,
  Icon,
  MenuItem,
  Select,
  TextField,
} from "@mui/material";
import { DEFAULT_RECORDS_PER_PAGE } from "components/common/constant";
import {
  wrapCell,
  formatDateTime,
} from "components/CommonComponent/CommonFunction";
import DataTable from "examples/Tables/DataTable";
import { useDispatch, useSelector } from "react-redux";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import MDBox from "components/MDBox";
import MDButton from "components/MDButton";
import MDTypography from "components/MDTypography";
import { getAttendanceLogList } from "actions/attendanceLog";
import { attendanceLogList } from "constants/attendanceLog";
import { ATTENDANCE_LOG_STATE } from "components/common/constant";
import { selectSx } from "components/CommonComponent/CommonFunction";
import { ArrowDropDown } from "@mui/icons-material";

function AttendanceLogList() {
  const dispatch = useDispatch();
  const attendanceLogListData = useSelector(
    (state) => state.attendanceLog.attendanceLogList,
  );
  const [listState, setListState] = useState({
    isLoading: true,
    error: "",
    noRecords: false,
    showLoaderOnClick: false,
  });
  const [filterState, setFilterState] = useState({
    searchString: "",
    attendanceState: null,
    startDate: "",
    endDate: "",
    filterParams: {
      filters: [],
      searchString: "",
    },
  });
  const [pageNum, setPageNum] = useState(1);

  useEffect(() => {
    getAttendanceLogList(
      dispatch,
      pageNum,
      filterState.filterParams.filters,
      filterState.filterParams.searchString,
    );
    return () => {
      dispatch({ type: attendanceLogList, payload: {} });
    };
  }, [
    dispatch,
    pageNum,
    filterState.filterParams.filters,
    filterState.filterParams.searchString,
  ]);

  useEffect(() => {
    if (attendanceLogListData?.code === 200) {
      setListState({
        isLoading: false,
        error: "",
        noRecords: attendanceLogListData?.data?.AttendanceLog?.length === 0,
        showLoaderOnClick: false,
      });
    } else if (attendanceLogListData?.code || attendanceLogListData?.message) {
      setListState({
        isLoading: false,
        error:
          attendanceLogListData?.message ||
          "Failed to load attendance log list! Please try again.",
        noRecords: false,
        showLoaderOnClick: false,
      });
    }
  }, [attendanceLogListData]);

  const handleKeyDown = (e) => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleSearch();
    }
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

    if (filterState.startDate && filterState.endDate) {
      filterParams.push({
        field: "Timestamp",
        condition: "btw",
        filterValues: [filterState.startDate, filterState.endDate],
      });
    } else if (filterState.startDate) {
      filterParams.push({
        field: "Timestamp",
        condition: "gteq",
        filterValues: [filterState.startDate],
      });
    } else if (filterState.endDate) {
      filterParams.push({
        field: "Timestamp",
        condition: "lteq",
        filterValues: [`${filterState.endDate} 23:59:59`],
      });
    }

    if (filterState.attendanceState) {
      filterParams.push({
        field: "AttendanceState",
        condition: "eq",
        filterValues: [filterState.attendanceState],
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
      attendanceState: null,
      startDate: "",
      endDate: "",
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

  // Table data columns
  const formattedTableData =
    attendanceLogListData?.data?.AttendanceLog?.map((item, index) => ({
      ...item,
      id: item?.ID,
      srNo:
        (pageNum - 1) *
          (attendanceLogListData?.data?.noOfRecordsPerPage ||
            DEFAULT_RECORDS_PER_PAGE) +
        index +
        1,
      uID: item?.uID ?? "-",
      empID: item?.empID || "-",
      empName: item?.empName || "-",
      timestamp: item?.timestamp || null,
      status: item?.status ?? "-",
      punch: item?.punch ?? "-",
      attendanceState: item?.attendanceState || "-",
      deviceName: item?.deviceName || "-",
    })) || [];

  const columns = [
    {
      Header: "Sr No.",
      accessor: "srNo",
      align: "left",
      width: "85px",
    },
    {
      Header: "Employee ID",
      accessor: "empID",
      align: "left",
      width: "200px",
    },
    {
      Header: "Employee Name",
      accessor: "empName",
      align: "left",
      width: "200px",
      Cell: (cell) => wrapCell(cell.value || "-", "200px"),
    },
    {
      Header: "Device Name",
      accessor: "deviceName",
      align: "left",
      width: "200px",
      Cell: (cell) => wrapCell(cell.value || "-", "200px"),
    },
    {
      Header: "Status",
      accessor: "status",
      align: "left",
      width: "100px",
    },
    {
      Header: "Punch",
      accessor: "punch",
      align: "left",
      width: "100px",
    },
    {
      Header: "Attendance State",
      accessor: "attendanceState",
      align: "left",
      width: "150px",
    },
    {
      Header: "Punch Time",
      accessor: "timestamp",
      align: "left",
      width: "160px",
      Cell: (cell) => (cell.value ? formatDateTime(cell.value) : "-"),
    },
  ];

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
                    Attendance Log Lists
                  </MDTypography>
                </MDBox>

                {/* Filters */}
                <MDBox px={3} pb={2} onKeyDown={handleKeyDown}>
                  <Grid container spacing={2} alignItems="center">
                    <Grid item xs={12} sm={6} md={4} lg={2}>
                      <TextField
                        fullWidth
                        placeholder="Search by emp ID or device name "
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
                          value={filterState.attendanceState || ""}
                          displayEmpty
                          onChange={(e) =>
                            setFilterState({
                              ...filterState,
                              attendanceState: e.target.value,
                            })
                          }
                          IconComponent={ArrowDropDown}
                          renderValue={(selected) => {
                            if (!selected) {
                              return (
                                <span style={{ color: "rgba(0, 0, 0, 0.38)" }}>
                                  Select Attendance Log state
                                </span>
                              );
                            }

                            const selectedState = ATTENDANCE_LOG_STATE.find(
                              (p) => p.label === selected,
                            );
                            return selectedState?.value || selected;
                          }}
                        >
                          {ATTENDANCE_LOG_STATE.map((p) => (
                            <MenuItem key={p.label} value={p.label}>
                              {p.value}
                            </MenuItem>
                          ))}
                        </Select>
                      </FormControl>
                    </Grid>

                    <Grid item xs={12} sm={6} md={4} lg={2}>
                      <TextField
                        fullWidth
                        type="date"
                        label="Start Date"
                        size="small"
                        InputLabelProps={{ shrink: true }}
                        value={filterState.startDate}
                        onChange={(e) =>
                          setFilterState({
                            ...filterState,
                            startDate: e.target.value,
                          })
                        }
                      />
                    </Grid>

                    <Grid item xs={12} sm={6} md={4} lg={2}>
                      <TextField
                        fullWidth
                        type="date"
                        label="End Date"
                        size="small"
                        InputLabelProps={{ shrink: true }}
                        inputProps={{ min: filterState.startDate || undefined }}
                        value={filterState.endDate}
                        onChange={(e) =>
                          setFilterState({
                            ...filterState,
                            endDate: e.target.value,
                          })
                        }
                      />
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
                            variant="gradient"
                            color="light"
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
                    Total Records:{" "}
                    {attendanceLogListData?.data?.totalRecords || 0}
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
                        (attendanceLogListData?.data?.totalRecords || 0) /
                          (attendanceLogListData?.data?.noOfRecordsPerPage ||
                            DEFAULT_RECORDS_PER_PAGE),
                      )}
                      //   onRowClick={handleRowClick}
                    />
                  )}
                </MDBox>
              </Card>
            </Grid>
          </Grid>
        </MDBox>
      </DashboardLayout>
    </>
  );
}

export default AttendanceLogList;
