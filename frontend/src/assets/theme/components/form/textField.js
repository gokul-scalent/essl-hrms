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

// Material Dashboard 2 React Base Styles
import colors from "assets/theme/base/colors";

const { transparent } = colors;

const textField = {
  styleOverrides: {
    root: {
      backgroundColor: transparent.main,
      height: "2.5rem",

      "& input": {
        padding: "0 0.875rem",
        height: "2.7rem",
        boxSizing: "border-box",
      },

      "& .MuiOutlinedInput-root": {
        "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
          borderColor: colors.deepBlue.main,
        },
      },
    },
  },
};
export default textField;
