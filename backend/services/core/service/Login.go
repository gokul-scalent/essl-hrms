package service

import (
	"context"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/commonConstants"
	"github.com/scalent.io/scalent-hrms/internal/auth"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type LoginServiceImpl struct {
	LoginRepo LoginRepo
	Auth      *auth.AuthImpl
	UserRepo  UserRepo
}

func NewLoginServiceImpl(LoginRepo LoginRepo, Auth *auth.AuthImpl, UserRepo UserRepo) (*LoginServiceImpl, error) {
	return &LoginServiceImpl{
		LoginRepo: LoginRepo,
		Auth:      Auth,
		UserRepo:  UserRepo,
	}, nil
}

func (s *LoginServiceImpl) Login(ctx context.Context, identifier, password string) (*entity.User, string, errors.Response) {

	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>login: Login started", reqID)

	var (
		userEntity   *entity.User
		sessionToken *string
		errResp      errors.Response
	)

	// 1. Fetch user

	userEntity, sessionToken, errResp =
		s.LoginRepo.GetUserDetailsForLogin(ctx, identifier)

	if errResp != nil {
		return nil, "", errResp
	}

	// 2. Status check
	if userEntity.Status != commonConstants.STATUS_ACTIVE {
		return nil, "", errors.ResponseUnauthorizedError("Account is inactive. Please contact the administrator.")
	}
	// 4. Password check
	if userEntity.Password == "" {
		return nil, "", errors.ResponseBadRequestError("Password not set")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(userEntity.Password),
		[]byte(password),
	); err != nil {
		if _, costErr := bcrypt.Cost([]byte(userEntity.Password)); costErr == nil {
			return nil, "", errors.ResponseUnauthorizedError("Invalid email or password")
		}

		if userEntity.Password != password {
			return nil, "", errors.ResponseUnauthorizedError("Invalid email or password")
		}

		// Legacy plain-text password detected. Migrate to bcrypt after successful match.
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			log.Error("failed to hash legacy password for user migration: "+hashErr.Error(), reqID)
		} else {
			errResp := s.UserRepo.ChangePassword(ctx, string(hashedPassword), userEntity.ID)
			if errResp != nil {
				log.Error("failed to migrate legacy password hash: "+errResp.Error(), reqID)
			}
		}
	}

	roleEntity, errResp := s.LoginRepo.GetUserRoleForLogin(ctx, userEntity.ID)
	if errResp != nil {
		return nil, "", errResp
	}

	if roleEntity == nil || len(roleEntity.RoleIDs) == 0 || len(roleEntity.Roles) == 0 {
		return nil, "", errors.ResponseUnauthorizedError(
			"User has no active role assigned",
		)
	}
	userEntity.RoleIDs = roleEntity.RoleIDs
	userEntity.Roles = roleEntity.Roles
	roleCode := userEntity.Roles[0].Code

	if roleCode == "" {
		return nil, "", errors.ResponseUnauthorizedError("User has no valid active role assigned")
	}

	log.Info("core>service>login: Session role:"+roleCode, reqID)

	if roleCode == "" {
		return nil, "", errors.ResponseUnauthorizedError("User has no valid active role assigned")
	}

	log.Info("core>service>login: Session role:"+roleCode, reqID)

	// 7. Remove previous session safely
	if sessionToken != nil && *sessionToken != "" {
		if err := s.Auth.RemovePreviousSessionToken(ctx, *sessionToken); err != nil {
			log.Error(err.Error(), reqID)
		}
	}

	// 8. Create new session
	sessionResp, errResp := s.Auth.CreateSession(
		ctx,
		auth.CreateSessionRequest{
			UserID: userEntity.ID,
			Role:   roleCode,
			TTL:    commonConstants.SESSION_TTL,
		},
	)

	if errResp != nil {
		return nil, "", errResp
	}

	// 8. Update DB
	if errResp := s.LoginRepo.UpdateLoginMeta(ctx, userEntity.ID, sessionResp.Token); errResp != nil {
		return nil, "", errResp
	}

	log.Info("core>service>login: Login completed", reqID)
	return userEntity, sessionResp.Token, nil
}

func (s *LoginServiceImpl) LogOut(ctx context.Context) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>login: logout started", reqID)

	session, err := mailoraContext.GetSessionFromContext(ctx)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseUnauthorizedError("Invalid session.")
	}

	errResp := s.LoginRepo.UpdateUserSessionToken(ctx, session.UserID, nil)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	errResp = s.Auth.DeleteSession(ctx)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>login: logout completed", reqID)
	return nil
}
