package service

import (
	"context"
	"time"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
)

type LoginService interface {
	Login(ctx context.Context, identifier, password string) (*entity.User, string, errors.Response)
	LogOut(ctx context.Context) errors.Response
}

type HomeService interface {
	Home(ctx context.Context) errors.Response
}

type UserService interface {
	CreateUser(ctx context.Context, user entity.User) (int, errors.Response)
	PartialUpdateUser(ctx context.Context, user entity.User) errors.Response
	UpdateUser(ctx context.Context, user entity.User) errors.Response
	DeleteUser(ctx context.Context, userID int) errors.Response
	GetUserbyID(ctx context.Context, userID int) (entity.User, errors.Response)
	ListUser(ctx context.Context, filter *filters.ListFilter) (int, []entity.User, errors.Response)
	ChangePassword(ctx context.Context, oldPassword string, newPassword string) errors.Response
	SendUserMail(ctx context.Context, userID int) errors.Response
}

type EmployeeService interface {
	CreateEmployee(ctx context.Context, employee entity.Employee) (int, errors.Response)
	CreateEmployeeForUser(ctx context.Context, employee entity.Employee) (int, errors.Response)
	PartialUpdateEmployee(ctx context.Context, employee entity.Employee) errors.Response
	UpdateEmployee(ctx context.Context, employee entity.Employee) errors.Response
	DeleteEmployee(ctx context.Context, employeeID int) errors.Response
	GetEmployeebyID(ctx context.Context, employeeID int) (entity.Employee, errors.Response)
	ListEmployee(ctx context.Context, filter *filters.ListFilter) (int, []entity.Employee, errors.Response)
}

type AttendanceLogService interface {
	CreateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) (int, errors.Response)
	PartialUpdateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) errors.Response
	UpdateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) errors.Response
	GetAttendanceLogbyID(ctx context.Context, attendanceLogID int) (entity.AttendanceLog, errors.Response)
	ListAttendanceLog(ctx context.Context, filter *filters.ListFilter) (int, []entity.AttendanceLog, errors.Response)
	//	ListDailyAttendanceLog(ctx context.Context, filter *filters.ListFilter, empID, fromDate, toDate string) (int, []entity.DailyAttendanceLog, errors.Response)
	ListDailyAttendanceByHoursLog(ctx context.Context, filter *filters.ListFilter, empID, targetDate string) (int, []entity.DailyAttendanceLogHours, errors.Response)
}

type CronService interface {
	CalculateWorkingHours(ctx context.Context, targetDate time.Time) ([]entity.WorkingHoursSummary, errors.Response)
}

type RepoService interface {
	GetRoles(ctx context.Context) ([]entity.Role, errors.Response)
	GetRoleByID(ctx context.Context, id int) (*entity.Role, errors.Response)
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
