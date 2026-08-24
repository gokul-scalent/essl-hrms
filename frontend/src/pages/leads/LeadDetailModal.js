import React, { useEffect, useState } from "react";
import { Grid, Icon, Radio, RadioGroup, FormControlLabel } from "@mui/material";
import CenterPopup from "components/CommonComponent/CenterPopup";
import MDBox from "components/MDBox";
import MDTypography from "components/MDTypography";
import MDInput from "components/MDInput";
import MDButton from "components/MDButton";
import {
  formatDateTime,
  radioLabelSx,
} from "components/CommonComponent/CommonFunction";
import { useNotify } from "components/CommonComponent/NotificationProvider";
import { getLeadByID, updateLeadListByID } from "actions/leads";
import { renderPriorityBadge } from "pages/emailList/EmailListing";
import { renderStatusBadge } from "./LeadListing";
import colors from "assets/theme/base/colors";
import { PRIORITY } from "components/common/constant";
import { getLeadsList } from "actions/leads";

const EditableItem = ({ label, value, onChange }) => (
  <Grid item xs={12} md={4}>
    <MDTypography
      variant="body2"
      sx={{
        fontWeight: 600,
        color: colors.grey[600],
        mb: 0.75,
        display: "block",
      }}
    >
      {label}
    </MDTypography>

    <MDInput
      fullWidth
      value={value}
      onChange={onChange}
      size="small"
      inputProps={{ maxLength: 50 }}
    />
  </Grid>
);

const DetailItem = ({ label, value }) => (
  <Grid item xs={12} md={4}>
    <MDTypography
      variant="body2"
      sx={{
        fontWeight: 600,
        color: colors.grey[600],
        mb: 0.5,
        display: "block",
      }}
    >
      {label}
    </MDTypography>

    <MDTypography
      variant="body2"
      sx={{
        color: "deepBlue",
        minHeight: 24,
        display: "flex",
        alignItems: "center",
        fontWeight: 400,
      }}
    >
      {value || "-"}
    </MDTypography>
  </Grid>
);

const getEditData = (lead) => ({
  firstName: lead?.firstName || "",
  lastName: lead?.lastName || "",
  email: lead?.email || "",
  priority: lead?.priority?.toLowerCase() || "",
  company: lead?.company || "",
  jobTitle: lead?.jobTitle || "",
  industry: lead?.industry || "",
  city: lead?.city || "",
  country: lead?.country || "",
});

const getErrorMessage = (response, defaultMessage) => {
  if (Array.isArray(response?.message)) {
    return response.message
      .map((err) => (err.Msg ? err.Msg : JSON.stringify(err)))
      .join(", ");
  }

  return response?.message || defaultMessage;
};

