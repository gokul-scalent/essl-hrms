import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  InputAdornment,
  Icon,
  LinearProgress,
  Divider,
  Grid,
} from "@mui/material";

import MDBox from "components/MDBox";
import MDInput from "components/MDInput";
import MDButton from "components/MDButton";
import MDTypography from "components/MDTypography";
import RequiredLabel from "components/CommonComponent/RequiredLabel";
import colors from "assets/theme/base/colors";
import { useNotify } from "components/CommonComponent/NotificationProvider";
import { changePasswordAPI } from "actions/password";
import { encodePassword } from "components/CommonComponent/CommonFunction";
import { useAuth } from "pages/Auth/AuthContext";

const initialState = {
  currentPassword: { value: "", error: "" },
  newPassword: { value: "", error: "" },
  confirmPassword: { value: "", error: "" },
};

function ChangePasswordModal({ open, handleClose }) {
  const { updateAuth, ...auth } = useAuth();
  const { notify } = useNotify();
  const [passwordForm, setPasswordForm] = useState(initialState);
  const [loading, setLoading] = useState(false);
  const [show, setShow] = useState({
    current: false,
    new: false,
    confirm: false,
  });

  useEffect(() => {
    if (!open) {
      setPasswordForm(initialState);
      setLoading(false);
      setShow({ current: false, new: false, confirm: false });
    }
  }, [open]);

  const getPasswordStrength = (password) => {
    if (!password) {
      return { score: 0, label: "Weak", error: "Password is required" };
    }

    if (password.length < 12) {
      return {
        score: 0,
        label: "Weak",
        error: "Password must be at least 12 characters long",
      };
    }

    const hasUpper = /[A-Z]/.test(password);
    const hasLower = /[a-z]/.test(password);
    const hasNumber = /[0-9]/.test(password);
    const hasSpecial = /[^A-Za-z0-9]/.test(password);

    if (!hasUpper) {
      return {
        score: 1,
        label: "Weak",
        error: "Must contain at least 1 uppercase letter",
      };
    }

    if (!hasLower) {
      return {
        score: 1,
        label: "Weak",
        error: "Must contain at least 1 lowercase letter",
      };
    }

    if (!hasNumber) {
      return {
        score: 1,
        label: "Weak",
        error: "Must contain at least 1 number",
      };
    }

    if (!hasSpecial) {
      return {
        score: 1,
        label: "Weak",
        error: "Must contain at least 1 special character",
      };
    }

    // scoring (safe + consistent)
    let score = 0;

    if (password.length >= 12) score++;
    if (password.length >= 16) score++;
    if (hasUpper) score++;
    if (hasLower) score++;
    if (hasNumber) score++;
    if (hasSpecial) score++;

    score = Math.min(score, 5);

    if (score <= 2) return { score, label: "Weak", error: "" };
    if (score <= 4) return { score, label: "Medium", error: "" };

    return { score, label: "Strong", error: "" };
  };

  const getStrengthValue = (strength) => {
    if (strength?.label === "Strong") return 100;
    if (strength?.label === "Medium") return 60;
    return 30;
  };

  const getStrengthColor = (strength) => {
    if (strength?.label === "Strong") return "success";
    if (strength?.label === "Medium") return "warning";
    return "error";
  };

  const toggleVisibility = (key) =>
    setShow((prev) => ({ ...prev, [key]: !prev[key] }));

  const handleChange = (name, value) => {
    setPasswordForm((prev) => ({
      ...prev,
      [name]: { value, error: "" },
    }));
  };

  const getValidationError = (name, value) => {
    switch (name) {
      case "currentPassword":
        return !value ? "Please enter current password" : "";
      case "newPassword":
        if (!value) return "Please enter a new password";

        const strength = getPasswordStrength(value);
        if (strength.error) return strength.error;

        if (
          passwordForm.confirmPassword.value &&
          value !== passwordForm.confirmPassword.value
        ) {
          return "The new password and confirm password must match";
        }

        return "";
      case "confirmPassword":
        if (!value) return "Please confirm password";

        if (value !== passwordForm.newPassword.value) {
          return "The new password and confirm password must match";
        }

        return "";
      default:
        return "";
    }
  };

  const handleOnBlur = (name, value) => {
    setPasswordForm((prev) => ({
      currentPassword: { ...prev.currentPassword, error: "" },
      newPassword: { ...prev.newPassword, error: "" },
      confirmPassword: { ...prev.confirmPassword, error: "" },
    }));
    const error = getValidationError(name, value);
    setPasswordForm((prev) => ({
      ...prev,
      [name]: { ...prev[name], error },
    }));
  };

  const validateForm = () => {
    let temp = { ...passwordForm };
    Object.keys(temp).forEach((key) => {
      temp[key] = {
        ...temp[key],
        error: getValidationError(key, temp[key].value),
      };
    });
    const newPassword = temp.newPassword.value;
    const strength = getPasswordStrength(newPassword);

    if (strength.label !== "Strong") {
      temp.newPassword.error = "Only STRONG password is allowed";
    }
    setPasswordForm(temp);
    return Object.values(temp).every((f) => f.error === "");
  };

  const handleSubmit = async () => {
    if (!validateForm()) return;
    try {
      setLoading(true);

      const payload = {
        oldPassword: encodePassword(passwordForm.currentPassword.value),
        newPassword: encodePassword(passwordForm.newPassword.value),
        confirmPassword: encodePassword(passwordForm.confirmPassword.value),
      };

      const res = await changePasswordAPI(payload);

      if (res?.code === 200) {
        notify(res?.message || "Password updated successfully", "success");
        //after succesfully update or change password -> set yes
        updateAuth({
          ...auth,
          isPasswordSet: "YES",
        });
        // refresh page to force logout/session reset
        window.location.reload();

        handleClose();
      } else {
        const message = Array.isArray(res?.message)
          ? res?.message[0]?.Msg || JSON.stringify(res?.message)
          : res?.message || "Failed to update password";
        notify(res?.message || "Failed to update password", "error");
      }
    } catch (error) {
      notify(
        error?.message || "Something went wrong while updating password",
        "error",
      );
    } finally {
      setLoading(false);
    }
  };

  const strength = passwordForm.newPassword.value
    ? getPasswordStrength(passwordForm.newPassword.value)
    : null;

  const canCloseDialog = () => {
    if (auth?.isPasswordSet === "NO") {
      notify("Please set your password first", "info");
      return false;
    }
    return true;
  };
  return (
    <Dialog open={open} onClose={canCloseDialog} maxWidth="xs" fullWidth>
      <DialogTitle sx={{ position: "relative" }}>
        <IconButton
          onClick={() => {
            if (!canCloseDialog()) return;
            handleClose();
          }}
          sx={{
            position: "absolute",
            right: 8,
            top: 8,
            color: (theme) => theme.palette.grey[500],
          }}
        >
          <Icon>close</Icon>
        </IconButton>

        <MDBox textAlign="center">
          <Icon fontSize="large" sx={{ color: colors.deepBlue.main }}>
            lock
          </Icon>
          <MDTypography variant="h5" mt={1}>
            Change Password
          </MDTypography>
          <MDTypography variant="caption" color="text">
            Choose a strong password to keep your account secure
          </MDTypography>
        </MDBox>
      </DialogTitle>

      {/* <Divider /> */}

      <DialogContent>
        <Grid container spacing={2}>
          {/* Current Password */}
          <Grid item xs={12}>
            <RequiredLabel required>Current Password</RequiredLabel>
            <MDInput
              variant="outlined"
              placeholder="Please enter current password"
              fullWidth
              name="currentPassword"
              type={show.current ? "text" : "password"}
              value={passwordForm.currentPassword.value}
              onChange={(e) => handleChange(e.target.name, e.target.value)}
              onBlur={(e) => handleOnBlur(e.target.name, e.target.value)}
              error={!!passwordForm.currentPassword.error}
              inputProps={{ maxLength: 16 }}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => toggleVisibility("current")}
                      edge="end"
                    >
                      <Icon>
                        {show.current ? "visibility_off" : "visibility"}
                      </Icon>
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
            {passwordForm.currentPassword.error && (
              <MDTypography variant="caption" color="error">
                {passwordForm.currentPassword.error}
              </MDTypography>
            )}
          </Grid>

          {/* New Password */}
          <Grid item xs={12}>
            <RequiredLabel required>New Password</RequiredLabel>
            <MDInput
              variant="outlined"
              placeholder="Please enter new password"
              fullWidth
              name="newPassword"
              type={show.new ? "text" : "password"}
              value={passwordForm.newPassword.value}
              onChange={(e) => handleChange(e.target.name, e.target.value)}
              onBlur={(e) => handleOnBlur(e.target.name, e.target.value)}
              error={!!passwordForm.newPassword.error}
              inputProps={{ maxLength: 16 }}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => toggleVisibility("new")}
                      edge="end"
                    >
                      <Icon>{show.new ? "visibility_off" : "visibility"}</Icon>
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
            {passwordForm.newPassword.error && (
              <MDTypography variant="caption" color="error">
                {passwordForm.newPassword.error}
              </MDTypography>
            )}

            {/* Strength bar */}
            {passwordForm.newPassword.value && (
              <MDBox mt={1}>
                <LinearProgress
                  variant="determinate"
                  value={getStrengthValue(strength)}
                  color={getStrengthColor(strength)}
                  sx={{ height: 6, borderRadius: 5 }}
                />
                <MDTypography
                  variant="caption"
                  color={getStrengthColor(strength)}
                >
                  {strength.label} password
                </MDTypography>
              </MDBox>
            )}
          </Grid>

          {/* Confirm Password */}
          <Grid item xs={12}>
            <RequiredLabel required>Confirm Password</RequiredLabel>
            <MDInput
              variant="outlined"
              fullWidth
              placeholder="Please confirm new password"
              name="confirmPassword"
              type={show.confirm ? "text" : "password"}
              value={passwordForm.confirmPassword.value}
              onChange={(e) => handleChange(e.target.name, e.target.value)}
              onBlur={(e) => handleOnBlur(e.target.name, e.target.value)}
              error={!!passwordForm.confirmPassword.error}
              inputProps={{ maxLength: 16 }}
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => toggleVisibility("confirm")}
                      edge="end"
                    >
                      <Icon>
                        {show.confirm ? "visibility_off" : "visibility"}
                      </Icon>
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
            {passwordForm.confirmPassword.error && (
              <MDTypography variant="caption" color="error">
                {passwordForm.confirmPassword.error}
              </MDTypography>
            )}
          </Grid>

          {/* Security note */}
          <Grid item xs={12}>
            <MDTypography
              variant="caption"
              color="text"
              display="flex"
              justifyContent="center"
              alignItems="center"
            >
              <Icon fontSize="small" sx={{ mr: 0.5 }}>
                security
              </Icon>
              Your password is encrypted and securely stored
            </MDTypography>
          </Grid>
        </Grid>
      </DialogContent>

      <DialogActions sx={{ p: 2, gap: 1 }}>
        <MDButton
          variant="gradient"
          color="light"
          onClick={() => {
            if (!canCloseDialog()) return;
            handleClose();
          }}
          disabled={loading}
        >
          Cancel
        </MDButton>

        <MDButton
          variant="contained"
          color="deepBlue"
          onClick={handleSubmit}
          disabled={loading}
        >
          {loading ? (
            <>
              <i className="fa fa-spinner fa-spin fa-lg text-white" />
            </>
          ) : (
            "Save"
          )}
        </MDButton>
      </DialogActions>
    </Dialog>
  );
}

ChangePasswordModal.propTypes = {
  open: PropTypes.bool.isRequired,
  handleClose: PropTypes.func.isRequired,
};

export default ChangePasswordModal;
