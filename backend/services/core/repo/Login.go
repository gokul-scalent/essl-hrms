package repo

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/model"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type LoginRepoImpl struct {
	db *sqlx.DB
}

func NewLoginRepoImpl(db *sqlx.DB) (*LoginRepoImpl, error) {
	return &LoginRepoImpl{
		db: db,
	}, nil
}

// Get User Details
func (r *LoginRepoImpl) GetUserDetailsForLogin(ctx context.Context, email string) (*entity.User, *string, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>login: GetUserDetailsForLogin started", reqID)

	userModel := model.User{}

	query := `
		SELECT id, email, password, status, is_password_set, session_token
		FROM users 
		WHERE email = ? AND deleted_at IS NULL
	`

	err := r.db.Get(&userModel, query, email)
	if err != nil {
		log.Error(err.Error(), reqID)
		if err == sql.ErrNoRows {
			return nil, nil, errors.ResponseNotFoundError("User not found.")
		}
		return nil, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	userEntity := entity.User{
		ID:            userModel.ID,
		Email:         userModel.Email.String,
		Password:      userModel.Password.String,
		Status:        userModel.Status.String,
		IsPasswordSet: userModel.IsPasswordSet,
	}

	var token *string
	if userModel.SessionToken.Valid {
		token = &userModel.SessionToken.String
	} else {
		token = nil
	}

	log.Info("core>repo>login: GetUserDetailsForLogin completed", reqID)
	return &userEntity, token, nil
}

func (r *LoginRepoImpl) UpdateLoginMeta(ctx context.Context, userID int, token string) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>login: UpdateLoginMeta started", reqID)

	query := `
		UPDATE users 
		SET last_login_at = NOW(),
		    session_token = ?
		WHERE id = ?
		AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, token, userID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>login: UpdateLoginMeta completed", reqID)
	return nil
}

func (r *LoginRepoImpl) UpdateUserSessionToken(ctx context.Context, userID int, token *string) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>login: UpdateUserSessionToken started", reqID)

	updateQuery := ` UPDATE users SET session_token = ?, last_login_at = ? WHERE id = ? AND deleted_at IS NULL `
	_, err := r.db.Exec(updateQuery, token, time.Now(), userID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>login: UpdateUserSessionToken completed", reqID)
	return nil
}

// Get User Role
func (r *LoginRepoImpl) GetUserRoleForLogin(
	ctx context.Context,
	userID int,
) (*entity.User, errors.Response) {

	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)

	log.Info("core>repo>login: GetUserRoleForLogin started", reqID)
	log.Info(
		"core>repo>login: GetUserRoleForLogin started for userID "+strconv.Itoa(userID),
		reqID,
	)

	type UserRoleModel struct {
		UserID     int    `db:"user_id"`
		RoleID     int    `db:"role_id"`
		RoleCode   string `db:"role_code"`
		RoleName   string `db:"role_name"`
		RoleStatus string `db:"role_status"`
	}

	var roleModels []UserRoleModel

	query := `
		SELECT
			ur.user_id,
			ur.role_id,
			r.code AS role_code,
			r.name AS role_name,
			r.status AS role_status
		FROM user_roles ur
		INNER JOIN roles r ON r.id = ur.role_id
		WHERE ur.user_id = ?
		  AND ur.deleted_at IS NULL
		  AND r.deleted_at IS NULL
		  AND r.status = 'ACTIVE'
		ORDER BY ur.role_id
	`

	err := r.db.SelectContext(ctx, &roleModels, query, userID)

	if err != nil {
		log.Error(
			"core>repo>login: failed to get user roles: "+err.Error(),
			reqID,
		)

		return nil, errors.ResponseInternalServerError(
			errors.INTERNAL_SERVER_ERROR,
		)
	}

	log.Info(
		"core>repo>login: roles found: "+strconv.Itoa(len(roleModels)),
		reqID,
	)

	if len(roleModels) == 0 {
		return nil, errors.ResponseUnauthorizedError(
			"User has no active role assigned",
		)
	}

	roleIDs := make([]int, 0, len(roleModels))
	roles := make([]entity.Role, 0, len(roleModels))

	for _, roleModel := range roleModels {

		log.Info(
			"core>repo>login: role found - ID: "+
				strconv.Itoa(roleModel.RoleID)+
				", Code: "+roleModel.RoleCode+
				", Name: "+roleModel.RoleName,
			reqID,
		)

		roleIDs = append(roleIDs, roleModel.RoleID)

		roles = append(roles, entity.Role{
			ID:     roleModel.RoleID,
			Code:   roleModel.RoleCode,
			Name:   roleModel.RoleName,
			Status: roleModel.RoleStatus,
		})
	}

	return &entity.User{
		ID:      userID,
		RoleIDs: roleIDs,
		Roles:   roles,
	}, nil
}
