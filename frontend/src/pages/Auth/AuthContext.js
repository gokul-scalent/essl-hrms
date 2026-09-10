/* eslint-disable react/prop-types */
import { login, logout } from "actions/auth";
import { ROLES } from "components/common/constant";
import React, { createContext, useContext, useEffect, useState } from "react";

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [auth, setAuth] = useState({
    email: null,
    role: [],
    token: null,
    isPasswordSet: null,
  });

  const [loading, setLoading] = useState(true);

  // Restore session
  useEffect(() => {
    const stored = localStorage.getItem("p");
    if (stored) {
      try {
        const parsedAuth = JSON.parse(stored);

        setAuth({
          email: parsedAuth?.email || null,
          role: Array.isArray(parsedAuth?.role) ? parsedAuth.role : [],
          token: parsedAuth?.token || null,
          isPasswordSet: parsedAuth?.isPasswordSet || null,
        });
      } catch (error) {
        localStorage.removeItem("p");
      }
    }
    setLoading(false);
  }, []);

  const loginWithPassword = async (loginData) => {
    const res = await login(loginData);

    if (res.code === 200) {
      const roles = Array.isArray(res.data?.role) ? res.data.role : [];

      const hasValidRole = roles.some((role) =>
        Object.values(ROLES).includes(role),
      );

      if (!hasValidRole) {
        return { success: false, message: "Access denied" };
      }

      const authData = {
        ...res.data,
        role: roles,
      };

      setAuth(authData);
      localStorage.setItem("p", JSON.stringify(authData));

      return { success: true, data: authData };
    }

    return { success: false, message: res?.message || "Login failed" };
  };

  const logoutUser = async () => {
    try {
      await logout();
    } finally {
      localStorage.removeItem("p");

      setAuth({
        email: null,
        role: [],
        token: null,
        isPasswordSet: null,
      });
    }
  };

  const updateAuth = (data) => {
    const updatedAuth = {
      ...data,
      role: Array.isArray(data?.role) ? data.role : [],
    };

    setAuth(updatedAuth);
    localStorage.setItem("p", JSON.stringify(updatedAuth));
  };

  return (
    <AuthContext.Provider
      value={{
        ...auth,
        isAuthenticated: !!auth?.token,
        loading,
        loginWithPassword,
        logoutUser,
        updateAuth,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
