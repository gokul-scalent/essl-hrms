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

import { useEffect, useRef, useState } from "react";
import { useLocation, NavLink } from "react-router-dom";
import PropTypes from "prop-types";

// @mui material components
import List from "@mui/material/List";
import Divider from "@mui/material/Divider";
import Link from "@mui/material/Link";
import Icon from "@mui/material/Icon";

// Material Dashboard components
import MDBox from "components/MDBox";
import MDTypography from "components/MDTypography";

// Example components
import SidenavCollapse from "examples/Sidenav/SidenavCollapse";
import SidenavRoot from "examples/Sidenav/SidenavRoot";
import sidenavLogoLabel from "examples/Sidenav/styles/sidenav";

// Context
import { useMaterialUIController, setMiniSidenav } from "context";
import { useAuth } from "pages/Auth/AuthContext";

import {
  Avatar,
  Menu,
  MenuItem,
  IconButton,
  ListItemIcon,
  ListItemText,
  Popper,
  ClickAwayListener,
  Paper,
} from "@mui/material";

import MoreVertIcon from "@mui/icons-material/MoreVert";
import SettingsIcon from "@mui/icons-material/Settings";
import LogoutIcon from "@mui/icons-material/Logout";
function Sidenav({ brand, brandName, routes, onChangePassword, ...rest }) {
  // Access global UI controller state (Material Dashboard context)
  const [controller, dispatch] = useMaterialUIController();
  const { miniSidenav } = controller;

  // Get current URL location
  const location = useLocation();
  const { logoutUser } = useAuth();
  const profileRef = useRef();
  const [anchorEl, setAnchorEl] = useState(null);
  const authData = JSON.parse(localStorage.getItem("p")) || {};

  const userEmail = authData?.email;

  const open = Boolean(anchorEl);

  const handleMenuOpen = (e) => {
    setAnchorEl(e.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  // Extract route name from URL (currently unused but kept if needed later)
  const collapseName = location.pathname.replace("/", "");

  // Close sidenav manually (used in mobile view)
  const closeSidenav = () => setMiniSidenav(dispatch, true);

  // Handle responsive behaviour of sidenav
  useEffect(() => {
    function handleMiniSidenav() {
      // Automatically collapse sidenav if screen width < 1200px
      setMiniSidenav(dispatch, window.innerWidth < 1200);
    }

    window.addEventListener("resize", handleMiniSidenav);
    handleMiniSidenav();

    return () => window.removeEventListener("resize", handleMiniSidenav);
  }, [dispatch]);

  //check child nested route is active or not
  // Check if any child (or nested child) route is active
  const isChildActive = (collapseRoutes) => {
    return collapseRoutes.some((route) => {
      // If nested collapse exists, check recursively
      if (route.collapse) {
        return isChildActive(route.collapse);
      }

      const fullPath = `${route.layout || ""}${route.route}`;

      // Open parent if:
      //Exact match / current path starts with child path (for nested routes)
      return (
        location.pathname === fullPath ||
        location.pathname.startsWith(fullPath + "/")
      );
    });
  };
  const renderRoutes = (routes) =>
    routes
      .filter((route) => route.show !== false) // Skip hidden routes
      .map((route) => {
        const { type, name, icon, key, route: path, collapse, layout } = route;

        if (type === "collapse") {
          // Parent collapse (has children)
          if (collapse) {
            // Check if any child route is active
            const open = isChildActive(collapse);

            return (
              <SidenavCollapse
                key={key}
                name={name}
                icon={icon}
                // active={open}
                open={open}

                //"active={open}"because that would apply deepBlue color
              >
                {renderRoutes(collapse)}
              </SidenavCollapse>
            );
          }

          // Single clickable route
          return (
            <NavLink
              key={key}
              to={`${layout || ""}${path}`}
              style={{ textDecoration: "none" }}
            >
              {({ isActive }) => (
                <SidenavCollapse name={name} icon={icon} active={isActive} />
              )}
            </NavLink>
          );
        }

        return null;
      });
  return (
    <SidenavRoot {...rest} variant="permanent" ownerState={{ miniSidenav }}>
      <MDBox pt={3} pb={1} px={4} textAlign="center">
        {/* Close button (visible only on small screens) */}
        <MDBox
          display={{ xs: "block", xl: "none" }}
          position="absolute"
          top={0}
          right={0}
          p={1.625}
          onClick={closeSidenav}
          sx={{ cursor: "pointer" }}
        >
          <MDTypography variant="h6" color="secondary">
            <Icon sx={{ fontWeight: "bold" }}>close</Icon>
          </MDTypography>
        </MDBox>

        {/* Logo + Brand Name */}
        <MDBox display="flex" alignItems="center">
          {brand && (
            <MDBox component="img" src={brand} alt="Brand" width="10rem" />
          )}

          {!miniSidenav && brandName && (
            <MDBox
              sx={{
                ml: 1,
                display: "flex",
                alignItems: "center",
              }}
            >
              <MDTypography
                component="h6"
                fontWeight="bold"
                sx={(theme) => ({
                  color: theme.palette.deepBlue.main,
                  fontSize: theme.typography.pxToRem(20),
                  letterSpacing: "0.5px",
                })}
              >
                {brandName}
              </MDTypography>
            </MDBox>
          )}
        </MDBox>
      </MDBox>

      {/* ---------------- Navigation List ---------------- */}
      <MDBox
        sx={{
          flex: 1,
          overflowY: "auto",
        }}
      >
        <List>{renderRoutes(routes)}</List>
      </MDBox>

      <MDBox
        ref={profileRef}
        p={2}
        display="flex"
        alignItems="center"
        justifyContent="space-between"
        sx={{
          borderTop: "0.063rem solid #ececec",
          position: "relative",
        }}
      >
        <MDBox display="flex" alignItems="center">
          <Avatar
            sx={{
              bgcolor: "deepBlue.main",
              width: 25,
              height: 25,
            }}
          ></Avatar>

          {!miniSidenav && (
            <MDBox ml={1.5}>
              <MDTypography
                variant="button"
                fontWeight="medium"
                sx={{
                  maxWidth: 170,
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                  whiteSpace: "nowrap",
                }}
              >
                {userEmail || "Administrator"}
              </MDTypography>
            </MDBox>
          )}
        </MDBox>

        {!miniSidenav && (
          <IconButton onClick={handleMenuOpen}>
            <MoreVertIcon />
          </IconButton>
        )}
      </MDBox>

      <Popper
        open={open}
        anchorEl={profileRef.current}
        placement="top"
        style={{
          width: profileRef.current?.offsetWidth,
          zIndex: 1300,
        }}
      >
        <ClickAwayListener onClickAway={handleMenuClose}>
          <Paper
            elevation={8}
            sx={{
              p: 1,
              mb: 1,
              borderRadius: 2,
              overflow: "hidden",
            }}
          >
            <MenuItem
              onClick={() => {
                handleMenuClose();
                onChangePassword();
              }}
            >
              <ListItemIcon sx={{ minWidth: 32 }}>
                <SettingsIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText
                primary="Change Password"
                primaryTypographyProps={{
                  fontSize: "0.875rem",
                  fontWeight: 500,
                }}
              />
            </MenuItem>
            <Divider sx={{ my: 0.5 }} />

            <MenuItem
              onClick={() => {
                handleMenuClose();
                logoutUser();
              }}
            >
              <ListItemIcon sx={{ minWidth: 32 }}>
                <LogoutIcon fontSize="small" />
              </ListItemIcon>

              <ListItemText
                primary="Logout"
                primaryTypographyProps={{
                  fontSize: "0.875rem",
                  fontWeight: 500,
                }}
              />
            </MenuItem>
          </Paper>
        </ClickAwayListener>
      </Popper>
    </SidenavRoot>
  );
}

Sidenav.propTypes = {
  brand: PropTypes.string,
  brandName: PropTypes.string.isRequired,
  routes: PropTypes.arrayOf(PropTypes.object).isRequired,
  onChangePassword: PropTypes.func.isRequired,
};

Sidenav.defaultProps = {
  brand: "",
};

export default Sidenav;
