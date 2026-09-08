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

	roles := make([]coreAPIModel.RoleResponse, 0, len(e.Roles))

	for _, role := range e.Roles {
		roles = append(roles, coreAPIModel.RoleResponse{
			ID:     role.ID,
			Name:   role.Name,
			Code:   role.Code,
			Status: role.Status,
		})
	}

	return coreAPIModel.UserResponse{
		ID:            e.ID,
		Email:         e.Email,
		Status:        e.Status,
		IsPasswordSet: e.IsPasswordSet,
		LastLoginAt:   lastLoginAt,
		Roles:         roles,
		EmpID:         e.EmpID,
		EmpName:       e.EmpName,
	}
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

func RoleEntityToRoleAPIModelResponse(e entity.Role) coreAPIModel.RoleResponse {
	return coreAPIModel.RoleResponse{
		ID:     e.ID,
		Name:   e.Name,
		Code:   e.Code,
		Status: e.Status,
	}
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
