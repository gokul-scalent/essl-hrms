import React, { useEffect, useState } from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import Sidenav from "examples/Sidenav";
import Configurator from "examples/Configurator";
import routes from "routes";
import { useMaterialUIController, setMiniSidenav } from "context";
import { useAuth } from "pages/Auth/AuthContext";
import LoadingScreen from "./LoadingScreen";
import { DEFAULT_ROUTES } from "components/common/constant";
import ChangePasswordModal from "pages/ChangePassword/ChangePasswordModal";

export default function AdminLayout() {
  const [controller, dispatch] = useMaterialUIController();
  const { miniSidenav, sidenavColor } = controller;
  const [onMouseEnter, setOnMouseEnter] = useState(false);
  const [changePasswordOpen, setChangePasswordOpen] = useState(false);

  const {
    role,
    isAuthenticated,
    isPasswordSet,
    loading,
  } = useAuth();

  // Always treat role as an array
  const roles = Array.isArray(role) ? role : [];

  // Scroll to top when layout loads
  useEffect(() => {
    document.documentElement.scrollTop = 0;
    document.scrollingElement.scrollTop = 0;
  }, []);

  // Open change password modal when password is not set
  useEffect(() => {
    if (isAuthenticated && isPasswordSet === "NO") {
      setChangePasswordOpen(true);
    }
  }, [isAuthenticated, isPasswordSet]);

  if (loading) {
    return <LoadingScreen />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/auth/sign-in" replace />;
  }

  // Expand sidenav on hover
  const handleOnMouseEnter = () => {
    if (miniSidenav && !onMouseEnter) {
      setMiniSidenav(dispatch, false);
      setOnMouseEnter(true);
    }
  };

  // Collapse sidenav when mouse leaves
  const handleOnMouseLeave = () => {
    if (onMouseEnter) {
      setMiniSidenav(dispatch, true);
      setOnMouseEnter(false);
    }
  };

  // Get first valid role for default route
  const userRole = roles.find((currentRole) =>
    Object.prototype.hasOwnProperty.call(DEFAULT_ROUTES, currentRole),
  );

  const defaultRoute =
    DEFAULT_ROUTES[userRole] || "/admin/dashboard";

  // Check whether user has access to route
  const hasRoleAccess = (routeRoles) => {
    // Public route
    if (!routeRoles || routeRoles.length === 0) {
      return true;
    }

    // User must have at least one matching role
    return routeRoles.some((routeRole) =>
      roles.includes(routeRole),
    );
  };

  // Generate admin routes
  const getAdminRoutes = (allRoutes) =>
    allRoutes.flatMap((route) => {
      if (route.collapse) {
        return getAdminRoutes(route.collapse);
      }

      if (route.layout === "/admin" && route.route) {
        const hasAccess = hasRoleAccess(route.roles);

        return (
          <Route
            key={route.key}
            path={route.route.replace("/", "")}
            element={
              hasAccess ? (
                route.component
              ) : (
                <Navigate to="/admin/access-denied" replace />
              )
            }
          />
        );
      }

      return [];
    });

  // Filter sidenav routes according to user roles
  const filterRoutesByRole = (allRoutes) =>
    allRoutes
      .map((route) => {
        // Parent/collapse route
        if (route.collapse) {
          const filteredChildren = filterRoutesByRole(route.collapse);

          if (filteredChildren.length > 0) {
            return {
              ...route,
              collapse: filteredChildren,
            };
          }
          return null;
        }

        // Route without role restriction
        if (!route.roles || route.roles.length === 0) {
          return route;
        }

        // Route with role restriction
        if (hasRoleAccess(route.roles)) {
          return route;
        }

        return null;
      })
      .filter(Boolean);
  const filteredRoutes = filterRoutesByRole(routes);
  return (
    <>
      <Sidenav
        color={sidenavColor}
        routes={filteredRoutes}
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
            element={<Navigate to={defaultRoute} replace />}
          />
        </Routes>
      </div>
    </>
  );
}