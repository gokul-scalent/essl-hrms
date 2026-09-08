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
		IsPasswordSet: user.IsPasswordSet,
		EmpID:         user.EmpID,
		EmpName:       user.EmpName,
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
