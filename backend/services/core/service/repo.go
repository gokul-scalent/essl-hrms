package service

import (
	"context"
	"time"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/reacher"
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

type EmailListRepo interface {
	CreateEmailList(ctx context.Context, emailList entity.EmailList, userID int) (int, errors.Response)
	PartialUpdateEmailList(ctx context.Context, emailList entity.EmailList) errors.Response
	UpdateEmailList(ctx context.Context, emailList entity.EmailList) errors.Response
	DeleteEmailList(ctx context.Context, emailListID int) errors.Response
	GetEmailListbyID(ctx context.Context, emailListID int, userID int) (entity.EmailList, errors.Response)
	ListEmailList(ctx context.Context, filter *filters.ListFilter, userID int) (int, []entity.EmailList, errors.Response)
	GetEmailListDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.EmailList, errors.Response)
	GetNextPendingEmailList(ctx context.Context) (*entity.EmailList, errors.Response)
	GetNextEmailListForVerification(ctx context.Context) (*entity.EmailList, errors.Response)
}

type LeadRepo interface {
	CreateLead(ctx context.Context, lead entity.Lead) (int, errors.Response)
	PartialUpdateLead(ctx context.Context, lead entity.Lead, emailChanged bool) errors.Response
	UpdateLead(ctx context.Context, lead entity.Lead) errors.Response
	DeleteLead(ctx context.Context, leadID int) errors.Response
	GetLeadbyID(ctx context.Context, leadID int) (entity.Lead, errors.Response)
	ListLead(ctx context.Context, emailListID int, filter *filters.ListFilter) (int, []entity.Lead, entity.LeadStatusCount, errors.Response)
	GetLeadDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.Lead, errors.Response)
	GetPendingLeadsByEmailListID(ctx context.Context, emailListID int, limit int) ([]entity.Lead, errors.Response)
	IsLeadExists(ctx context.Context, emailListID int, email string) (bool, errors.Response)
	UpdateEmailListCounts(ctx context.Context, emailListID int) errors.Response
	UpdateLeadVerification(ctx context.Context, leadID int, resp *reacher.VerifyEmailResponse) errors.Response
	GetSafeLeads(ctx context.Context, emailListID int, status string) ([]entity.Lead, errors.Response)
	ReverifyLead(ctx context.Context, leadID int) errors.Response              //re verify the lead
	GetLeadEmailListID(ctx context.Context, leadID int) (int, errors.Response) //to update the parent after reverify or email change
	HandleLeadTimeout(ctx context.Context, leadID int, maxRetries int, retryAfter time.Duration) errors.Response
}

type UserSettingRepo interface {
	CreateUserSetting(ctx context.Context, userSetting entity.UserSetting) (int, errors.Response)
	PartialUpdateUserSetting(ctx context.Context, userSetting entity.UserSetting) errors.Response
	UpdateUserSetting(ctx context.Context, userSetting entity.UserSetting) errors.Response
	DeleteUserSetting(ctx context.Context, userSettingID int) errors.Response
	GetUserSettingbyID(ctx context.Context, userSettingID int) (entity.UserSetting, errors.Response)
	ListUserSetting(ctx context.Context, filter *filters.ListFilter) (int, []entity.UserSetting, errors.Response)
	GetUserSettingDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.UserSetting, errors.Response)
	GetVerificationInterval(ctx context.Context, userID int) (int, errors.Response)
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
