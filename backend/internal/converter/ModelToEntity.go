package converter

import (
	"fmt"
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

func DailyAttendanceLogModelToEntity(m model.DailyAttendanceLog) entity.DailyAttendanceLog {
	e := entity.DailyAttendanceLog{

		EmpID:   m.EmpID,
		EmpName: m.EmpName,
		Date:    m.LogDate,
	}

	if m.CheckIn.Valid {
		checkInTime := m.CheckIn.Time
		e.CheckIn = &checkInTime
	}

	if m.CheckOut.Valid {
		checkOutTime := m.CheckOut.Time
		e.CheckOut = &checkOutTime
	}

	switch {
	case e.CheckIn != nil && e.CheckOut != nil:
		e.Status = "PRESENT"
		e.WorkingHours = formatWorkingHours(e.CheckOut.Sub(*e.CheckIn))
	case e.CheckIn != nil && e.CheckOut == nil:
		e.Status = "INCOMPLETE"
		e.WorkingHours = "00:00"
	default:
		e.Status = "ABSENT"
		e.WorkingHours = "00:00"
	}

	return e
}

func formatWorkingHours(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
