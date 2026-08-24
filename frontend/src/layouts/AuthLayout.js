import React, { useEffect } from "react";
import { Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "pages/Auth/AuthContext";
import routes from "routes";
import LoadingScreen from "./LoadingScreen";
import { DEFAULT_ROUTES } from "components/common/constant";

export default function AuthLayout() {
  const { role, isAuthenticated, loading } = useAuth();

  //   if (loading) {
  //     return <LoadingScreen />;
  //   }

  useEffect(() => {
    document.body.classList.add("bg-white");
    return () => document.body.classList.remove("bg-white");
  }, []);

  const defaultRoute = DEFAULT_ROUTES[role] || "/admin/dashboard"; //DEFAULT_ROUTES[role] ||

  const getAuthRoutes = (allRoutes) =>
    allRoutes.flatMap((route) => {
      if (route.collapse) {
        return getAuthRoutes(route.collapse);
      }

      if (route.layout === "/auth" && route.route) {
        return (
          <Route
            key={route.key}
            path={route.route.replace("/", "")}
            element={
              isAuthenticated ? (
                <Navigate to={defaultRoute} replace />
              ) : (
                route.component
              )
            }
          />
        );
      }

      return [];
    });

  return (
    <div className="auth-wrapper">
      <Routes>
        {getAuthRoutes(routes)}
        <Route path="*" element={<Navigate to="/auth/sign-in" replace />} />
      </Routes>
    </div>
  );
}
