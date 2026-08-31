import React, { useEffect, useState } from "react";
import { Card, Grid, Icon, TextField } from "@mui/material";
import { getEmployeeList } from "actions/employee";
import { DEFAULT_RECORDS_PER_PAGE } from "components/common/constant";
import { wrapCell, selectSx } from "components/CommonComponent/CommonFunction";
import { employeeList } from "constants/employee";
import DataTable from "examples/Tables/DataTable";
import { useDispatch, useSelector } from "react-redux";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import MDBox from "components/MDBox";
import MDButton from "components/MDButton";
import MDTypography from "components/MDTypography";

function EmployeeList() {
  const dispatch = useDispatch();
  const employeeListData = useSelector((state) => state.employee.employeeList);
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
  const [pageNum, setPageNum] = useState(1);

  useEffect(() => {
    getEmployeeList(
      dispatch,
      pageNum,
      filterState.filterParams.filters,
      filterState.filterParams.searchString,
    );
    return () => {
      dispatch({ type: employeeList, payload: {} });
    };
  }, [
    dispatch,
    pageNum,
    filterState.filterParams.filters,
    filterState.filterParams.searchString,
  ]);

  useEffect(() => {
    if (employeeListData?.code === 200) {
      setListState({
        isLoading: false,
        error: "",
        noRecords: employeeListData?.data?.Employee?.length === 0,
        showLoaderOnClick: false,
      });
    } else if (employeeListData?.code || employeeListData?.message) {
      setListState({
        isLoading: false,
        error:
          employeeListData?.message ||
          "Failed to load employee list! Please try again.",
        noRecords: false,
        showLoaderOnClick: false,
      });
    }
  }, [employeeListData]);

  // Table data columns
  const formattedTableData =
    employeeListData?.data?.Employee?.map((item, index) => ({
      ...item,
      id: item?.ID,
      srNo:
        (pageNum - 1) *
          (employeeListData?.data?.noOfRecordsPerPage ||
            DEFAULT_RECORDS_PER_PAGE) +
        index +
        1,
      empID: item?.empID || "-",
      empName: item?.empName?.trim() || "-",
      privilege: item?.privilege ?? 0,
      groupID: item?.groupID || "-",
      card: item?.card || "-",
      uID: item?.uID ?? "-",
    })) || [];

  const handleKeyDown = (e) => {
    if (e.key === "Enter") {
      e.preventDefault();
        handleSearch();
    }
  };

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
      Header: "Name",
      accessor: "empName",
      align: "left",
      width: "200px",
      Cell: (cell) => wrapCell(cell.value || "-", "200px"),
    },
    {
      Header: "Privilege",
      accessor: "privilege",
      align: "left",
      width: "100px",
    },
    {
      Header: "Group ID",
      accessor: "groupID",
      align: "left",
      width: "120px",
    },
    {
      Header: "Card",
      accessor: "card",
      align: "left",
      width: "120px",
    },
  ];

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
                    Employees Lists
                  </MDTypography>
                </MDBox>

                {/* Filters */}
                <MDBox px={3} pb={2} onKeyDown={handleKeyDown}>
                  <Grid container spacing={2} alignItems="center">
                    <Grid item xs={12} sm={6} md={4} lg={2}>
                      <TextField
                        fullWidth
                        placeholder="Search by Emp ID or name"
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
                    Total Records: {employeeListData?.data?.totalRecords || 0}
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
                        (employeeListData?.data?.totalRecords || 0) /
                          (employeeListData?.data?.noOfRecordsPerPage ||
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

export default EmployeeList;
