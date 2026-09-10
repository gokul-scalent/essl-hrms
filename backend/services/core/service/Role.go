package service

import (
	"context"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
)

type RoleService interface {
	GetRoles(ctx context.Context) ([]entity.Role, errors.Response)
	GetRoleByID(ctx context.Context, id int) (*entity.Role, errors.Response)
}

type RoleServiceImpl struct {
	roleRepo RoleRepo
}

func NewRoleServiceImpl(roleRepo RoleRepo) (*RoleServiceImpl, error) {
	return &RoleServiceImpl{
		roleRepo: roleRepo,
	}, nil
}

func (s *RoleServiceImpl) GetRoles(ctx context.Context) ([]entity.Role, errors.Response) {
	roles, errResp := s.roleRepo.GetRoles(ctx)
	if errResp != nil {
		return nil, errResp
	}
	return roles, nil
}

func (s *RoleServiceImpl) GetRoleByID(ctx context.Context, id int) (*entity.Role, errors.Response) {
	role, errResp := s.roleRepo.GetRoleByID(ctx, id)
	if errResp != nil {
		return nil, errResp
	}
	return role, nil
}
