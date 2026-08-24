package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strconv"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/commonConstants"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	packageEmail "github.com/scalent.io/scalent-hrms/pkg/email"
	"github.com/scalent.io/scalent-hrms/pkg/email/smtp"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	"github.com/scalent.io/scalent-hrms/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepo UserRepo
	config   Config
}

func NewUserServiceImpl(userRepo UserRepo, config Config) (*UserServiceImpl, error) {
	return &UserServiceImpl{
		userRepo: userRepo,
		config:   config,
	}, nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user entity.User) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>user: create user started", reqID)

	// Generate temporary password
	password, err := utils.GenerateRandomPassword(12)
	if err != nil {
		log.Error("failed to generate password: "+err.Error(), reqID)

		return 0, errors.ResponseInternalServerError(
			errors.INTERNAL_SERVER_ERROR,
		)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Error("failed to hash password: "+err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(
			errors.INTERNAL_SERVER_ERROR,
		)
	}

	user.Password = string(hashedPassword)
	// User has NOT set their own password yet.
	user.IsPasswordSet = "NO"
	// New users are ACTIVE.
	if user.Status == "" {
		user.Status = "ACTIVE"
	}

	log.Info("creating user with email="+user.Email+" status="+user.Status, reqID)
	// Create user
	userID, errResp := s.userRepo.CreateUser(ctx, user)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}
	//for now we have role as admin so
	errResp = s.userRepo.AssignUserRole(ctx, userID, 1)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}
	// Send welcome email-
	if s.config.IsEmailSendingEnabled == "ACTIVE" {

		type WelcomeEmailData struct {
			UserName      string
			LoginEmail    string
			LoginPassword string
			LoginLink     string
		}

		emailData := WelcomeEmailData{
			UserName:      user.Email,
			LoginEmail:    user.Email,
			LoginPassword: password,
			LoginLink:     commonConstants.LOGIN_LINK,
		}

		templatePath := fmt.Sprintf(
			s.config.TemplatePath,
			s.config.WelcomeUserTemplateConstant,
		)

		templateConstant := s.config.WelcomeUserTemplateConstant
		emailSubject := s.config.WelcomeUserEmailSubject

		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			log.Error("Failed to parse welcome email template: "+err.Error(), reqID)
			return userID, nil
		}

		var emailHTML bytes.Buffer

		err = tmpl.ExecuteTemplate(&emailHTML, fmt.Sprintf("%s.html", templateConstant), emailData)
		if err != nil {
			log.Error("Failed to execute welcome email template: "+err.Error(), reqID)
			return userID, nil
		}

		emailPayload := packageEmail.EmailPayload{
			EmailSubject: emailSubject,
			ToEmail:      user.Email,
			HtmlContent:  emailHTML.String(),
			SenderEmail:  s.config.SenderEmail,
			SenderName:   s.config.SenderName,
		}

		emailSender := smtp.NewSmtpClient(
			&packageEmail.EmailConfig{
				ServerUrl:      s.config.SmtpHost,
				Port:           s.config.SmtpPort,
				AccessKey:      s.config.SmtpUsername,
				SecretKey:      s.config.SmtpPassword,
				EncryptionType: s.config.SmtpEncryptionType,
			},
		)
		_, _, err = emailSender.Send(ctx, emailPayload, "")
		if err != nil {
			log.Error("failed to send welcome email: "+err.Error(), reqID)
		}
	}

	log.Info("core>service>user: create user completed & user id is "+strconv.Itoa(userID), reqID)
	return userID, nil
}

func (s *UserServiceImpl) PartialUpdateUser(ctx context.Context, user entity.User) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>user: partila update user started for user id "+strconv.Itoa(user.ID), reqID)

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password: "+err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}
		user.Password = string(hashedPassword)
	}

	errResp := s.userRepo.PartialUpdateUser(ctx, user)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>user: update user completed for user id "+strconv.Itoa(user.ID), reqID)
	return nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, user entity.User) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>user: update user started for user id "+strconv.Itoa(user.ID), reqID)

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to hash password: "+err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}
		user.Password = string(hashedPassword)
	}

	errResp := s.userRepo.UpdateUser(ctx, user)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>user: update user completed for user id "+strconv.Itoa(user.ID), reqID)
	return nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>user: delete user started for user id "+strconv.Itoa(userID), reqID)

	errResp := s.userRepo.DeleteUser(ctx, userID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>user: delete user completed for user id "+strconv.Itoa(userID), reqID)
	return nil
}

func (s *UserServiceImpl) GetUserbyID(ctx context.Context, userID int) (entity.User, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>user: get user started for user id "+strconv.Itoa(userID), reqID)

	user, errResp := s.userRepo.GetUserbyID(ctx, userID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return entity.User{}, errResp
	}

	log.Info("core>service>user: user fetched successfully for user id "+strconv.Itoa(userID), reqID)
	return user, nil
}

func (s *UserServiceImpl) ListUser(ctx context.Context, filter *filters.ListFilter) (int, []entity.User, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>user: user list started", reqID)

	totalRecords, usersEntity, errResp := s.userRepo.ListUser(ctx, filter)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, nil, errResp
	}

	log.Info("core>service>user: user list completed", reqID)
	return totalRecords, usersEntity, nil
}

func (s *UserServiceImpl) ChangePassword(ctx context.Context, oldPassword string, newPassword string) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>user: ChangePassword started", reqID)

	sessionEntity, err := mailoraContext.GetSessionFromContext(ctx)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	userEntity, errResp := s.userRepo.GetUserDetails(ctx, []string{"id", "password", "is_password_set"}, "users", []string{"id=?"}, []interface{}{sessionEntity.UserID})
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	//  First-time password check set password or not
	if userEntity.Password == "" {
		return errors.ResponseBadRequestError("Please set your password first.")
	}

	//  Validate old password
	if userEntity.Password != "" {
		err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(oldPassword))
		if err != nil {
			if _, costErr := bcrypt.Cost([]byte(userEntity.Password)); costErr == nil {
				log.Error(err.Error(), reqID)
				return errors.ResponseBadRequestError("The old password is incorrect.")
			}

			if userEntity.Password != oldPassword {
				return errors.ResponseBadRequestError("The old password is incorrect.")
			}
		}
	}

	//  check same password reuse
	if oldPassword == newPassword {
		log.Error("New password same as old password", reqID)
		return errors.ResponseBadRequestError("New password cannot be same as old password.")
	}

	//  Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password: "+err.Error(), reqID)
		return errors.ResponseInternalServerError("Failed to change password. Please try again.")
	}

	//  Update password in DB
	errResp = s.userRepo.ChangePassword(ctx, string(hashedPassword), userEntity.ID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>user: ChangePassword completed for user id "+strconv.Itoa(userEntity.ID), reqID)
	return nil
}
