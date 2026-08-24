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
function collapseItem(theme, ownerState) {
  const { palette, transitions, breakpoints, borders, functions } = theme;
  const { active } = ownerState;

  const { grey, deepBlue, white } = palette;
  const { borderRadius } = borders;
  const { pxToRem } = functions;

  return {
    backgroundColor: active ? deepBlue.main : "transparent",
    color: active ? white.main : grey[700],
    display: "flex",
    alignItems: "center",
    width: "100%",
    padding: `${pxToRem(8)} ${pxToRem(10)}`,
    margin: `${pxToRem(4)} ${pxToRem(16)}`,
    borderRadius: borderRadius.md,
    cursor: "pointer",
    transition: transitions.create(["background-color", "color"], {
      easing: transitions.easing.easeInOut,
      duration: transitions.duration.shorter,
    }),

    "&:hover": {
      backgroundColor: active ? deepBlue.main : grey[200],
    },
  };
}

function collapseIconBox(theme, ownerState) {
  const { palette, borders, functions } = theme;
  const { active } = ownerState;

  const { grey, white } = palette;
  const { borderRadius } = borders;
  const { pxToRem } = functions;

  return {
    minWidth: pxToRem(32),
    minHeight: pxToRem(32),
    color: active ? white.main : grey[600],
    borderRadius: borderRadius.md,
    display: "grid",
    placeItems: "center",

    "& svg, svg g": {
      color: active ? white.main : grey[600],
    },
  };
}

const collapseIcon = ({ palette: { white, gradients } }, { active }) => ({
  color: active ? white.main : gradients.dark.state,
});

function collapseText(theme, ownerState) {
  const { typography, transitions, breakpoints, functions } = theme;
  const { miniSidenav, active } = ownerState;

  const { size, fontWeightRegular, fontWeightLight } = typography;
  const { pxToRem } = functions;

  return {
    marginLeft: pxToRem(10),

    [breakpoints.up("xl")]: {
      opacity: miniSidenav ? 0 : 1,
      maxWidth: miniSidenav ? 0 : "100%",
      marginLeft: miniSidenav ? 0 : pxToRem(10),
      transition: transitions.create(["opacity", "margin"], {
        easing: transitions.easing.easeInOut,
        duration: transitions.duration.standard,
      }),
    },

    "& span": {
      fontWeight: active ? fontWeightRegular : fontWeightLight,
      fontSize: size.sm,
      lineHeight: 0,
    },
  };
}
export { collapseItem, collapseIconBox, collapseIcon, collapseText };
