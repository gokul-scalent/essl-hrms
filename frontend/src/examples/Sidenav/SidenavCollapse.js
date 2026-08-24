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

import { useState, useEffect } from "react";
import PropTypes from "prop-types";

import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Collapse from "@mui/material/Collapse";
import Icon from "@mui/material/Icon";
import ExpandLess from "@mui/icons-material/ExpandLess";
import ExpandMore from "@mui/icons-material/ExpandMore";

import MDBox from "components/MDBox";

import {
  collapseItem,
  collapseIconBox,
  collapseIcon,
  collapseText,
} from "examples/Sidenav/styles/sidenavCollapse";

import { useMaterialUIController } from "context";

function SidenavCollapse({ icon, name, active, children, open: openProp }) {
  const [controller] = useMaterialUIController();
  const { miniSidenav } = controller;

  // internal state
  const [open, setOpen] = useState(openProp || false);

  // when any of child of collapse parent are active stay open it
  useEffect(() => {
    if (typeof openProp === "boolean") {
      setOpen(openProp);
    }
  }, [openProp]);

  const handleToggle = () => {
    if (children) {
      setOpen((prev) => !prev);
    }
  };

  const hasChildren = Boolean(children);

  return (
    <>
      <ListItem component="li" onClick={hasChildren ? handleToggle : undefined}>
        <MDBox
          sx={(theme) =>
            collapseItem(theme, {
              active,
              whiteSidenav: true,
            })
          }
        >
          <ListItemIcon
            sx={(theme) =>
              collapseIconBox(theme, { whiteSidenav: true, active })
            }
          >
            {typeof icon === "string" ? (
              <Icon sx={(theme) => collapseIcon(theme, { active })}>
                {icon}
              </Icon>
            ) : (
              icon
            )}
          </ListItemIcon>

          <ListItemText
            primary={name}
            sx={(theme) =>
              collapseText(theme, {
                miniSidenav,
                whiteSidenav: true,
                active,
              })
            }
          />

          {hasChildren && (open ? <ExpandLess /> : <ExpandMore />)}
        </MDBox>
      </ListItem>

      {hasChildren && (
        <Collapse in={open} timeout="auto" unmountOnExit>
          {children}
        </Collapse>
      )}
    </>
  );
}

SidenavCollapse.propTypes = {
  icon: PropTypes.node.isRequired,
  name: PropTypes.string.isRequired,
  active: PropTypes.bool,
  children: PropTypes.node,
  open: PropTypes.bool, 
};

export default SidenavCollapse;
