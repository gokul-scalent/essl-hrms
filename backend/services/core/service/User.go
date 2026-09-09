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
	userRepo        UserRepo
	employeeService EmployeeService //added emp service to send insert records
	config          Config
}

func NewUserServiceImpl(userRepo UserRepo, employeeService EmployeeService, config Config) (*UserServiceImpl, error) {
	return &UserServiceImpl{
		userRepo:        userRepo,
		employeeService: employeeService,
		config:          config,
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
		user.Status = "INACTIVE"
	}

	log.Info("creating user with email="+user.Email+" status="+user.Status, reqID)
	// Create user
	userID, errResp := s.userRepo.CreateUser(ctx, user)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}
	//insert the new user to employee table
	employee := entity.Employee{
		UID:       userID,
		EmpID:     user.EmpID,
		EmpName:   user.EmpName,
		Privilege: user.Privilege,
	}

	_, errResp = s.employeeService.CreateEmployeeForUser(ctx, employee)
	if errResp != nil {
		log.Error("Failed to create employee: "+errResp.Error(), reqID)
		return 0, errResp
	}
	if user.RoleIDs != nil {
		for _, roleID := range user.RoleIDs {

			if roleID <= 0 {
				continue
			}

			errResp = s.userRepo.AssignUserRole(ctx, userID, roleID)
			if errResp != nil {
				log.Error(errResp.Error(), reqID)
				return 0, errResp
			}
		}
	}
	log.Info("CreateUser roleIDs="+fmt.Sprintf("%v", user.RoleIDs), reqID)
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

	// Update employee name
	if user.EmpName != "" {
		errResp = s.userRepo.UpdateEmployeeName(ctx, user.ID, user.EmpName)
		if errResp != nil {
			log.Error(errResp.Error(), reqID)
			return errResp
		}
	}

	if user.RoleIDs != nil {
		errResp = s.userRepo.UpdateUserRoles(ctx, user.ID, user.RoleIDs)

		if errResp != nil {
			log.Error(errResp.Error(), reqID)
			return errResp
		}
	}

	log.Info("core>service>user: partial update user completed for user id "+strconv.Itoa(user.ID), reqID)
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

	// Update employees.emp_name
	if user.EmpName != "" {
		errResp = s.userRepo.UpdateEmployeeName(ctx, user.ID, user.EmpName)

		if errResp != nil {
			log.Error(errResp.Error(), reqID)
			return errResp
		}
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

func (s *UserServiceImpl) SendUserMail(ctx context.Context, userID int) errors.Response {

	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>user: send user mail started for user id "+strconv.Itoa(userID), reqID)

	// 1. Get logged-in user's session
	sessionEntity, err := mailoraContext.GetSessionFromContext(ctx)
	if err != nil {
		log.Error("failed to get session: "+err.Error(), reqID)

		return errors.ResponseUnauthorizedError("access denied")
	}

	// 2. Only ADMIN and HR can send mail
	if sessionEntity.Role != "ADMIN" &&
		sessionEntity.Role != "HR" {

		log.Error("only ADMIN or HR can send user mail", reqID)

		return errors.ResponseUnauthorizedError("access denied")
	}

	// 3. Get target user
	user, errResp := s.userRepo.GetUserbyID(ctx, userID)
	if errResp != nil {
		log.Error("failed to get user: "+errResp.Error(), reqID)

		return errResp
	}

	if user.Email == "" {
		log.Error("user email is not available", reqID)

		return errors.ResponseBadRequestError("User email is not available")
	}

	// 4. Check email sending configuration
	if s.config.IsEmailSendingEnabled != "ACTIVE" {
		log.Error("email sending is disabled", reqID)

		return errors.ResponseBadRequestError("Email sending is disabled")
	}

	// 5. Generate new password
	password, err := utils.GenerateRandomPassword(12)
	if err != nil {
		log.Error("failed to generate password: "+err.Error(), reqID)

		return errors.ResponseInternalServerError(
			errors.INTERNAL_SERVER_ERROR,
		)
	}

	// 6. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Error("failed to hash generated password: "+err.Error(), reqID)

		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	// 7. Update password in database
	errResp = s.userRepo.UpdateUserPassword(ctx, userID, string(hashedPassword))

	if errResp != nil {
		log.Error("failed to update user password: "+errResp.Error(), reqID)
		return errResp
	}

	// 8. Prepare email template data
	type UserEmailData struct {
		UserName      string
		LoginEmail    string
		LoginPassword string
		LoginLink     string
	}

	emailData := UserEmailData{
		UserName:      user.Email,
		LoginEmail:    user.Email,
		LoginPassword: password,
		LoginLink:     commonConstants.LOGIN_LINK,
	}
	// 9. Load email template

	templatePath := fmt.Sprintf(
		s.config.TemplatePath,
		s.config.WelcomeUserTemplateConstant,
	)

	templateConstant := s.config.WelcomeUserTemplateConstant

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Error("failed to parse welcome email template: "+err.Error(), reqID)

		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	// 10. Execute template
	var emailHTML bytes.Buffer

	err = tmpl.ExecuteTemplate(
		&emailHTML,
		fmt.Sprintf("%s.html", templateConstant),
		emailData,
	)

	if err != nil {
		log.Error("failed to execute welcome email template: "+err.Error(), reqID)

		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	// 11. Build email payload
	emailPayload := packageEmail.EmailPayload{
		EmailSubject: s.config.WelcomeUserEmailSubject,
		ToEmail:      user.Email,
		HtmlContent:  emailHTML.String(),
		SenderEmail:  s.config.SenderEmail,
		SenderName:   s.config.SenderName,
	}

	// 12. Create SMTP client
	emailSender := smtp.NewSmtpClient(
		&packageEmail.EmailConfig{
			ServerUrl:      s.config.SmtpHost,
			Port:           s.config.SmtpPort,
			AccessKey:      s.config.SmtpUsername,
			SecretKey:      s.config.SmtpPassword,
			EncryptionType: s.config.SmtpEncryptionType,
		},
	)

	// 13. Send email
	_, _, err = emailSender.Send(ctx, emailPayload, "")

	if err != nil {
		log.Error("failed to send user mail: "+err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}
	log.Info("core>service>user: send user mail completed for user id "+strconv.Itoa(userID), reqID)

	return nil
}
