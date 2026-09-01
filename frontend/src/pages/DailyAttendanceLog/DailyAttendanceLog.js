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
import { useDispatch, useSelector } from "react-redux";
import { getDailyAttendanceLogList } from "actions/attendanceLog";
import { dailyAttendanceLogList } from "constants/attendanceLog";
import { DEFAULT_RECORDS_PER_PAGE } from "components/common/constant";
import {
  wrapCell,
  formatDate,
} from "components/CommonComponent/CommonFunction";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import MDBox from "components/MDBox";
import MDButton from "components/MDButton";
import MDTypography from "components/MDTypography";
import DataTable from "examples/Tables/DataTable";

function DailyAttendanceLog() {
  const dispatch = useDispatch();
  const dailyAttendanceLogListData = useSelector(
    (state) => state.attendanceLog.dailyAttendanceLogList,
  );
  const [listState, setListState] = useState({
    isLoading: true,
    error: "",
    noRecords: false,
    showLoaderOnClick: false,
  });
  const [pageNum, setPageNum] = useState(1);
  const [filterState, setFilterState] = useState({
    searchString: "",
    filterParams: {
      filters: [],
      searchString: "",
    },
  });

  useEffect(() => {
    getDailyAttendanceLogList(
      dispatch,
      pageNum,
      filterState.filterParams.filters,
      filterState.filterParams.searchString,
    );
    return () => {
      dispatch({ type: dailyAttendanceLogList, payload: {} });
    };
  }, [
    dispatch,
    pageNum,
    filterState.filterParams.filters,
    filterState.filterParams.searchString,
  ]);

  useEffect(() => {
    if (dailyAttendanceLogListData?.code === 200) {
      setListState({
        isLoading: false,
        error: "",
        noRecords:
          dailyAttendanceLogListData?.data?.AttendanceLog?.length === 0,
        showLoaderOnClick: false,
      });
    } else if (
      dailyAttendanceLogListData?.code ||
      dailyAttendanceLogListData?.message
    ) {
      setListState({
        isLoading: false,
        error:
          dailyAttendanceLogListData?.message ||
          "Failed to load daily attendance log list! Please try again.",
        noRecords: false,
        showLoaderOnClick: false,
      });
    }
  }, [dailyAttendanceLogListData]);

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

  const buildFilterParams = () => {
    const filterParams = [];

    return {
      filters: filterParams.length ? JSON.stringify(filterParams) : [],
      searchString: filterState.searchString.trim(),
    };
  };

  const handleClear = () => {
    setFilterState({
      searchString: "",
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
    dailyAttendanceLogListData?.data?.dailyAttendanceLog?.map(
      (item, index) => ({
        ...item,
        id: item?.ID,
        srNo:
          (pageNum - 1) *
            (dailyAttendanceLogListData?.data?.noOfRecordsPerPage ||
              DEFAULT_RECORDS_PER_PAGE) +
          index +
          1,
        empID: item?.empID || "-",
        empName: item?.empName || "-",
        date: item?.date || null,
        checkIn:
          item?.punches?.map((punch) => punch?.checkIn || "-").join(", ") ||
          "-",
        checkOut:
          item?.punches?.map((punch) => punch?.checkOut || "-").join(", ") ||
          "-",
        workingHours: item?.workingHours || "-",
        status: item?.status || "-",
      }),
    ) || [];

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
      Header: "Date",
      accessor: "date",
      align: "left",
      width: "120px",
      Cell: (cell) => (cell.value ? formatDate(cell.value) : "-"),
    },
    {
      Header: "Check In",
      accessor: "checkIn",
      align: "left",
      width: "130px",
    },
    {
      Header: "Check Out",
      accessor: "checkOut",
      align: "left",
      width: "130px",
    },
    {
      Header: "Working Hours",
      accessor: "workingHours",
      align: "left",
      width: "150px",
    },
    {
      Header: "Status",
      accessor: "status",
      align: "left",
      width: "120px",
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
                    Daily Attendance Log Lists
                  </MDTypography>
                </MDBox>

                {/* Filters */}
                <MDBox px={3} pb={2} onKeyDown={handleKeyDown}>
                  <Grid container spacing={2} alignItems="center">
                    <Grid item xs={12} sm={6} md={4} lg={2}>
                      <TextField
                        fullWidth
                        placeholder="Search by emp ID or name "
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
                    {dailyAttendanceLogListData?.data?.totalRecords || 0}
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
                        (dailyAttendanceLogListData?.data?.totalRecords || 0) /
                          (dailyAttendanceLogListData?.data
                            ?.noOfRecordsPerPage || DEFAULT_RECORDS_PER_PAGE),
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

export default DailyAttendanceLog;
