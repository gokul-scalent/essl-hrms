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

	punches := []coreAPIModel.AttendancePunchResponse{}

	for _, punch := range e.Punches {

		checkIn := ""
		if punch.CheckIn != nil {
			checkIn = punch.CheckIn.Format("15:04")
		}
		checkOut := ""
		if punch.CheckOut != nil {
			checkOut = punch.CheckOut.Format("15:04")
		}
		punches = append(
			punches,
			coreAPIModel.AttendancePunchResponse{
				CheckIn:  checkIn,
				CheckOut: checkOut,
			},
		)
	}

	return coreAPIModel.DailyAttendanceLogResponse{
		EmpID:        e.EmpID,
		EmpName:      e.EmpName,
		Date:         e.Date.Format("2006-01-02"),
		Punches:      punches,
		WorkingHours: e.WorkingHours,
		Status:       e.Status,
	}
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
