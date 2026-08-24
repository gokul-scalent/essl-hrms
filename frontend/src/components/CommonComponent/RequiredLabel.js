import React from "react";
import PropTypes from "prop-types";
import MDTypography from "components/MDTypography";

const RequiredLabel = ({ children, required = false, ...props }) => (
  <MDTypography
    {...props}
    sx={{ display: "flex", alignItems: "center", mb: 0.5 , fontSize : "0.9rem" }}
  >
    {children}
    {required && <span style={{ color: "red", marginLeft: 2 }}>*</span>}
  </MDTypography>
);

RequiredLabel.propTypes = {
  children: PropTypes.node.isRequired,
  required: PropTypes.bool,
};

export default RequiredLabel;
