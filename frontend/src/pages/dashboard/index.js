/**
=========================================================
* Material Dashboard 2 React - v2.2.0
=========================================================

* Product Page: https://www.creative-tim.com/product/material-dashboard-react
* Copyright 2023 Creative Tim (https://www.creative-tim.com)

Coded by www.creative-tim.com

 =========================================================

* The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
*/

// @mui material components
import Grid from "@mui/material/Grid";

// Material Dashboard 2 React components
import MDBox from "components/MDBox";

// Material Dashboard 2 React example components
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";
import DashboardNavbar from "examples/Navbars/DashboardNavbar";
import ReportsBarChart from "examples/Charts/BarCharts/ReportsBarChart";
import ReportsLineChart from "examples/Charts/LineCharts/ReportsLineChart";
import ComplexStatisticsCard from "examples/Cards/StatisticsCards/ComplexStatisticsCard";

import reportsLineChartData from "./data/reportsLineChartData";
import reportsBarChartData from "./data/reportsBarChartData";
import { useDispatch, useSelector } from "react-redux";
import { getEmailLists } from "actions/email";
import { emailList } from "constants/email";
import { useEffect, useState } from "react";
import MDTypography from "components/MDTypography";
import { Card } from "@mui/material";

function Dashboard() {
  const dispatch = useDispatch();
  const emailListData = useSelector((state) => state.email.emailList);
  const { sales, tasks } = reportsLineChartData;
  const [listState, setListState] = useState({
    isLoading: true,
    error: "",
    noRecords: false,
    showLoaderOnClick: false,
  });

  useEffect(() => {
    getEmailLists(dispatch);
    return () => {
      dispatch({ type: emailList, payload: {} });
    };
  }, [dispatch]);

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

  const summary = emailListData?.data?.emailList?.reduce(
    (acc, item) => {
      acc.safe += item.safeCount;
      acc.risky += item.riskyCount;
      acc.invalid += item.invalIDCount;
      acc.unknown += item.unknownCount;
      acc.pending += item.pendingCount;
      return acc;
    },
    {
      safe: 0,
      risky: 0,
      invalid: 0,
      unknown: 0,
      pending: 0,
    },
  );

  return (
    <DashboardLayout>
      <DashboardNavbar />
      <MDBox py={3}>
        <MDBox>
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
          ) : (
            <>
              <Grid container spacing={3}>
                <Grid item xs={12} sm={6} md={3}>
                  <ComplexStatisticsCard
                    color="success"
                    icon="verified"
                    title="Safe"
                    count={listState.noRecords ? "-" : summary?.safe ?? 0}
                    percentage={{
                      color: "success",
                      amount: "",
                      label: listState.noRecords
                        ? "No emails verified yet"
                        : "Verified emails",
                    }}
                  />
                </Grid>

                <Grid item xs={12} sm={6} md={3}>
                  <ComplexStatisticsCard
                    color="warning"
                    icon="warning"
                    title="Risky"
                    count={listState.noRecords ? "-" : summary?.risky ?? 0}
                    percentage={{
                      color: "warning",
                      amount: "",
                      label: listState.noRecords
                        ? "No emails verified yet"
                        : "Require attention",
                    }}
                  />
                </Grid>

                <Grid item xs={12} sm={6} md={3}>
                  <ComplexStatisticsCard
                    color="info"
                    icon="help_outline"
                    title="Unknown"
                    count={listState.noRecords ? "-" : summary?.unknown ?? 0}
                    percentage={{
                      color: "info",
                      amount: "",
                      label: listState.noRecords
                        ? "No pending emails"
                        : "Pending verification",
                    }}
                  />
                </Grid>

                <Grid item xs={12} sm={6} md={3}>
                  <ComplexStatisticsCard
                    color="error"
                    icon="cancel"
                    title="Invalid"
                    count={listState.noRecords ? "-" : summary?.invalid ?? 0}
                    percentage={{
                      color: "error",
                      amount: "",
                      label: listState.noRecords
                        ? "No emails verified yet"
                        : "Invalid emails",
                    }}
                  />
                </Grid>
              </Grid>

              {/* {listState.noRecords && (
                <Card>
                  <MDBox
                    display="flex"
                    justifyContent="center"
                    alignItems="center"
                    flexDirection="column"
                    mt={4}
                  >
                    <MDTypography variant="h6">
                      No verification data yet
                    </MDTypography>

                    <MDTypography variant="body2" color="text">
                      Upload your first email list to start verifying emails.
                    </MDTypography>
                  </MDBox>
                </Card>
              )} */}
            </>
          )}

          <Grid container spacing={3} mt={6}>
            <Grid item xs={12} md={6} lg={4}>
              <MDBox mb={3}>
                <ReportsBarChart
                  color="info"
                  title="website views"
                  description="Last Campaign Performance"
                  date="campaign sent 2 days ago"
                  chart={reportsBarChartData}
                />
              </MDBox>
            </Grid>
            <Grid item xs={12} md={6} lg={4}>
              <MDBox mb={3}>
                <ReportsLineChart
                  color="success"
                  title="daily sales"
                  description={
                    <>
                      (<strong>+15%</strong>) increase in today sales.
                    </>
                  }
                  date="updated 4 min ago"
                  chart={sales}
                />
              </MDBox>
            </Grid>
            <Grid item xs={12} md={6} lg={4}>
              <MDBox mb={3}>
                <ReportsLineChart
                  color="dark"
                  title="completed tasks"
                  description="Last Campaign Performance"
                  date="just updated"
                  chart={tasks}
                />
              </MDBox>
            </Grid>
          </Grid>
        </MDBox>
      </MDBox>
    </DashboardLayout>
  );
}

export default Dashboard;
