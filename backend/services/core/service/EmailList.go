package service

import (
	"context"
	"strconv"
	"strings"

	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type EmailListServiceImpl struct {
	emailListRepo EmailListRepo
}

func NewEmailListServiceImpl(emailListRepo EmailListRepo) (*EmailListServiceImpl, error) {
	return &EmailListServiceImpl{
		emailListRepo: emailListRepo,
	}, nil
}

func (s *EmailListServiceImpl) CreateEmailList(ctx context.Context, emailList entity.EmailList, req *coreAPIModel.CreateEmailListRequest) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>emailList: CreateEmailList started", reqID)

	session, _ := mailoraContext.GetSessionFromContext(ctx)

	// Check duplicate name
	existingEmailList, errResp := s.emailListRepo.GetEmailListDetails(
		ctx,
		[]string{"id"},
		"email_lists",
		[]string{
			"user_id = ?",
			"LOWER(name) = LOWER(?)",
			"deleted_at IS NULL",
		},
		[]interface{}{
			session.UserID,
			strings.TrimSpace(emailList.Name),
		},
	)

	if errResp == nil && existingEmailList != nil {
		return 0, errors.ResponseBadRequestError("Email list name already exists")
	}

	emailListID, errResp := s.emailListRepo.CreateEmailList(ctx, emailList, session.UserID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}

	log.Info("core>service>emailList: CreateEmailList completed & email list id is "+strconv.Itoa(emailListID), reqID)
	return emailListID, nil
}

func (s *EmailListServiceImpl) PartialUpdateEmailList(ctx context.Context, emailList entity.EmailList) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>emailList: partila update email list started for email list id "+strconv.Itoa(emailList.ID), reqID)

	// Check duplicate name if name is being updated
	if emailList.Name != "" {
		existingEmailList, errResp := s.emailListRepo.GetEmailListDetails(
			ctx,
			[]string{"id"},
			"email_lists",
			[]string{
				"id != ?",
				"LOWER(name) = LOWER(?)",
				"deleted_at IS NULL",
			},
			[]interface{}{
				emailList.ID,
				strings.TrimSpace(emailList.Name),
			},
		)

		if errResp == nil && existingEmailList != nil {
			return errors.ResponseBadRequestError("Email list name already exists")
		}
	}

	errResp := s.emailListRepo.PartialUpdateEmailList(ctx, emailList)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>emailList: update email list completed for email list id "+strconv.Itoa(emailList.ID), reqID)
	return nil
}

func (s *EmailListServiceImpl) UpdateEmailList(ctx context.Context, emailList entity.EmailList) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>emailList: update email list started for email list id "+strconv.Itoa(emailList.ID), reqID)

	// Check duplicate name if name is being updated
	if emailList.Name != "" {
		existingEmailList, errResp := s.emailListRepo.GetEmailListDetails(
			ctx,
			[]string{"id"},
			"email_lists",
			[]string{
				"id != ?",
				"LOWER(name) = LOWER(?)",
				"deleted_at IS NULL",
			},
			[]interface{}{
				emailList.ID,
				strings.TrimSpace(emailList.Name),
			},
		)

		if errResp == nil && existingEmailList != nil {
			return errors.ResponseBadRequestError("Email list name already exists")
		}
	}

	errResp := s.emailListRepo.UpdateEmailList(ctx, emailList)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>emailList: update email list completed for email list id "+strconv.Itoa(emailList.ID), reqID)
	return nil
}

func (s *EmailListServiceImpl) DeleteEmailList(ctx context.Context, emailListID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>emailList: delete email list started for email list id "+strconv.Itoa(emailListID), reqID)

	errResp := s.emailListRepo.DeleteEmailList(ctx, emailListID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>emailList: delete email list completed for email list id "+strconv.Itoa(emailListID), reqID)
	return nil
}

func (s *EmailListServiceImpl) GetEmailListbyID(ctx context.Context, emailListID int, userID int) (entity.EmailList, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>emailList: get email list started for email list id "+strconv.Itoa(emailListID), reqID)

	emailList, errResp := s.emailListRepo.GetEmailListbyID(ctx, emailListID, userID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return entity.EmailList{}, errResp
	}

	log.Info("core>service>emailList: email list fetched successfully for email list id "+strconv.Itoa(emailListID), reqID)
	return emailList, nil
}

func (s *EmailListServiceImpl) ListEmailList(ctx context.Context, filter *filters.ListFilter, userID int) (int, []entity.EmailList, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>emailList: email list list started", reqID)

	totalRecords, emailListsEntity, errResp := s.emailListRepo.ListEmailList(ctx, filter, userID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, nil, errResp
	}

	log.Info("core>service>emailList: email list list completed", reqID)
	return totalRecords, emailListsEntity, nil
}

func (s *EmailListServiceImpl) VerifyPendingEmails(ctx context.Context) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>emailList: VerifyPendingEmails started", reqID)

	log.Info("core>service>emailList: VerifyPendingEmails completed", reqID)
	return nil
}
