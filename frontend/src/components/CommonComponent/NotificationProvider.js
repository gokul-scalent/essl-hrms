import React, { createContext, useContext, useState } from "react";
import PropTypes from "prop-types";
import Snackbar from "@mui/material/Snackbar";
import Alert from "@mui/material/Alert";
import Slide from "@mui/material/Slide";
import colors from "assets/theme/base/colors";

const NotificationContext = createContext();

function SlideTransition(props) {
  return (
    <Slide
      {...props}
      direction="down"
      timeout={{ enter: 600, exit: 400 }}
      easing={{
        enter: "cubic-bezier(0.4, 0, 0.2, 1)",
        exit: "cubic-bezier(0.4, 0, 1, 1)",
      }}
    />
  );
}

export const NotificationProvider = ({ children }) => {
  const [state, setState] = useState({
    open: false,
    message: "",
    severity: "success",
  });

  const notify = (message, severity = "success") => {
    setState({
      open: true,
      message,
      severity,
    });
  };

  const handleClose = (_, reason) => {
    if (reason === "clickaway") return;
    setState((prev) => ({ ...prev, open: false }));
  };

  return (
    <NotificationContext.Provider value={{ notify }}>
      {children}

      <Snackbar
        open={state.open}
        autoHideDuration={4000}
        onClose={handleClose}
        anchorOrigin={{ vertical: "top", horizontal: "right" }}
        TransitionComponent={SlideTransition}
        sx={{
          mt: 2,
        }}
      >
        <Alert
          severity={state.severity}
          variant="filled"
          onClose={handleClose}
          sx={{
            width: {
              xs: "90vw",
              sm: 420,
              md: 450,
            },
            maxWidth: "500px",
            fontWeight: 500,
            borderRadius: 3,
            px: 3,
            py: 1.5,
            fontSize: "0.92rem",
            color: colors.white.main,
            letterSpacing: 0.3,
          }}
        >
          {state.message}
        </Alert>
      </Snackbar>
    </NotificationContext.Provider>
  );
};

NotificationProvider.propTypes = {
  children: PropTypes.node.isRequired,
};

export const useNotify = () => useContext(NotificationContext);
