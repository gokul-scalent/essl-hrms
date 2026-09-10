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

// Material Dashboard 2 React base styles
import boxShadows from "assets/theme/base/boxShadows";
import typography from "assets/theme/base/typography";
import colors from "assets/theme/base/colors";
import borders from "assets/theme/base/borders";

// Material Dashboard 2 React helper functions
import pxToRem from "assets/theme/functions/pxToRem";

const { lg } = boxShadows;
const { size } = typography;
const { text, white, transparent, light, dark, gradients } = colors;
const { borderRadius } = borders;

const autocomplete = {
  styleOverrides: {
    root: {
      minHeight: "2.73rem",

      "& .MuiFormControl-root": {
        height: "auto",
      },
    },

    inputRoot: {
      minHeight: "2.73rem",
      padding: "0 !important",
      display: "flex",
      alignItems: "center",
      flexWrap: "wrap",

      "& .MuiOutlinedInput-notchedOutline": {
        borderRadius: "0.4285rem",
        borderColor: "#D0D7DE",
      },

      // "&:hover .MuiOutlinedInput-notchedOutline": {
      //   borderColor: colors.deepBlue.main,
      // },

      "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
        borderColor: colors.deepBlue.main,
      },

      "&.Mui-error .MuiOutlinedInput-notchedOutline": {
        borderColor: colors.vividRed.main,
        borderWidth: 2,
      },

      "&.Mui-error.Mui-focused .MuiOutlinedInput-notchedOutline": {
        borderColor: colors.vividRed.main,
        borderWidth: 2,
      },

      "& input": {
        padding: "0 0.75rem !important",
        height: "auto",
        boxSizing: "border-box",
        fontSize: "0.75rem",
      },

      "& .MuiAutocomplete-tag": {
        margin: "2px", // better chip spacing
      },
    },

    endAdornment: {
      right: "0.5rem",
    },

    popupIndicator: {
      "& svg": {
        fontSize: "1.5rem", // try 1.5rem or 1.6rem
      },
      color: "#6C757D",
    },

    clearIndicator: {
      transform: "scale(1.1)",
    },

    paper: {
      boxShadow: "none",
      margin: 0,
    },

    option: {
      padding: "6px 16px",
      fontSize: "0.75rem",
      colors: colors.text.main,
    },
  },
};

export default autocomplete;
