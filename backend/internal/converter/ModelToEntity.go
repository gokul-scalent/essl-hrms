package converter

import (
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
		LastLoginAt:   m.LastLoginAt.Time,
		SessionToken:  m.SessionToken.String,
		CreatedAt:     m.CreatedAt.Time,
		UpdatedAt:     m.UpdatedAt.Time,
		DeletedAt:     m.DeletedAt.Time,
		Role: entity.Role{
			ID:     int(m.RoleID.Int64),
			Name:   m.RoleName.String,
			Code:   m.RoleCode.String,
			Status: m.RoleStatus.String,
		},
		EmpID:   m.EmpID.String,
		EmpName: m.EmpName.String,
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

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
