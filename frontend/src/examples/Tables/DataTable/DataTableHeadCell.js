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

// prop-types is a library for typechecking of props
import PropTypes from "prop-types";

// @mui material components
import Icon from "@mui/material/Icon";

// Material Dashboard 2 React components
import MDBox from "components/MDBox";

// Material Dashboard 2 React contexts
import { useMaterialUIController } from "context";

function DataTableHeadCell({ width, children, sorted, align, ...rest }) {
  const [controller] = useMaterialUIController();
  const { darkMode } = controller;

  return (
    <MDBox
      component="th"
      width={width}
      py={1.5}
      px={3}
      sx={({ palette: { light }, borders: { borderWidth } }) => ({
        minWidth: width,
        maxWidth: width,
        whiteSpace: "nowrap",
        borderBottom: `${borderWidth[1]} solid ${light.main}`,
        backgroundColor: "#eef3f8",
      })}
    >
      <MDBox
        {...rest}
        position="relative"
        textAlign={align}
        color={darkMode ? "white" : "primary"}
        opacity={0.8}
        sx={({ typography: { size, fontWeightBold } }) => ({
          fontSize: size.sm,
          fontWeight: fontWeightBold,
          textTransform: "uppercase",
          cursor: sorted && "pointer",
          userSelect: sorted && "none",
        })}
      >
        {children}
        {sorted && (
          <MDBox display="inline-flex" alignItems="left" ml={0.5}>
            {sorted === "asc" && (
              <Icon fontSize="small" color="action">
                arrow_drop_up
              </Icon>
            )}

            {sorted === "desc" && (
              <Icon fontSize="small" color="action">
                arrow_drop_down
              </Icon>
            )}

            {sorted === "none" && (
              <MDBox
                display="flex"
                flexDirection="column"
                sx={{ lineHeight: 0 }}
              >
                <Icon
                  fontSize="small"
                  color="disabled"
                  sx={{ marginBottom: "-6px" }}
                >
                  arrow_drop_up
                </Icon>
                <Icon
                  fontSize="small"
                  color="disabled"
                  sx={{ marginTop: "-6px" }}
                >
                  arrow_drop_down
                </Icon>
              </MDBox>
            )}
          </MDBox>
        )}
      </MDBox>
    </MDBox>
  );
}

// Setting default values for the props of DataTableHeadCell
DataTableHeadCell.defaultProps = {
  width: "auto",
  sorted: "none",
  align: "left",
};

// Typechecking props for the DataTableHeadCell
DataTableHeadCell.propTypes = {
  width: PropTypes.oneOfType([PropTypes.string, PropTypes.number]),
  children: PropTypes.node.isRequired,
  sorted: PropTypes.oneOf([false, "none", "asc", "desc"]),
  align: PropTypes.oneOf(["left", "right", "center"]),
};

export default DataTableHeadCell;
