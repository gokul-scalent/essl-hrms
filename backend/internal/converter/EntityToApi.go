package converter

import (
	apimodel "github.com/scalent.io/scalent-hrms/apimodel/core"
	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
	"github.com/scalent.io/scalent-hrms/entity"
)

func UserEntityToUserAPIModelResponse(user entity.User) apimodel.UserResponse {
	response := apimodel.UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Status:        user.Status,
		City:          user.City,
		BiometricSync: user.BiometricSync,
		IsPasswordSet: user.IsPasswordSet,
		EmpID:         user.EmpID,
		EmpName:       user.EmpName,
		Privilege:     user.Privilege,
	}

	if !user.LastLoginAt.IsZero() {
		response.LastLoginAt = &user.LastLoginAt
	}

	for i, roleID := range user.RoleIDs {
		role := apimodel.RoleResponse{
			ID: roleID,
		}

		if i < len(user.RoleNames) {
			role.Name = user.RoleNames[i]
		}

		if i < len(user.RoleCodes) {
			role.Code = user.RoleCodes[i]
		}

		if i < len(user.RoleStatus) {
			role.Status = user.RoleStatus[i]
		}

		response.Roles = append(response.Roles, role)
	}

	return response
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
		EmpName:         e.EmpName,
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
func RoleEntityToRoleAPIModelResponse(e entity.Role) coreAPIModel.RoleResponse {
	return coreAPIModel.RoleResponse{
		ID:     e.ID,
		Name:   e.Name,
		Code:   e.Code,
		Status: e.Status,
	}
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
