package service

import (
	"context"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	"github.com/scalent.io/scalent-hrms/model"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
)

type LoginRepo interface {
	GetUserDetailsForLogin(ctx context.Context, email string) (*entity.User, *string, errors.Response)
	UpdateUserSessionToken(ctx context.Context, userID int, token *string) errors.Response
	UpdateLoginMeta(ctx context.Context, userID int, token string) errors.Response
	GetUserRoleForLogin(ctx context.Context, userID int) (*entity.User, errors.Response)
}

type HomeRepo interface {
	Home(ctx context.Context) errors.Response
}

type UserRepo interface {
	CreateUser(ctx context.Context, user entity.User) (int, errors.Response)
	PartialUpdateUser(ctx context.Context, user entity.User) errors.Response
	UpdateUser(ctx context.Context, user entity.User) errors.Response
	DeleteUser(ctx context.Context, userID int) errors.Response
	GetUserbyID(ctx context.Context, userID int) (entity.User, errors.Response)
	ListUser(ctx context.Context, filter *filters.ListFilter) (int, []entity.User, errors.Response)
	GetUserDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.User, errors.Response)
	AssignUserRole(ctx context.Context, userID int, roleID int) errors.Response
	ChangePassword(ctx context.Context, password string, userID int) errors.Response
}

type EmployeeRepo interface {
	CreateEmployee(ctx context.Context, employee entity.Employee) (int, errors.Response)
	PartialUpdateEmployee(ctx context.Context, employee entity.Employee) errors.Response
	UpdateEmployee(ctx context.Context, employee entity.Employee) errors.Response
	DeleteEmployee(ctx context.Context, employeeID int) errors.Response
	GetEmployeebyID(ctx context.Context, employeeID int) (entity.Employee, errors.Response)
	ListEmployee(ctx context.Context, filter *filters.ListFilter) (int, []entity.Employee, errors.Response)
	GetEmployeeDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.Employee, errors.Response)
}

type AttendanceLogRepo interface {
	CreateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) (int, errors.Response)
	PartialUpdateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) errors.Response
	UpdateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) errors.Response
	GetAttendanceLogbyID(ctx context.Context, attendanceLogID int) (entity.AttendanceLog, errors.Response)
	ListAttendanceLog(ctx context.Context, filter *filters.ListFilter) (int, []entity.AttendanceLog, errors.Response)
	GetAttendanceLogDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.AttendanceLog, errors.Response)
	ListDailyAttendanceLog(ctx context.Context, filter *filters.ListFilter, empID, fromDate, toDate string) (int, []model.DailyAttendanceLog, errors.Response)
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
