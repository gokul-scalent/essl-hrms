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

// react-router-dom components
import { Link, useNavigate } from "react-router-dom";

// @mui material components
import Card from "@mui/material/Card";
import Icon from "@mui/material/Icon";
import IconButton from "@mui/material/IconButton";
import { CircularProgress, InputAdornment } from "@mui/material";

// Material Dashboard 2 React components
import MDBox from "components/MDBox";
import MDTypography from "components/MDTypography";
import MDInput from "components/MDInput";
import MDButton from "components/MDButton";
import BasicLayout from "../components/BasicLayout";
import { useState } from "react";
import RequiredLabel from "components/CommonComponent/RequiredLabel";
import { REGEX, DEFAULT_ROUTES } from "components/common/constant";
import colors from "assets/theme/base/colors";
import { useNotify } from "components/CommonComponent/NotificationProvider";
import { encodePassword } from "components/CommonComponent/CommonFunction";
import { useAuth } from "pages/Auth/AuthContext";

// Authentication layout components

function Basic() {
  const navigate = useNavigate();
  const { notify } = useNotify();
  const { loginWithPassword } = useAuth();

  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);

  const [loginInfo, setLoginInfo] = useState({
    email: { value: "", error: "" },
    password: { value: "", error: "" },
  });

  const handleChange = (field, value) => {
    setLoginInfo((prev) => ({
      ...prev,
      [field]: { ...prev[field], value, error: "" },
    }));
  };

  const validateField = (field, value) => {
    switch (field) {
      case "email":
        if (value.trim() === "") {
          return "Email is required.";
        }
        if (!REGEX.EMAIL_REGEX.test(value.trim())) {
          return "Please enter a valid email address.";
        }
        return "";

      case "password":
        if (value.trim() === "") {
          return "Password is required.";
        }
        return "";

      default:
        return "";
    }
  };

  const validateAllFields = () => {
    const errors = {
      email: validateField("email", loginInfo.email.value),
      password: validateField("password", loginInfo.password.value),
    };

    setLoginInfo((prev) => ({
      email: { ...prev.email, error: errors.email },
      password: { ...prev.password, error: errors.password },
    }));

    return errors.email === "" && errors.password === "";
  };

  const handleBlur = (name) => {
    const value = loginInfo[name].value.trim();
    const error = validateField(name, value);

    setLoginInfo((prev) => ({
      ...prev,
      [name]: {
        ...prev[name],
        value,
        error,
      },
    }));
  };

  const handleLogin = async () => {
    if (!validateAllFields()) return;

    setLoading(true);

    try {
      const loginData = {
        email: loginInfo.email.value.trim(),
        password: encodePassword(loginInfo.password.value.trim()),

      };
      const res = await loginWithPassword(loginData);

      if (res.success) {
        notify("Login successful", "success");

        const role = res.data?.role;
        const redirectPath = DEFAULT_ROUTES[role] || "/dashboard";

        navigate(redirectPath, { replace: true });
      } else {
        const message = Array.isArray(res?.message)
          ? res.message
              .map((err) => (err.Msg ? err.Msg : JSON.stringify(err)))
              .join(", ")
          : res?.message || "Login Failed";
        notify(message, "error");
      }
    } catch (err) {
      notify(err?.message || "Something went wrong", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === "Enter") {
      handleLogin();
    }
  };
  return (
    <BasicLayout
      sx={{
        minHeight: "100vh",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        justifyContent: "center",
        pt: 6,
        background: "linear-gradient(180deg, #f5f6fa 0%, #e8edf7 100%)",
      }}
    >
      {/* Logo + Welcome text above the card */}
      <MDBox textAlign="center" mt={2} mb={2}>
        <MDBox mb={1}>
          {/* <img
            src={require("../../../assets/images/brand/rpg-logo.png")}
            alt="Mailora"
            style={{ width: 300, height: "auto" }}
          /> */}
        </MDBox>

        <MDTypography variant="h4" pt={4} fontWeight="bold" color="dark">
          Welcome!
        </MDTypography>

        {/* Smaller description text */}
        <MDTypography variant="body2" color="text.secondary" mt={1}>
          Scalent HRMS - Email Verification
        </MDTypography>
      </MDBox>

      <Card
        sx={{
          maxWidth: 550, // wider
          width: "100%",
          minHeight: 400, // taller
          borderRadius: 4,
          boxShadow: 15,
          overflow: "hidden",
        }}
      >
        {/* Header */}
        <MDBox textAlign="center" p={4}>
          <MDTypography
            variant="body2"
            color="text.secondary"
            mt={0.5}
            display="block"
          >
            Please login with your credentials
          </MDTypography>
        </MDBox>

        {/* Form */}
        <MDBox pt={2} pb={3} px={4} mt={0}>
          <MDBox component="form" role="form" onKeyDown={handleKeyPress}>
            {/* Email */}
            <MDBox mb={2}>
              <RequiredLabel>Email</RequiredLabel>
              <MDInput
                type="email"
                placeholder="Enter your email "
                fullWidth
                value={loginInfo.email.value}
                onChange={(e) => handleChange("email", e.target.value)}
                onBlur={() => handleBlur("email")}
                error={!!loginInfo.email.error}
                variant="outlined"
                inputProps={{ maxLength: 50 }}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <Icon
                        fontSize="small"
                        sx={{ color: colors.deepBlue.main }}
                      >
                        email
                      </Icon>
                    </InputAdornment>
                  ),
                  sx: { borderRadius: 2 },
                }}
              />
              {loginInfo.email.error && (
                <MDTypography variant="caption" color="error">
                  {loginInfo.email.error}
                </MDTypography>
              )}
            </MDBox>

            {/* Password */}
            <MDBox mb={1}>
              <RequiredLabel>Password</RequiredLabel>
              <MDInput
                type={showPassword ? "text" : "password"}
                placeholder="Please enter password"
                fullWidth
                value={loginInfo.password.value}
                onChange={(e) => handleChange("password", e.target.value)}
                onBlur={() => handleBlur("password")}
                error={!!loginInfo.password.error}
                variant="outlined"
                inputProps={{ maxLength: 16 }}
                InputProps={{
                  startAdornment: (
                    <InputAdornment position="start">
                      <Icon
                        fontSize="small"
                        sx={{ color: colors.deepBlue.main }}
                      >
                        lock
                      </Icon>
                    </InputAdornment>
                  ),
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        onClick={() => setShowPassword((prev) => !prev)}
                        edge="end"
                      >
                        <Icon fontSize="small">
                          {showPassword ? "visibility_off" : "visibility"}
                        </Icon>
                      </IconButton>
                    </InputAdornment>
                  ),
                  sx: { borderRadius: 2 },
                }}
              />
              {loginInfo.password.error && (
                <MDTypography variant="caption" color="error">
                  {loginInfo.password.error}
                </MDTypography>
              )}
            </MDBox>

            {/* <MDTypography
              variant="body2"
              sx={{
                display: "flex",
                justifyContent: "flex-end",
                fontSize: "0.75rem",
              }}
            >
              <Link
                // to="/auth/forget-password"
                style={{ fontWeight: "bold", color: colors.deepBlue.main }}
              >
                Forgot Password ?
              </Link>
            </MDTypography> */}

            {/* Login Button */}
            <MDBox mt={4} mb={1}>
              <MDButton
                variant="contained"
                fullWidth
                onClick={handleLogin}
                disabled={loading}
                color="deepBlue"
                sx={{
                  py: 1.75,
                  borderRadius: 3,
                  fontWeight: "bold",
                  fontSize: "0.95rem",
                  textTransform: "uppercase",
                }}
              >
                {loading ? (
                  <CircularProgress size={20} color="inherit" />
                ) : (
                  "Login"
                )}
              </MDButton>
            </MDBox>
          </MDBox>
        </MDBox>
      </Card>

      {/* Footer */}
      <MDBox
        textAlign="center"
        py={3}
        px={4}
        mt={5}
        mb={0}
        bgcolor="rgba(0,0,0,0.03)"
      >
        <MDTypography variant="body2" color="text.secondary">
          © 2026 Mailora. All rights reserved.
        </MDTypography>
      </MDBox>
    </BasicLayout>
  );
}

export default Basic;
