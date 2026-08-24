/* eslint-disable react/prop-types */
import { login, logout } from "actions/auth";
import { ROLES } from "components/common/constant";
import React, { createContext, useContext, useEffect, useState } from "react";

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [auth, setAuth] = useState({
    email: null,
    role: null,
    token: null,
  });

  const [loading, setLoading] = useState(true);

  // Restore session
  useEffect(() => {
    const stored = localStorage.getItem("p");
    if (stored) {
      setAuth(JSON.parse(stored));
    }
    setLoading(false);
  }, []);

  const loginWithPassword = async (loginData) => {
    const res = await login(loginData);

    if (res.code === 200) {
      if (!Object.values(ROLES).includes(res.data?.role)) {
        return { success: false, message: "Access denied" };
      }

      setAuth(res.data);
      localStorage.setItem("p", JSON.stringify(res.data));

      return { success: true, data: res.data };
    }

    return { success: false, message: res?.message };
  };

  const logoutUser = async () => {
    try {
      await logout();
    } finally {
      localStorage.removeItem("p");
      setAuth({ email: null, token: null, role: null });
    }
  };

  const updateAuth = (data) => {
    setAuth(data);
    localStorage.setItem("p", JSON.stringify(data));
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
