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

import React from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "App";
import { Provider } from "react-redux";
import "@fortawesome/fontawesome-free/css/all.min.css";
// Material Dashboard 2 React Context Provider
import { MaterialUIControllerProvider } from "context";
import { AuthProvider } from "pages/Auth/AuthContext";
import { reducers } from "reducers";
import { configureStore } from "@reduxjs/toolkit";

const container = document.getElementById("app");
const root = createRoot(container);
const store = configureStore({ reducer: reducers });
root.render(
  <Provider store={store}>
    <BrowserRouter>
      <AuthProvider>
        <MaterialUIControllerProvider>
          <App />
        </MaterialUIControllerProvider>
      </AuthProvider>
    </BrowserRouter>
  </Provider>,
);
