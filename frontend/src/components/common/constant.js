import colors from "assets/theme/base/colors";

// constants.js
export const ROLES = {
  ADMIN: "ADMIN",
};

// Role-based default route for now
export const DEFAULT_ROUTES = {
  [ROLES.ADMIN]: "/admin/dashboard",
};

export const LEADS_STATUS = [
  { label: "RISKY", value: "Risky" },
  { label: "INVALID", value: "Invalid" },
  { label: "SAFE", value: "Safe" },
  { label: "UNKNOWN", value: "Unknown" },
  { label: "PENDING", value: "Pending" },
  { label: "TIMEOUT", value: "Timeout" },
];

//Regex
export const REGEX = {
  //email regex
  EMAIL_REGEX:
    /^[a-zA-Z0-9]+([._%+-]?[a-zA-Z0-9]+)*@([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]{2,}$/,
};

export const DEFAULT_RECORDS_PER_PAGE = 25;

//priority
export const PRIORITY = [
  { label: "P0", value: "p0", color: "#fa8181" },
  { label: "P1", value: "p1", color: "#fdbe71" },
];

//color to badge of lead status
export const STATUS_BADGE_COLORS = {
  SAFE: colors.badgeColors.success,
  RISKY: colors.badgeColors.warning, 
  INVALID: colors.badgeColors.error, 
  UNKNOWN: colors.badgeColors.secondary, 
  PENDING: colors.badgeColors.info, 
  TIMEOUT: colors.badgeColors.dark, 
};

//verification status for lead
export const VERIFICATION_STATUS = [
  { label: "PENDING", value: "Pending" },
  { label: "COMPLETED", value: "Completed" },
  { label: "FAILED", value: "Failed" },
  { label: "NOT_VERIFIED", value: "Not-Verified" },
];


export const DOWNLOAD_OPTIONS = [
  { label: "All", value: "" },
  { label: "SAFE", value: "safe" },
  { label: "RISKY", value: "risky" },
  { label: "INVALID", value: "invalid" },
  { label: "UNKNOWN", value: "unknown" },
  { label: "PENDING", value: "pending" },
];

export const STATUS = [
  { label: "ACTIVE", value: "Active" },
  { label: "INACTIVE", value: "Inactive" },
];

export const ATTENDANCE_LOG_STATE = [
  { label: "CHECK_OUT", value: "Check Out" },
  { label: "CHECK_IN", value: "Check In" },
];