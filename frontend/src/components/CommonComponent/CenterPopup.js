import React, { useState, useEffect, useRef } from "react";
import PropTypes from "prop-types";
import { Backdrop, Paper } from "@mui/material";

const CenterPopup = ({ isOpen, onClose, children, width = 450 }) => {
  const [visible, setVisible] = useState(false);
  const [active, setActive] = useState(false);
  const popupRef = useRef(null);

  useEffect(() => {
    if (isOpen) {
      setVisible(true);
      setTimeout(() => setActive(true), 10);
    } else {
      setActive(false);
      const timeout = setTimeout(() => setVisible(false), 300);
      return () => clearTimeout(timeout);
    }
  }, [isOpen]);

  const handleClose = () => {
    setActive(false);
    setTimeout(() => onClose && onClose(), 300);
  };

  if (!visible) return null;

  return (
    <>
      <Backdrop
        open={active}
        onClick={handleClose}
        sx={{
          zIndex: (theme) => theme.zIndex.drawer + 1,
          transition: "opacity 0.3s",
        }}
      />

      <Paper
        ref={popupRef}
        elevation={8}
        sx={{
          position: "fixed",
          top: "50%",
          left: "50%",
          transform: active
            ? "translate(-50%, -50%) scale(1)"
            : "translate(-50%, -50%) scale(0.8)",
          opacity: active ? 1 : 0,
          width: width, // adjust as needed
          maxWidth: "90vw",
          maxHeight: "90vh",
          overflowY: "auto",
          transition: "opacity 0.3s, transform 0.3s",
          p: 3,
          borderRadius: 2,
          display: "flex",
          flexDirection: "column",
          lineHeight: "1",
          zIndex: (theme) => theme.zIndex.drawer + 2,
        }}
      >
        {children}
      </Paper>
    </>
  );
};

CenterPopup.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  children: PropTypes.node,
  width: PropTypes.oneOfType([PropTypes.number, PropTypes.string]),
};

export default CenterPopup;
