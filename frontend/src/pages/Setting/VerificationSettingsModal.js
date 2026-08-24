import { useEffect, useState } from "react";
import { Card, Grid } from "@mui/material";

import SettingsOutlinedIcon from "@mui/icons-material/SettingsOutlined";
import InfoOutlinedIcon from "@mui/icons-material/InfoOutlined";

import MDBox from "components/MDBox";
import MDButton from "components/MDButton";
import MDTypography from "components/MDTypography";
import MDInput from "components/MDInput";

import {
  saveUserSettingConfiguration,
  getUserSettingList,
  updateUserSettingConfiguration,
} from "actions/userSetting";
import { useNotify } from "components/CommonComponent/NotificationProvider";
import DashboardLayout from "examples/LayoutContainers/DashboardLayout";

function VerificationSettingsModal() {
  const [settingId, setSettingId] = useState(null);
  const [interval, setInterval] = useState(15);
  const [loading, setLoading] = useState(false);

  const { notify } = useNotify();

  useEffect(() => {
    const fetchUserSetting = async () => {
      try {
        const res = await getUserSettingList();

        if (res?.code === 200) {
          const settings = res?.data?.UserSetting || [];

          if (settings.length > 0) {
            const setting = settings[0];

            setSettingId(setting.ID);
            setInterval(setting.verificationInterval || 15);
          } else {
            setSettingId(null);
            setInterval(15);
          }
        }
      } catch (error) {
        console.error(error);
      }
    };

    fetchUserSetting();
  }, []);

  const handleSave = async () => {
    if (!interval || interval < 10 || interval > 60) {
      return;
    }

    setLoading(true);
    try {
      const bodyData = {
        verificationInterval: interval,
      };

      let res;
      if (settingId) {
        res = await updateUserSettingConfiguration(settingId, bodyData);
      } else {
        res = await saveUserSettingConfiguration(bodyData);
      }

      if (res?.code === 200) {
        notify(res?.message || "Verification interval saved successfully");
      } else {
        const message = Array.isArray(res?.message)
          ? res.message
              .map((err) => (err.Msg ? err.Msg : JSON.stringify(err)))
              .join(", ")
          : res?.message || "Failed to save the verification interval";

        notify(message, "error");
      }
    } catch (error) {
      notify(error?.message || "Something went wrong", "error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <DashboardLayout>
      <MDBox>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <Card>
              <MDBox p={3}>
                <MDBox display="flex" alignItems="center" mb={1}>
                  <SettingsOutlinedIcon color="primary" sx={{ mr: 1 }} />

                  <MDTypography variant="h4" fontWeight="bold">
                    Verification Schedule
                  </MDTypography>
                </MDBox>

                {/* Description */}
                <MDTypography
                  variant="button"
                  color="text"
                  display="block"
                  mb={4}
                >
                  Choose how frequently Mailora verifies each pending email
                  record.
                </MDTypography>

                <MDBox
                  border="0.063rem solid"
                  borderColor="#E5E7EB"
                  borderRadius="0.625rem"
                  p={2.5}
                >
                  <MDTypography variant="h6" fontWeight="bold">
                    Verification Interval
                  </MDTypography>

                  <MDTypography
                    variant="caption"
                    color="text"
                    display="block"
                    mb={2}
                  >
                    Set the time interval between processing pending email
                    records.
                  </MDTypography>

                  <Grid container spacing={2}>
                    <Grid item xs={12} md={6}>
                      <MDInput
                        fullWidth
                        type="number"
                        value={interval}
                        placeholder="Enter interval"
                        inputProps={{
                          min: 10,
                          max: 60,
                        }}
                        onChange={(e) => {
                          const value = e.target.value;

                          if (value === "") {
                            setInterval("");
                            return;
                          }

                          if (value.length > 2) return;

                          setInterval(Number(value));
                        }}
                      />
                    </Grid>
                  </Grid>

                  {(!interval || interval < 10 || interval > 60) && (
                    <MDTypography
                      variant="caption"
                      color="error"
                      mt={1}
                      display="block"
                    >
                      Verification interval must be between 10 and 60 seconds.
                    </MDTypography>
                  )}

                  <MDBox
                    mt={2}
                    display="flex"
                    alignItems="center"
                    flexWrap="wrap"
                    gap={1}
                  >
                    <MDTypography variant="caption" color="text">
                      Quick Select:
                    </MDTypography>

                    {[15, 30, 45, 60].map((item) => (
                      <MDButton
                        key={item}
                        size="small"
                        color="deepBlue"
                        variant={interval === item ? "contained" : "outlined"}
                        onClick={() => setInterval(item)}
                      >
                        {item}s
                      </MDButton>
                    ))}
                  </MDBox>
                </MDBox>

                <MDBox
                  mt={3}
                  p={2}
                  borderRadius="0.5rem"
                  bgcolor="#F4F8FF"
                  border="0.063rem solid #D6E4FF"
                >
                  <MDBox display="flex" alignItems="flex-start">
                    <InfoOutlinedIcon
                      color="primary"
                      fontSize="small"
                      sx={{ mr: 1, mt: "2px" }}
                    />

                    <MDTypography variant="caption" color="text">
                      <strong>Verification Schedule</strong>
                      <br />
                      Mailora checks your pending email queue at the configured
                      interval and automatically processes emails awaiting
                      verification.
                    </MDTypography>
                  </MDBox>
                </MDBox>

                <MDBox mt={4} display="flex" justifyContent="flex-start" gap={2}>
                  <MDButton variant="gradient" color="light" disabled={loading}>
                    Cancel
                  </MDButton>

                  <MDButton
                    variant="contained"
                    color="deepBlue"
                    disabled={
                      loading || !interval || interval < 10 || interval > 60
                    }
                    onClick={handleSave}
                  >
                    {loading ? (
                      <>
                        <i
                          className="fa fa-spinner fa-spin"
                          style={{ marginRight: 8 }}
                        />
                        Saving...
                      </>
                    ) : (
                      "Save Changes"
                    )}
                  </MDButton>
                </MDBox>
              </MDBox>
            </Card>
          </Grid>
        </Grid>
      </MDBox>
    </DashboardLayout>
  );
}

export default VerificationSettingsModal;