function LeadDetailModal({ open, leadId, onClose, onUpdateSuccess }) {
  const { notify } = useNotify();

  const [leadDetails, setLeadDetails] = useState(null);
  const [editData, setEditData] = useState({});
  const [loading, setLoading] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    if (!open || !leadId) return;
    setIsEditing(false);

    const fetchLeadDetails = async () => {
      try {
        setLoading(true);
        const res = await getLeadByID(leadId);
        if (res?.code === 200) {
          setLeadDetails(res.data);
        } else {
          const message = Array.isArray(res?.message)
            ? res.message
                .map((err) => (err.Msg ? err.Msg : JSON.stringify(err)))
                .join(", ")
            : res?.message || "Failed to fetch lead details";

          notify(message, "error");
          onClose();
        }
      } catch (err) {
        const message = Array.isArray(err?.response?.data?.message)
          ? err.response.data.message
              .map((e) => (e.Msg ? e.Msg : JSON.stringify(e)))
              .join(", ")
          : err?.response?.data?.message || "Failed to fetch lead details";

        notify(message, "error");
        onClose();
      } finally {
        setLoading(false);
      }
    };

    fetchLeadDetails();
  }, [open, leadId]);

  const handleEdit = () => {
    setEditData(getEditData(leadDetails));
    setIsEditing(true);
  };

  const handleCancel = () => {
    setEditData(getEditData(leadDetails));
    setIsEditing(false);
  };

  const handleChange = (field) => (event) => {
    setEditData((prev) => ({
      ...prev,
      [field]: event.target.value,
    }));
  };

  const emailChanged =
    isEditing &&
    editData.email?.trim().toLowerCase() !==
      leadDetails?.email?.trim().toLowerCase();

  const handleSave = async () => {
    try {
      setLoading(true);

      const bodyData = {
        firstName: editData.firstName,
        lastName: editData.lastName,
        email: editData.email,
        priority: editData.priority?.toUpperCase(),
        company: editData.company,
        jobTitle: editData.jobTitle,
        industry: editData.industry,
        city: editData.city,
        country: editData.country,
      };

      const res = await updateLeadListByID(leadId, bodyData);

      if (res?.code === 200) {
        const updatedRes = await getLeadByID(leadId);
        if (updatedRes?.code === 200) {
          setLeadDetails(updatedRes.data);
          setEditData(getEditData(updatedRes.data));
        }
        setIsEditing(false);
        onClose();
        //after successfully update call get api
        await onUpdateSuccess?.();
        notify(res.message || "Lead updated successfully.", "success");
      } else {
        notify(getErrorMessage(res, "Failed to update lead"), "error");
      }
    } catch (err) {
      notify(
        getErrorMessage(err?.response?.data, "Failed to update lead"),
        "error",
      );
    } finally {
      setLoading(false);
    }
  };

  if (!open) return null;

  const isPending = leadDetails?.verificationStatus === "PENDING";
  const isCompleted = leadDetails?.verificationStatus === "COMPLETED";
  const isFailed = leadDetails?.verificationStatus === "FAILED";

  return (
    <CenterPopup isOpen={open} onClose={onClose} width={950}>
      <MDBox p={1}>
        <MDBox
          display="flex"
          justifyContent="space-between"
          alignItems="center"
          mb={3}
          pb={2}
          sx={{
            borderBottom: "0.063rem solid #e5e7eb",
          }}
        >
          <MDTypography variant="h6">Lead Details</MDTypography>

          <MDBox display="flex" alignItems="center" gap={1}>
            {isEditing ? (
              <>
                <MDButton
                  variant="outlined"
                  color="secondary"
                  onClick={handleCancel}
                  disabled={loading}
                >
                  Cancel
                </MDButton>

                <MDButton
                  variant="contained"
                  color="deepBlue"
                  onClick={handleSave}
                  disabled={loading}
                >
                  {loading ? <i className="fa fa-spinner fa-spin" /> : "Save"}
                </MDButton>
              </>
            ) : (
              <MDButton
                variant="contained"
                color="deepBlue"
                startIcon={<Icon>edit</Icon>}
                onClick={handleEdit}
              >
                Edit
              </MDButton>
            )}

            <Icon
              onClick={onClose}
              sx={{
                cursor: "pointer",
                color: "text.secondary",
                ml: 0.5,
              }}
            >
              close
            </Icon>
          </MDBox>
        </MDBox>
        {loading ? (
          <MDBox display="flex" justifyContent="center" py={5}>
            <i className="fa fa-spinner fa-spin fa-xl loader" />
          </MDBox>
        ) : leadDetails ? (
          <Grid container spacing={4}>
            {isEditing ? (
              <>
                {/* Editable fields */}
                <EditableItem
                  label="First Name"
                  value={editData.firstName}
                  onChange={handleChange("firstName")}
                />

                <EditableItem
                  label="Last Name"
                  value={editData.lastName}
                  onChange={handleChange("lastName")}
                />

                <EditableItem
                  label="Email"
                  value={editData.email}
                  onChange={handleChange("email")}
                />

                <Grid item xs={12} md={4}>
                  <MDTypography
                    variant="body2"
                    sx={{
                      fontWeight: 600,
                      color: colors.grey[600],
                      mb: 0.75,
                      display: "block",
                    }}
                  >
                    Priority
                  </MDTypography>

                  <RadioGroup
                    row
                    name="priority"
                    value={editData.priority || ""}
                    onChange={handleChange("priority")}
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
                </Grid>

                <EditableItem
                  label="Company"
                  value={editData.company}
                  onChange={handleChange("company")}
                />

                <EditableItem
                  label="Job Title"
                  value={editData.jobTitle}
                  onChange={handleChange("jobTitle")}
                />

                <EditableItem
                  label="Industry"
                  value={editData.industry}
                  onChange={handleChange("industry")}
                />

                <EditableItem
                  label="City"
                  value={editData.city}
                  onChange={handleChange("city")}
                />

                <EditableItem
                  label="Country"
                  value={editData.country}
                  onChange={handleChange("country")}
                />

                {/* Email change information */}
                {emailChanged && (
                  <Grid item xs={12}>
                    <MDBox
                      display="flex"
                      alignItems="center"
                      gap={1}
                      px={1.5}
                      py={1}
                      sx={{
                        backgroundColor: "#f5f8ff",
                        border: "0.063rem solid #dbe5ff",
                        borderRadius: "0.5rem",
                      }}
                    >
                      <Icon
                        sx={{
                          color: "deepBlue",
                          fontSize: "1.2rem",
                        }}
                      >
                        info
                      </Icon>

                      <MDTypography
                        variant="caption"
                        color="deepBlue"
                        fontWeight="600"
                      >
                        Changing the email address will require re-verification.
                      </MDTypography>
                    </MDBox>
                  </Grid>
                )}

                {/* Verification information while editing */}
                <DetailItem
                  label="Verification"
                  value={leadDetails.verificationStatus}
                />

                {isPending && (
                  <Grid item xs={12}>
                    <MDBox
                      display="flex"
                      alignItems="center"
                      gap={1}
                      px={1.5}
                      py={1.5}
                      sx={{
                        backgroundColor: "#f5f8ff",
                        border: "0.063rem solid #dbe5ff",
                        borderRadius: "0.5rem",
                      }}
                    >
                      <Icon
                        sx={{
                          color: "deepBlue",
                          fontSize: "1.2rem",
                        }}
                      >
                        pending
                      </Icon>

                      <MDTypography variant="body2" color="deepBlue">
                        This email is pending verification. Verification details
                        will appear once verification is completed.
                      </MDTypography>
                    </MDBox>
                  </Grid>
                )}

                {isCompleted && (
                  <>
                    <DetailItem
                      label="Provider"
                      value={leadDetails.emailProvIDer}
                    />

                    <DetailItem
                      label="Status"
                      value={renderStatusBadge(leadDetails.isSafe)}
                    />

                    <DetailItem
                      label="Reachable"
                      value={leadDetails.isReachable}
                    />

                    <DetailItem
                      label="Disposable"
                      value={leadDetails.isDisposable ? "Yes" : "No"}
                    />

                    <DetailItem
                      label="Role Account"
                      value={leadDetails.isRoleAccount ? "Yes" : "No"}
                    />

                    <DetailItem
                      label="Verified On"
                      value={
                        leadDetails.verifiedOn
                          ? formatDateTime(leadDetails.verifiedOn)
                          : "-"
                      }
                    />
                  </>
                )}

                {isFailed && (
                  <Grid item xs={12}>
                    <MDBox
                      display="flex"
                      alignItems="center"
                      gap={1}
                      px={1.5}
                      py={1.5}
                      sx={{
                        backgroundColor: "#fff8f8",
                        border: "0.063rem solid #ffd6d6",
                        borderRadius: "0.5rem",
                      }}
                    >
                      <Icon
                        sx={{
                          color: "error.main",
                          fontSize: "1.2rem",
                        }}
                      >
                        error_outline
                      </Icon>

                      <MDTypography variant="body2">
                        Email verification failed. You can reverify this email.
                      </MDTypography>
                    </MDBox>
                  </Grid>
                )}
              </>
            ) : (
              <>
                {/* Basic details */}
                <DetailItem label="First Name" value={leadDetails.firstName} />

                <DetailItem label="Last Name" value={leadDetails.lastName} />

                <DetailItem label="Email" value={leadDetails.email} />

                <DetailItem
                  label="Provider"
                  value={leadDetails.emailProvIDer}
                />

                <DetailItem
                  label="Priority"
                  value={renderPriorityBadge(
                    leadDetails.priority?.toLowerCase(),
                  )}
                />

                <DetailItem label="Company" value={leadDetails.company} />

                <DetailItem label="Job Title" value={leadDetails.jobTitle} />

                <DetailItem label="Industry" value={leadDetails.industry} />

                <DetailItem
                  label="Location"
                  value={
                    [leadDetails.city, leadDetails.country]
                      .filter(Boolean)
                      .join(", ") || "-"
                  }
                />

                {/* Verification */}
                <DetailItem
                  label="Verification"
                  value={leadDetails.verificationStatus}
                />

                {isPending && (
                  <Grid item xs={12}>
                    <MDBox
                      display="flex"
                      alignItems="center"
                      gap={1}
                      px={1.5}
                      py={1.5}
                      sx={{
                        backgroundColor: "#f5f8ff",
                        border: "0.063rem solid #dbe5ff",
                        borderRadius: "0.5rem",
                      }}
                    >
                      <Icon
                        sx={{
                          color: "deepBlue",
                          fontSize: "1.2rem",
                        }}
                      >
                        pending
                      </Icon>

                      <MDTypography variant="body2" color="deepBlue">
                        This email is pending verification. Verification details
                        will appear once verification is completed.
                      </MDTypography>
                    </MDBox>
                  </Grid>
                )}

                {isCompleted && (
                  <>
                    <DetailItem
                      label="Status"
                      value={renderStatusBadge(leadDetails.isSafe)}
                    />

                    <DetailItem
                      label="Reachable"
                      value={leadDetails.isReachable}
                    />

                    <DetailItem
                      label="Disposable"
                      value={leadDetails.isDisposable ? "Yes" : "No"}
                    />

                    <DetailItem
                      label="Role Account"
                      value={leadDetails.isRoleAccount ? "Yes" : "No"}
                    />

                    <DetailItem
                      label="Verified On"
                      value={
                        leadDetails.verifiedOn
                          ? formatDateTime(leadDetails.verifiedOn)
                          : "-"
                      }
                    />
                  </>
                )}

                {isFailed && (
                  <Grid item xs={12}>
                    <MDBox
                      display="flex"
                      alignItems="center"
                      gap={1}
                      px={1.5}
                      py={1.5}
                      sx={{
                        backgroundColor: "#fff8f8",
                        border: "0.063rem solid #ffd6d6",
                        borderRadius: "0.5rem",
                      }}
                    >
                      <Icon
                        sx={{
                          color: "error.main",
                          fontSize: "1.2rem",
                        }}
                      >
                        error_outline
                      </Icon>

                      <MDTypography variant="body2">
                        Email verification failed. You can reverify this email.
                      </MDTypography>
                    </MDBox>
                  </Grid>
                )}
              </>
            )}
          </Grid>
        ) : (
          <MDTypography align="center">No details found.</MDTypography>
        )}
      </MDBox>
    </CenterPopup>
  );
}

export default LeadDetailModal;
