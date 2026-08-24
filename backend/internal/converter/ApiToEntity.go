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

func EmailListAPIToEmailListEntity(request *coreAPIModel.EmailList) entity.EmailList {
	e := entity.EmailList{

		User:             entity.User{ID: request.UserID},
		Name:             request.Name,
		TotalRecords:     request.TotalRecords,
		ProcessedRecords: request.ProcessedRecords,
		SafeCount:        request.SafeCount,
		RiskyCount:       request.RiskyCount,
		InvalIDCount:     request.InvalIDCount,
		UnknownCount:     request.UnknownCount,
		PendingCount:     request.PendingCount,
		Priority:         request.Priority,
	}

	return e
}

func CreateEmailListAPIToEmailListEntity(request *coreAPIModel.CreateEmailListRequest) entity.EmailList {
	e := entity.EmailList{
		Name:     request.Name,
		Priority: request.Priority,
	}

	return e
}

func UpdateEmailListAPIRequestToEmailListEntity(request *coreAPIModel.UpdateEmailListRequest) entity.EmailList {
	e := entity.EmailList{
		Name:     request.Name,
		Priority: request.Priority,
	}

	return e
}

func LeadAPIToLeadEntity(request *coreAPIModel.Lead) entity.Lead {
	e := entity.Lead{

		EmailList:          entity.EmailList{ID: request.EmailListID},
		Priority:           request.Priority,
		Email:              request.Email,
		EmailProvIDer:      request.EmailProvIDer,
		FirstName:          request.FirstName,
		LastName:           request.LastName,
		JobTitle:           request.JobTitle,
		Company:            request.Company,
		City:               request.City,
		Country:            request.Country,
		Industry:           request.Industry,
		IsSafe:             request.IsSafe,
		FinalStatus:        request.FinalStatus,
		IsReachable:        request.IsReachable,
		IsDisposable:       request.IsDisposable,
		IsRoleAccount:      request.IsRoleAccount,
		VerificationStatus: request.VerificationStatus,
		VerifiedOn:         request.VerifiedOn,
	}

	return e
}

func CreateLeadAPIToLeadEntity(request *coreAPIModel.CreateLeadRequest) entity.Lead {
	return entity.Lead{
		EmailList: entity.EmailList{
			ID: request.EmailListID,
		},
		Email: request.Email,
	}
}
func UpdateLeadAPIRequestToLeadEntity(request *coreAPIModel.UpdateLeadRequest) entity.Lead {
	e := entity.Lead{

		EmailList:          entity.EmailList{ID: request.EmailListID},
		Priority:           request.Priority,
		Email:              request.Email,
		EmailProvIDer:      request.EmailProvIDer,
		FirstName:          request.FirstName,
		LastName:           request.LastName,
		JobTitle:           request.JobTitle,
		Company:            request.Company,
		City:               request.City,
		Country:            request.Country,
		Industry:           request.Industry,
		IsSafe:             request.IsSafe,
		FinalStatus:        request.FinalStatus,
		IsReachable:        request.IsReachable,
		IsDisposable:       request.IsDisposable,
		IsRoleAccount:      request.IsRoleAccount,
		VerificationStatus: request.VerificationStatus,
		VerifiedOn:         request.VerifiedOn,
	}

	return e
}

func UserSettingAPIToUserSettingEntity(request *coreAPIModel.UserSetting) entity.UserSetting {
	e := entity.UserSetting{

		VerificationInterval: request.VerificationInterval,
	}

	return e
}
func UpdateUserSettingAPIRequestToUserSettingEntity(request *coreAPIModel.UpdateUserSettingRequest) entity.UserSetting {
	e := entity.UserSetting{

		VerificationInterval: request.VerificationInterval,
	}

	return e
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
