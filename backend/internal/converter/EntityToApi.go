package converter

import (
	"time"

	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
	"github.com/scalent.io/scalent-hrms/entity"
)

func UserEntityToUserAPIModelResponse(e entity.User) coreAPIModel.UserResponse {
	var lastLoginAt *time.Time

	if !e.LastLoginAt.IsZero() {
		lastLoginAt = &e.LastLoginAt
	}
	list := coreAPIModel.UserResponse{

		ID:            e.ID,
		Email:         e.Email,
		Status:        e.Status,
		IsPasswordSet: e.IsPasswordSet,
		LastLoginAt:   lastLoginAt,
	}
	return list
}

func EmailListEntityToEmailListAPIModelResponse(e entity.EmailList) coreAPIModel.EmailListResponse {
	list := coreAPIModel.EmailListResponse{

		ID:               e.ID,
		UserID:           e.User.ID,
		Name:             e.Name,
		TotalRecords:     e.TotalRecords,
		ProcessedRecords: e.ProcessedRecords,
		SafeCount:        e.SafeCount,
		RiskyCount:       e.RiskyCount,
		InvalIDCount:     e.InvalIDCount,
		UnknownCount:     e.UnknownCount,
		PendingCount:     e.PendingCount,
		Priority:         e.Priority,
		CreatedAt:        e.CreatedAt,
	}
	return list
}

func LeadEntityToLeadAPIModelResponse(e entity.Lead) coreAPIModel.LeadResponse {
	list := coreAPIModel.LeadResponse{

		ID:                  e.ID,
		EmailListID:         e.EmailList.ID,
		Priority:            e.Priority,
		Email:               e.Email,
		EmailProvIDer:       e.EmailProvIDer,
		FirstName:           e.FirstName,
		LastName:            e.LastName,
		JobTitle:            e.JobTitle,
		Company:             e.Company,
		City:                e.City,
		Country:             e.Country,
		Industry:            e.Industry,
		IsSafe:              e.IsSafe,
		FinalStatus:         e.FinalStatus,
		IsReachable:         e.IsReachable,
		IsDisposable:        e.IsDisposable,
		IsRoleAccount:       e.IsRoleAccount,
		VerificationStatus:  e.VerificationStatus,
		VerifiedOn:          e.VerifiedOn,
		VerificationAttempt: e.RetryCount,
	}
	return list
}

func UserSettingEntityToUserSettingAPIModelResponse(e entity.UserSetting) coreAPIModel.UserSettingResponse {
	list := coreAPIModel.UserSettingResponse{

		ID:                   e.ID,
		UserID:               e.User.ID,
		VerificationInterval: e.VerificationInterval,
	}
	return list
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
