package converter

import (
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
		EmpName:         m.EmpName,
		Timestamp:       m.Timestamp,
		Status:          m.Status,
		Punch:           m.Punch,
		AttendanceState: m.AttendanceState,
		DeviceName:      m.DeviceName,
		CreatedAt:       m.CreatedAt.Time,
	}
	return e
}

func DailyAttendanceLogModelToEntity(m model.DailyAttendanceLog) entity.DailyAttendanceLog {
	return entity.DailyAttendanceLog{
		EmpID:   m.EmpID,
		EmpName: m.EmpName,
		Date:    m.LogDate,
		Punches: []entity.AttendancePunch{},
	}
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
