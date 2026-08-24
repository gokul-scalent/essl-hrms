package service

import (
	"context"
	"time"

	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
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
}

type EmailListService interface {
	CreateEmailList(ctx context.Context, emailList entity.EmailList, req *coreAPIModel.CreateEmailListRequest) (int, errors.Response)
	PartialUpdateEmailList(ctx context.Context, emailList entity.EmailList) errors.Response
	UpdateEmailList(ctx context.Context, emailList entity.EmailList) errors.Response
	DeleteEmailList(ctx context.Context, emailListID int) errors.Response
	GetEmailListbyID(ctx context.Context, emailListID int, userID int) (entity.EmailList, errors.Response)
	ListEmailList(ctx context.Context, filter *filters.ListFilter, userID int) (int, []entity.EmailList, errors.Response)
	VerifyPendingEmails(ctx context.Context) errors.Response
}

type LeadService interface {
	CreateLead(ctx context.Context, lead entity.Lead, req *coreAPIModel.CreateLeadRequest) (int, errors.Response)
	PartialUpdateLead(ctx context.Context, lead entity.Lead) errors.Response
	UpdateLead(ctx context.Context, lead entity.Lead) errors.Response
	DeleteLead(ctx context.Context, leadID int) errors.Response
	GetLeadbyID(ctx context.Context, leadID int) (entity.Lead, errors.Response)
	ListLead(ctx context.Context, emailListID int, filter *filters.ListFilter) (int, []entity.Lead, entity.LeadStatusCount, errors.Response)
	GetSafeLeads(ctx context.Context, emailListID int, status string) ([]entity.Lead, errors.Response)
	VerifyPendingLeads(ctx context.Context) (time.Duration, errors.Response)
	HasPendingVerification(ctx context.Context) (bool, errors.Response)
	ReverifyLead(ctx context.Context, leadID int) errors.Response
}

type UserSettingService interface {
	CreateUserSetting(ctx context.Context, userSetting entity.UserSetting) (int, errors.Response)
	PartialUpdateUserSetting(ctx context.Context, userSetting entity.UserSetting) errors.Response
	UpdateUserSetting(ctx context.Context, userSetting entity.UserSetting) errors.Response
	DeleteUserSetting(ctx context.Context, userSettingID int) errors.Response
	GetUserSettingbyID(ctx context.Context, userSettingID int) (entity.UserSetting, errors.Response)
	ListUserSetting(ctx context.Context, filter *filters.ListFilter) (int, []entity.UserSetting, errors.Response)
	GetVerificationInterval(ctx context.Context, userID int) (int, errors.Response)
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
