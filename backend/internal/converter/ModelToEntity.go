package converter

import (
	"strconv"
	"strings"
	"time"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/model"
)

func UserModelToUserEntity(m model.User) entity.User {
	e := entity.User{
		ID:            m.ID,
		Email:         m.Email.String,
		Password:      m.Password.String,
		IsPasswordSet: m.IsPasswordSet,
		Status:        m.Status.String,
		BiometricSync: m.BiometricSync,
		LastLoginAt:   m.LastLoginAt.Time,
		SessionToken:  m.SessionToken.String,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
		DeletedAt:     m.DeletedAt.Time,
		EmpID:         m.EmpID.String,
		EmpName:       m.EmpName.String,
		City:          m.City.String,
		Privilege:     m.Privilege,
	}

	// Role IDs
	if m.RoleIDs.Valid && m.RoleIDs.String != "" {
		roleIDs := strings.Split(m.RoleIDs.String, ",")

		for _, roleID := range roleIDs {
			id, err := strconv.Atoi(roleID)
			if err != nil {
				continue
			}

			e.RoleIDs = append(e.RoleIDs, id)
		}
	}

	// Role Codes
	if m.RoleCodes.Valid && m.RoleCodes.String != "" {
		e.RoleCodes = strings.Split(m.RoleCodes.String, ",")
	}
	// Role Names
	if m.RoleNames.Valid && m.RoleNames.String != "" {
		e.RoleNames = strings.Split(m.RoleNames.String, ",")
	}
	// Role Status
	if m.RoleStatus.Valid && m.RoleStatus.String != "" {
		e.RoleStatus = strings.Split(m.RoleStatus.String, ",")
	}

	return e
}

func EmployeeModelToEmployeeEntity(m model.Employee) entity.Employee {
	e := entity.Employee{

		ID:        m.ID,
		UID:       m.UID,
		EmpID:     m.EmpID,
		EmpName:   m.EmpName,
		Privilege: int(m.Privilege.Int64),
		Password:  m.Password.String,
		GroupID:   m.GroupID.String,
		Card:      m.Card.String,
		CreatedAt: m.CreatedAt.Time,
		DeletedAt: m.DeletedAt.Time,
	}
	return e
}

func AttendanceLogModelToAttendanceLogEntity(m model.AttendanceLog) entity.AttendanceLog {
	e := entity.AttendanceLog{

		ID:              m.ID,
		UID:             m.UID,
		EmpID:           m.EmpID,
		Timestamp:       m.Timestamp,
		Status:          m.Status,
		Punch:           m.Punch,
		AttendanceState: m.AttendanceState,
		DeviceName:      m.DeviceName,
		CreatedAt:       m.CreatedAt.Time,
	}
	return e
}

func RoleModelToRoleEntity(m model.Role) entity.Role {
	var updatedAt time.Time
	var deletedAt *time.Time

	if m.UpdatedAt.Valid {
		updatedAt = m.UpdatedAt.Time
	}
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return entity.Role{
		ID:        m.ID,
		Name:      m.Name,
		Code:      m.Code,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
}

func UserRoleModelsToUserEntity(userID int, roleModels []model.UserRoleModel) *entity.User {

	roleIDs := make([]int, 0, len(roleModels))
	roles := make([]entity.Role, 0, len(roleModels))

	for _, roleModel := range roleModels {

		roleIDs = append(roleIDs, roleModel.RoleID)

		roles = append(roles, entity.Role{
			ID:     roleModel.RoleID,
			Code:   roleModel.RoleCode,
			Name:   roleModel.RoleName,
			Status: roleModel.RoleStatus,
		})
	}

	return &entity.User{
		ID:      userID,
		RoleIDs: roleIDs,
		Roles:   roles,
	}
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
