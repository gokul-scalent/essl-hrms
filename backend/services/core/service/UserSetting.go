package service

import (
	"context"
	"strconv"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type UserSettingServiceImpl struct {
	userSettingRepo UserSettingRepo
}

func NewUserSettingServiceImpl(userSettingRepo UserSettingRepo) (*UserSettingServiceImpl, error) {
	return &UserSettingServiceImpl{
		userSettingRepo: userSettingRepo,
	}, nil
}

func (s *UserSettingServiceImpl) CreateUserSetting(ctx context.Context, userSetting entity.UserSetting) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>userSetting: create user setting started", reqID)

	userSettingID, errResp := s.userSettingRepo.CreateUserSetting(ctx, userSetting)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}

	log.Info("core>service>userSetting: create user setting completed & user setting id is "+strconv.Itoa(userSettingID), reqID)
	return userSettingID, nil
}

func (s *UserSettingServiceImpl) PartialUpdateUserSetting(ctx context.Context, userSetting entity.UserSetting) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>userSetting: partila update user setting started for user setting id "+strconv.Itoa(userSetting.ID), reqID)

	errResp := s.userSettingRepo.PartialUpdateUserSetting(ctx, userSetting)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>userSetting: update user setting completed for user setting id "+strconv.Itoa(userSetting.ID), reqID)
	return nil
}

func (s *UserSettingServiceImpl) UpdateUserSetting(ctx context.Context, userSetting entity.UserSetting) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>userSetting: update user setting started for user setting id "+strconv.Itoa(userSetting.ID), reqID)

	errResp := s.userSettingRepo.UpdateUserSetting(ctx, userSetting)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>userSetting: update user setting completed for user setting id "+strconv.Itoa(userSetting.ID), reqID)
	return nil
}

func (s *UserSettingServiceImpl) DeleteUserSetting(ctx context.Context, userSettingID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>userSetting: delete user setting started for user setting id "+strconv.Itoa(userSettingID), reqID)

	errResp := s.userSettingRepo.DeleteUserSetting(ctx, userSettingID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>userSetting: delete user setting completed for user setting id "+strconv.Itoa(userSettingID), reqID)
	return nil
}

func (s *UserSettingServiceImpl) GetUserSettingbyID(ctx context.Context, userSettingID int) (entity.UserSetting, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>userSetting: get user setting started for user setting id "+strconv.Itoa(userSettingID), reqID)

	userSetting, errResp := s.userSettingRepo.GetUserSettingbyID(ctx, userSettingID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return entity.UserSetting{}, errResp
	}

	log.Info("core>service>userSetting: user setting fetched successfully for user setting id "+strconv.Itoa(userSettingID), reqID)
	return userSetting, nil
}

func (s *UserSettingServiceImpl) ListUserSetting(ctx context.Context, filter *filters.ListFilter) (int, []entity.UserSetting, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>userSetting: user setting list started", reqID)

	totalRecords, userSettingsEntity, errResp := s.userSettingRepo.ListUserSetting(ctx, filter)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, nil, errResp
	}

	log.Info("core>service>userSetting: user setting list completed", reqID)
	return totalRecords, userSettingsEntity, nil
}

// get the verification interval
func (s *UserSettingServiceImpl) GetVerificationInterval(ctx context.Context, userID int,
) (int, errors.Response) {

	return s.userSettingRepo.GetVerificationInterval(ctx, userID)
}
