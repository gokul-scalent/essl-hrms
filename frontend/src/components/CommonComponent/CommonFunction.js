import colors from "assets/theme/base/colors";

// src/components/common/commonStyles.js
const PASSWORD_SALT = "G8sK8xNm0qAr2zY6";

export const selectSx = (width = "12rem") => ({
  width: { xs: "100%", sm: width },
  minHeight: "2.73rem",
  borderRadius: "0.4285rem",
  "& .MuiOutlinedInput-root": {
    borderRadius: "0.4285rem",
    minHeight: "2.73rem",
    padding: 0,
    display: "flex",
    alignItems: "center",
    "& fieldset": { borderColor: "#d0d7de" },
    // "&:hover fieldset": { borderColor: "#0654a2" },
    "&.Mui-focused fieldset": {
      borderColor: "#0654a2",
      // boxShadow: "0 0 0 0 rgba(0,123,255,.25)",
    },
  },
  "& .MuiSelect-select": {
    padding: "0 0.75rem",
    display: "flex",
    alignItems: "center",
    fontSize: "0.75rem",
    justifyContent: "space-between", // arrow stays on right
    position: "relative",
    textTransform: "capitalize",
  },
  "& .MuiSelect-icon": {
    right: "0.5rem",
    color: "#6c757d",
    display: "block",
    transform: "scale(1.8)", // scale 50% bigger
  },
  "& .MuiMenuItem-root": {
    display: "flex",
    alignItems: "center",
    gap: "0.5rem", // badge + text spacing
  },
});

export const radioLabelSx = {
  "& .MuiFormControlLabel-label": {
    fontWeight: "normal",
    fontSize: "0.8rem",
  },
};

export const leadCsvStyles = {
  uploadBox: {
    border: `0.125rem dashed ${colors.grey[300]}`,
    borderRadius: "0.75rem",
    padding: "1.875rem",
    textAlign: "center",
    backgroundColor: colors.grey[100],
    position: "relative",
    minHeight: "8.75rem",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    flexDirection: "column",
  },
  uploadRemoveButton: {
    position: "absolute",
    top: 8,
    right: 8,
    width: 20,
    height: 20,
    borderRadius: "50%",
    backgroundColor: colors.deepBlue.main,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
  mappingContainer: {
    border: `0.063rem solid ${colors.grey[300]}`,
    borderRadius: "0.75rem",
    overflow: "hidden",
    backgroundColor: colors.white.main,
  },
  mappingHeader: {
    backgroundColor: colors.grey[100],
  },
  tableHeaderCell: {
    padding: "0.75rem 1rem",
    textAlign: "left",
    fontSize: "0.8rem",
    color: colors.grey[600],
  },
  tableCell: {
    padding: "0.85rem 1rem",
    verticalAlign: "top",
  },
  tableRow: {
    borderTop: `0.063rem solid ${colors.grey[300]}`,
  },
  selectInput: {
    width: "100%",
    padding: "0.55rem",
    borderRadius: "0.4rem",
    border: `0.063rem solid ${colors.grey[300]}`,
    backgroundColor: colors.white.main,
  },
  sampleText: {
    wordBreak: "break-word",
    color: colors.grey[800],
  },
};

export const leadListUiStyles = {
  tableState: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    minHeight: "73vh",
  },
  modalTitle: {
    mb: 3,
  },
  uploadTitle: {
    fontWeight: 600,
    fontSize: "0.95rem",
    color: colors.deepBlue.main,
    mb: 0.5,
  },
  uploadSubtitle: {
    fontSize: "0.95rem",
    color: colors.grey[600],
  },
  helperText: {
    fontWeight: 500,
    color: colors.grey[600],
  },
  primaryButton: {
    py: 1.2,
    px: 3,
    fontWeight: 700,
    fontSize: "0.9rem",
    textTransform: "none",
  },
  secondaryButton: {
    mr: 2,
    fontWeight: 700,
    fontSize: "0.9rem",
    textTransform: "none",
  },
};

export const authorizedUser = (res) => {
  if (res?.code === 440) {
    localStorage.clear();
    window.location.href = `/auth/login?error=Your session has expired`;
  } else {
    return;
  }
};

export const wrapCell = (value, maxWidth = "10rem") => (
  <div
    style={{
      whiteSpace: "normal",
      wordBreak: "break-word",
      overflowWrap: "break-word",
      lineHeight: "1.4",
      maxWidth,
    }}
  >
    {value ?? "-"}
  </div>
);

export const formatDateTime = (date) => {
  if (!date) return "-";

  return new Date(date).toLocaleString("en-IN", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
    timeZone: "UTC",
  });
};

export const formatDate = (date) => {
  if (!date) return "-";

  return new Date(date).toLocaleString("en-IN", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
};

export const encodePassword = (password) => {
  const salted = PASSWORD_SALT + password;
  return btoa(salted); // Convert the salted string into Base64 format
};
