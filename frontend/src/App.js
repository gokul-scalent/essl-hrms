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

import { useState, useEffect, useMemo } from "react";

// react-router components
import { Routes, Route, Navigate, useLocation } from "react-router-dom";

// @mui material components
import { ThemeProvider } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";

// Material Dashboard 2 React example components
import Sidenav from "examples/Sidenav";
import Configurator from "examples/Configurator";

// Material Dashboard 2 React themes
import theme from "assets/theme";

// Material Dashboard 2 React Dark Mode themes
import themeDark from "assets/theme-dark";

// Material Dashboard 2 React routes
import routes from "routes";

// Material Dashboard 2 React contexts
import { useMaterialUIController, setMiniSidenav } from "context";
import { NotificationProvider } from "components/CommonComponent/NotificationProvider";
import AuthLayout from "layouts/AuthLayout";
import AdminLayout from "layouts/AdminLayout";
import { DEFAULT_ROUTES } from "components/common/constant";
import logo from "assets/images/brand/mailora-1.png";

export default function App() {
  const [controller, dispatch] = useMaterialUIController();
  const { miniSidenav, layout, sidenavColor } = controller;
  const [onMouseEnter, setOnMouseEnter] = useState(false);
  const { pathname } = useLocation();
  const hideSidenav = pathname.startsWith("/auth");

  const authData = JSON.parse(localStorage.getItem("p")) || {};
  const token = authData?.token;
  const role = authData?.role;

  const isAuthenticated = () => !!token;

  const getDefaultRoute = () => {
    return DEFAULT_ROUTES[role] || "/admin/dashboard";
  };

  // Open sidenav when mouse enter on mini sidenav
  const handleOnMouseEnter = () => {
    if (miniSidenav && !onMouseEnter) {
      setMiniSidenav(dispatch, false);
      setOnMouseEnter(true);
    }
  };

  // Close sidenav when mouse leave mini sidenav
  const handleOnMouseLeave = () => {
    if (onMouseEnter) {
      setMiniSidenav(dispatch, true);
      setOnMouseEnter(false);
    }
  };

  // Setting page scroll to 0 when changing the route
  useEffect(() => {
    document.documentElement.scrollTop = 0;
    document.scrollingElement.scrollTop = 0;
  }, [pathname]);

  const getRoutes = (allRoutes) =>
    allRoutes.flatMap((route) => {
      if (route.collapse) {
        return getRoutes(route.collapse);
      }

      if (route.route && route.layout) {
        const isPublic = !route.roles || route.roles.length === 0;
        const hasAccess = Array.isArray(route.roles)
          ? route.roles.includes(role)
          : true;

        return (
          <Route
            key={route.key}
            path={`${route.layout}${route.route}`}
            element={
              isPublic ? (
                route.component
              ) : isAuthenticated() ? (
                hasAccess ? (
                  route.component
                ) : (
                  <Navigate to="/access-denied" replace />
                )
              ) : (
                <Navigate to="/auth/sign-in" replace />
              )
            }
          />
        );
      }

      return [];
    });
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      {layout === "dashboard" && !hideSidenav && (
        <>
          <Sidenav
            color={sidenavColor}
            brand={logo}
            brandName="Scalent HRMS"
            routes={routes}
            onMouseEnter={handleOnMouseEnter}
            onMouseLeave={handleOnMouseLeave}
          />
          <Configurator />
        </>
      )}
      {layout === "vr" && <Configurator />}

      <NotificationProvider>
        <Routes>
          <Route path="/admin/*" element={<AdminLayout />} />
          <Route path="/auth/*" element={<AuthLayout />} />
          <Route path="*" element={<Navigate to="/auth/sign-in" replace />} />
        </Routes>
      </NotificationProvider>
    </ThemeProvider>
  );
}
