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
	}

	return e
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
