package converter

import (
	"time"

	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
	"github.com/scalent.io/scalent-hrms/entity"
)

func UserEntityToUserAPIModelResponse(e entity.User) coreAPIModel.UserResponse {
	var lastLoginAt *time.Time

	if !e.LastLoginAt.IsZero() {
		lastLoginAt = &e.LastLoginAt
	}
	list := coreAPIModel.UserResponse{

		ID:            e.ID,
		Email:         e.Email,
		Status:        e.Status,
		IsPasswordSet: e.IsPasswordSet,
		LastLoginAt:   lastLoginAt,
	}
	return list
}

func EmployeeEntityToUserAPIModelResponse(e entity.Employee) coreAPIModel.EmployeeResponse {
	list := coreAPIModel.EmployeeResponse{

		ID:        e.ID,
		UID:       e.UID,
		EmpID:     e.EmpID,
		EmpName:   e.EmpName,
		Privilege: e.Privilege,
		// Password:  e.Password,
		GroupID: e.GroupID,
		Card:    e.Card,
	}
	return list
}

func AttendanceLogEntityToAttendanceLogAPIModelResponse(e entity.AttendanceLog) coreAPIModel.AttendanceLogResponse {
	list := coreAPIModel.AttendanceLogResponse{

		ID:              e.ID,
		UID:             e.UID,
		EmpID:           e.EmpID,
		Timestamp:       e.Timestamp,
		Status:          e.Status,
		Punch:           e.Punch,
		AttendanceState: e.AttendanceState,
		DeviceName:      e.DeviceName,
	}
	return list
}

func DailyAttendanceLogEntityToAttendanceLogAPIModelResponse(e entity.DailyAttendanceLog) coreAPIModel.DailyAttendanceLogResponse {
	checkInStr := ""
	if e.CheckIn != nil {
		checkInStr = e.CheckIn.Format("15:04") //converting the utc date format into api response format eg: 2026-08-31 10:45:32 -> 10:45
	}

	checkOutStr := ""
	if e.CheckOut != nil {
		checkOutStr = e.CheckOut.Format("15:04")
	}

	r := coreAPIModel.DailyAttendanceLogResponse{

		EmpID:        e.EmpID,
		EmpName:      e.EmpName,
		Date:         e.Date.Format("2006-01-02"), //converting the utc date format into api response format eg: 2026-08-31 10:45:32 ->  2026-08-31
		CheckIn:      checkInStr,
		CheckOut:     checkOutStr,
		WorkingHours: e.WorkingHours,
		Status:       e.Status,
	}
	return r
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
