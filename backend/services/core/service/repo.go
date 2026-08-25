package service

import (
	"context"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/filters"
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

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
