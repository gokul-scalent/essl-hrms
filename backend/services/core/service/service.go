package service

import (
	"context"

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

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
