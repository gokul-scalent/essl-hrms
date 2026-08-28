import React, { useEffect, useState } from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import Sidenav from "examples/Sidenav";
import Configurator from "examples/Configurator";
import routes from "routes";
import { useMaterialUIController, setMiniSidenav } from "context";
import { useAuth } from "pages/Auth/AuthContext";
import LoadingScreen from "./LoadingScreen";
import { DEFAULT_ROUTES } from "components/common/constant";
import logo from "assets/images/brand/mailora-1.png";
import ChangePasswordModal from "pages/ChangePassword/ChangePasswordModal";

export default function AdminLayout() {
  const [controller, dispatch] = useMaterialUIController();
  const { miniSidenav, sidenavColor } = controller;
  const [onMouseEnter, setOnMouseEnter] = useState(false);
  const [changePasswordOpen, setChangePasswordOpen] = useState(false);

  const { role, isAuthenticated, isPasswordSet, loading } = useAuth();

  useEffect(() => {
    document.documentElement.scrollTop = 0;
    document.scrollingElement.scrollTop = 0;
  }, []);

  useEffect(() => {
    if (isAuthenticated && isPasswordSet === "NO") {
      setChangePasswordOpen(true);
    }
  }, [isAuthenticated, isPasswordSet]);

  useEffect(() => {
    if (isAuthenticated && isPasswordSet === "NO") {
      setChangePasswordOpen(true);
    }
  }, [isAuthenticated, isPasswordSet]);
    const defaultRoute = DEFAULT_ROUTES[role] || "/admin/email-lists";

  if (loading) {
    return <LoadingScreen />;
  }

  // Expand sidenav on hover
  const handleOnMouseEnter = () => {
    if (miniSidenav && !onMouseEnter) {
      setMiniSidenav(dispatch, false);
      setOnMouseEnter(true);
    }
  };

  const handleOnMouseLeave = () => {
    if (onMouseEnter) {
      setMiniSidenav(dispatch, true);
      setOnMouseEnter(false);
    }
  };

  const getAdminRoutes = (allRoutes) =>
    allRoutes.flatMap((route) => {
      if (route.collapse) {
        return getAdminRoutes(route.collapse);
      }

      if (route.layout === "/admin" && route.route) {
        const hasAccess = route.roles?.includes(role);

        return (
          <Route
            key={route.key}
            path={route.route.replace("/", "")}
            element={
              isAuthenticated ? (
                hasAccess ? (
                  route.component
                ) : (
                  <Navigate to="dashboard" replace />
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

  if (!isAuthenticated) {
    return <Navigate to="/auth/sign-in" replace />;
  }

  const filterRoutesByRole = (allRoutes) =>
    allRoutes
      .map((route) => {
        if (route.collapse) {
          const filteredChildren = filterRoutesByRole(route.collapse);

          // keep parent ONLY if it has visible children
          if (filteredChildren.length > 0) {
            return { ...route, collapse: filteredChildren };
          }
          return null;
        }

        // normal route
        if (!route.roles || route.roles.includes(role)) {
          return route;
        }

        return null;
      })
      .filter(Boolean);
  return (
    <>
      <Sidenav
        color={sidenavColor}
        // brand={logo}
        routes={filterRoutesByRole(routes)}
        onMouseEnter={handleOnMouseEnter}
        onMouseLeave={handleOnMouseLeave}
        onChangePassword={() => setChangePasswordOpen(true)}
      />

      <ChangePasswordModal
        open={changePasswordOpen}
        handleClose={() => setChangePasswordOpen(false)}
      />
      <Configurator />

      <div className="main-content">
        <Routes>
          {getAdminRoutes(routes)}
          <Route
            path="*"
            element={<Navigate to={defaultRoute} replace />} //default route for now
          />
        </Routes>
      </div>
    </>
  );
}
