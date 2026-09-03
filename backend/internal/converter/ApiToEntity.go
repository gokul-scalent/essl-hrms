package converter

import (
	"strings"
	"time"

	"github.com/scalent.io/scalent-hrms/apimodel"
	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/commonConstants"
	"github.com/scalent.io/scalent-hrms/entity/filters"
)

func FilterAPIRequestToFilterEntity(request apimodel.ListFiltersRequest) filters.ListFilter {

	tempFilters := []filters.Filter{}
	for _, v := range request.Filters {
		tempFilter := filters.Filter{
			Field:     v.Field,
			Condition: v.Condition,
		}

		var tempFilterValues []string
		for _, filterValue := range v.FilterValues {
			tempFilterValue := strings.TrimSpace(filterValue)
			tempFilterValues = append(tempFilterValues, tempFilterValue)

		}

		tempFilter.FilterValues = tempFilterValues
		tempFilters = append(tempFilters, tempFilter)
	}

	tempSorter := filters.SortOption{
		SortBy:   request.SortOption.SortBy,
		SortType: request.SortOption.SortType,
	}

	filter := filters.ListFilter{
		Page:         request.Page,
		Filters:      tempFilters,
		SortOption:   tempSorter,
		SearchString: strings.TrimSpace(request.SearchString),
	}

	return filter
}

func UserAPIToUserEntity(request *coreAPIModel.User) entity.User {
	var lastLoginAt time.Time

	if request.LastLoginAt != nil {
		lastLoginAt = *request.LastLoginAt
	}
	e := entity.User{

		Email:         request.Email,
		Password:      request.Password,
		IsPasswordSet: commonConstants.NO,
		Status:        request.Status,
		LastLoginAt:   lastLoginAt,
		SessionToken:  request.SessionToken,
		Role: entity.Role{
			ID: request.RoleID,
		},
	}

	return e
}

func CreateUserAPIToUserEntity(request *coreAPIModel.CreateUser) entity.User {
	return entity.User{
		Email:  request.Email,
		Status: request.Status,
	}
}

func UpdateUserAPIRequestToUserEntity(request *coreAPIModel.UpdateUserRequest) entity.User {
	e := entity.User{

		Email:        request.Email,
		Password:     request.Password,
		Status:       request.Status,
		LastLoginAt:  request.LastLoginAt,
		SessionToken: request.SessionToken,
		Role: entity.Role{
			ID: request.RoleID,
		},
	}

	return e
}

func CreateEmployeeAPIRequestToEmployeeEntity(request *coreAPIModel.Employee) entity.Employee {
	e := entity.Employee{

		UID:       request.UID,
		EmpID:     request.EmpID,
		EmpName:   request.EmpName,
		Privilege: request.Privilege,
		Password:  request.Password,
		GroupID:   request.GroupID,
		Card:      request.Card,
	}

	return e
}

func UpdateEmployeeAPIRequestToEmployeeEntity(request *coreAPIModel.UpdateEmployeeRequest) entity.Employee {
	e := entity.Employee{

		UID:       request.UID,
		EmpID:     request.EmpID,
		EmpName:   request.EmpName,
		Privilege: request.Privilege,
		Password:  request.Password,
		GroupID:   request.GroupID,
		Card:      request.Card,
	}

	return e
}

func CreateAttendanceLogAPIToAttendanceLogEntity(request *coreAPIModel.AttendanceLog) entity.AttendanceLog {
	e := entity.AttendanceLog{

		UID:             request.UID,
		EmpID:           request.EmpID,
		Timestamp:       request.Timestamp,
		Status:          request.Status,
		Punch:           request.Punch,
		AttendanceState: request.AttendanceState,
		DeviceName:      request.DeviceName,
	}

	return e
}

func UpdateAttendanceLogAPIRequestToAttendanceLogEntity(request *coreAPIModel.UpdateAttendanceLogRequest) entity.AttendanceLog {
	e := entity.AttendanceLog{

		UID:             request.UID,
		EmpID:           request.EmpID,
		Timestamp:       request.Timestamp,
		Status:          request.Status,
		Punch:           request.Punch,
		AttendanceState: request.AttendanceState,
		DeviceName:      request.DeviceName,
	}

	return e
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
