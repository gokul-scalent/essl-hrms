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
	}
	return e
}

func EmailListModelToEmailListEntity(m model.EmailList) entity.EmailList {
	e := entity.EmailList{

		ID:               m.ID,
		User:             entity.User{ID: m.UserID},
		Name:             m.Name,
		TotalRecords:     m.TotalRecords,
		ProcessedRecords: m.ProcessedRecords,
		SafeCount:        m.SafeCount,
		RiskyCount:       m.RiskyCount,
		InvalIDCount:     m.InvalIDCount,
		UnknownCount:     m.UnknownCount,
		PendingCount:     m.PendingCount,
		Priority:         m.Priority,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt.Time,
		DeletedAt:        m.DeletedAt.Time,
	}
	return e
}

func LeadModelToLeadEntity(m model.Lead) entity.Lead {
	var verifiedOn *time.Time
	if m.VerifiedOn.Valid {
		verifiedOn = &m.VerifiedOn.Time
	}
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}
	var nextRetryAt *time.Time
	if m.NextRetryAt.Valid {
		nextRetryAt = &m.NextRetryAt.Time
	}

	e := entity.Lead{
		ID:                 m.ID,
		EmailList:          entity.EmailList{ID: m.EmailListID},
		Priority:           m.Priority,
		Email:              m.Email,
		EmailProvIDer:      m.EmailProvIDer.String,
		FirstName:          m.FirstName.String,
		LastName:           m.LastName.String,
		JobTitle:           m.JobTitle.String,
		Company:            m.Company.String,
		City:               m.City.String,
		Country:            m.Country.String,
		Industry:           m.Industry.String,
		IsSafe:             m.IsSafe,
		FinalStatus:        m.FinalStatus,
		IsReachable:        m.IsReachable.String,
		IsDisposable:       m.IsDisposable.Bool,
		IsRoleAccount:      m.IsRoleAccount.Bool,
		VerificationStatus: m.VerificationStatus,
		VerifiedOn:         verifiedOn,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
		DeletedAt:          deletedAt,
		RetryCount:         m.RetryCount,
		NextRetryAt:        nextRetryAt,
	}

	return e
}

func UserSettingModelToUserSettingEntity(m model.UserSetting) entity.UserSetting {
	e := entity.UserSetting{

		ID:                   m.ID,
		User:                 entity.User{ID: m.UserID},
		VerificationInterval: m.VerificationInterval,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt.Time,
		DeletedAt:            m.DeletedAt.Time,
	}
	return e
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
