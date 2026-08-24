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

// @mui material components
// @mui material components
import TextField from "@mui/material/TextField";
import { styled } from "@mui/material/styles";
import colors from "assets/theme/base/colors"; // your colors file

export default styled(TextField)(({ theme, ownerState }) => {
  const { palette, functions } = theme;
  const { error, success, disabled } = ownerState;

  const {
    grey,
    transparent,
    error: colorError,
    success: colorSuccess,
  } = palette;
  const { pxToRem } = functions;

  // styles for the input with error={true}
  const errorStyles = () => ({
    "& .MuiOutlinedInput-root": {
      "& .MuiOutlinedInput-notchedOutline": {
        borderColor: colorError.main,
        borderWidth: 2,
      },
    },
    "& .MuiOutlinedInput-root.Mui-focused .MuiOutlinedInput-notchedOutline": {
      borderColor: colorError.main,
      borderWidth: 2,
    },
    "& .MuiInputLabel-root.Mui-focused": {
      color: colorError.main,
    },
  });

  // styles for the input with success={true} – KEEP AS IS
  const successStyles = () => ({
    backgroundImage:
      "url(\"data:image/svg+xml;charset=utf-8,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 10 8'%3E%3Cpath fill='%234CAF50' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/%3E%3C/svg%3E\")",
    backgroundRepeat: "no-repeat",
    backgroundPosition: `right ${pxToRem(12)} center`,
    backgroundSize: `${pxToRem(16)} ${pxToRem(16)}`,
    "& .Mui-focused": {
      "& .MuiOutlinedInput-notchedOutline, &:after": {
        borderColor: colorSuccess.main,
      },
    },
    "& .MuiInputLabel-root.Mui-focused": {
      color: colorSuccess.main,
    },
  });

  // default and focused border in deepBlue ---
  const normalStyles = () => ({
    "& .MuiOutlinedInput-root": {
      borderRadius: 8,

      "& .MuiOutlinedInput-notchedOutline": {
        borderColor: colors.deepBlue.main,
        borderWidth: 1.5,
      },

      "&:hover .MuiOutlinedInput-notchedOutline": {
        borderColor: colors.deepBlue.main,
      },

      "&.Mui-focused .MuiOutlinedInput-notchedOutline": {
        borderColor: `${colors.deepBlue.main} !important`,
        borderWidth: 2,
      },
    },
    "& .MuiInputLabel-root.Mui-focused": {
      color: `${colors.deepBlue.main} !important`,
    },
  });

  return {
    backgroundColor: disabled ? `${grey[200]} !important` : transparent.main,
    pointerEvents: disabled ? "none" : "auto",
    ...normalStyles(),
    ...(error && errorStyles()),
    ...(success && successStyles()),
  };
});
