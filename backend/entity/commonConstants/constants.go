package commonConstants

import "time"

const (
	NO_OF_RECORDS_PER_PAGE = 25
	SESSION_TTL            = 43200 // minutes (30 days)
)

// Status
const (
	STATUS_ACTIVE   = "ACTIVE"
	STATUS_INACTIVE = "INACTIVE"
)

// Lead Verification
const (
	VERIFICATION_BASE_DELAY = 15 * time.Second
	VERIFICATION_MAX_JITTER = 10 // seconds
)

// Lead Verification Status
const (
	VERIFICATION_STATUS_PENDING      = "PENDING"
	VERIFICATION_STATUS_COMPLETED    = "COMPLETED"
	VERIFICATION_STATUS_FAILED       = "FAILED"
	VERIFICATION_STATUS_NOT_VERIFIED = "NOT_VERIFIED"
)

// Lead Safety Status
const (
	LEAD_STATUS_PENDING = "PENDING"
	LEAD_STATUS_UNKNOWN = "UNKNOWN"
	LEAD_STATUS_TIMEOUT = "TIMEOUT"
	LEAD_STATUS_SAFE    = "SAFE"
	LEAD_STATUS_RISKY   = "RISKY"
	LEAD_STATUS_INVALID = "INVALID"
)

// Lead Final Status
const (
	FINAL_STATUS_GOOD      = "GOOD"
	FINAL_STATUS_BAD       = "BAD"
	FINAL_STATUS_RISKY     = "RISKY"
	FINAL_STATUS_UNKNOWN   = "UNKNOWN"
	FINAL_STATUS_CATCH_ALL = "CATCH_ALL"
)

// Upload Type
const (
	UPLOAD_TYPE_SINGLE = "single"
	UPLOAD_TYPE_CSV    = "csv"
)

// YES / NO
const (
	YES = "YES"
	NO  = "NO"
)

// login link
const (
	LOGIN_LINK = "http://mailora.scalent.io/auth/sign-in" //need to replace with live
)
