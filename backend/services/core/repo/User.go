package repo

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/commonConstants"
	"github.com/scalent.io/scalent-hrms/entity/filters"

	filterPkg "github.com/scalent.io/scalent-hrms/pkg/filter"

	"github.com/scalent.io/scalent-hrms/internal/converter"
	"github.com/scalent.io/scalent-hrms/model"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type UserRepoImpl struct {
	db *sqlx.DB
}

func NewUserRepoImpl(db *sqlx.DB) (*UserRepoImpl, error) {
	return &UserRepoImpl{
		db: db,
	}, nil
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, user entity.User) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>user: CreateUser started", reqID)

	query := "INSERT INTO users (email, password, is_password_set, status, session_token) VALUES(?, ?, ?, ?, ? )"

	result, err := r.db.Exec(query, user.Email, user.Password, user.IsPasswordSet, user.Status, user.SessionToken)
	if err != nil {
		log.Error(err.Error(), reqID)
		// Active user with same email.then show email already exits
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return 0, errors.ResponseBadRequestError("User already exists")
		}
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>user: CreateUser completed & user id is "+strconv.Itoa(int(userID)), reqID)
	return int(userID), nil
}

func (r *UserRepoImpl) PartialUpdateUser(ctx context.Context, user entity.User) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>user:  PartialUpdateUser started for user id "+strconv.Itoa(user.ID), reqID)

	var columns []string
	var args []interface{}

	if user.Email != "" {
		columns = append(columns, "email=?")
		args = append(args, user.Email)
	}

	// Update employee/user name
	if user.EmpName != "" {
		columns = append(columns, "empname=?")
		args = append(args, user.EmpName)
	}

	// Update last login
	if !user.LastLoginAt.IsZero() {
		columns = append(columns, "last_login_at=?")
		args = append(args, user.LastLoginAt)
	}

	if user.Password != "" {
		columns = append(columns, "password=?")
		args = append(args, user.Password)
	}

	if user.SessionToken != "" {
		columns = append(columns, "session_token=?")
		args = append(args, user.SessionToken)
	}

	if user.Status != "" {
		columns = append(columns, "status=?")
		args = append(args, user.Status)
	}

	// Nothing to update
	if len(columns) == 0 {
		log.Info("core>repo>user: No fields to update for user id "+strconv.Itoa(user.ID), reqID)
		return nil
	}

	query := `
		UPDATE users
		SET ` + strings.Join(columns, ", ") + `
		WHERE id = ?
		AND deleted_at IS NULL
	`

	args = append(args, user.ID)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		log.Error(err.Error(), reqID) // Duplicate email
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {

			return errors.ResponseBadRequestError("Email already exists")
		}
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>user: PartialUpdateUser completed for user id "+strconv.Itoa(user.ID), reqID)
	return nil
}

func (r *UserRepoImpl) UpdateUser(ctx context.Context, user entity.User) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>user: UpdateUser started for user id "+strconv.Itoa(user.ID), reqID)

	query := "UPDATE users SET email=?, password=?, status=?, last_login_at=?, session_token=? WHERE id=?  	AND deleted_at IS NULL"

	_, err := r.db.Exec(query, user.Email, user.Password, user.Status, user.LastLoginAt, user.SessionToken, user.ID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>user: UpdateUser completed for user id "+strconv.Itoa(user.ID), reqID)
	return nil
}

func (r *UserRepoImpl) DeleteUser(ctx context.Context, userID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>user: DeleteUser started for user id "+strconv.Itoa(userID), reqID)

	query := "UPDATE users SET deleted_at = ? WHERE id = ?"

	_, err := r.db.Exec(query, time.Now(), userID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>user: DeleteUser completed for user id "+strconv.Itoa(userID), reqID)
	return nil
}

func (r *UserRepoImpl) GetUserbyID(ctx context.Context, userID int) (entity.User, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)

	log.Info(
		"core>repo>user: GetUserbyID started for user id "+
			strconv.Itoa(userID),
		reqID,
	)

	query := `
		SELECT
			u.id,
			u.email,
			u.password,
			u.is_password_set,
			u.status,

			GROUP_CONCAT(
				DISTINCT ur.role_id
				ORDER BY ur.role_id
			) AS role_ids,

			GROUP_CONCAT(
				DISTINCT r.code
				ORDER BY ur.role_id
			) AS role_codes,

			GROUP_CONCAT(
				DISTINCT r.name
				ORDER BY ur.role_id
			) AS role_names,

			GROUP_CONCAT(
				DISTINCT r.status
				ORDER BY ur.role_id
			) AS role_status,

			e.emp_id AS emp_id,
			COALESCE(e.emp_name, u.empname) AS emp_name,

			u.last_login_at,
			u.session_token,
			u.created_at,
			u.updated_at,
			u.deleted_at

		FROM users u

		LEFT JOIN user_roles ur
			ON ur.user_id = u.id
			AND ur.deleted_at IS NULL

		LEFT JOIN roles r
			ON r.id = ur.role_id
			AND r.deleted_at IS NULL
			AND r.status = 'ACTIVE'

		LEFT JOIN employees e
			ON e.uid = u.id
			AND e.deleted_at IS NULL

		WHERE u.id = ?
		AND u.deleted_at IS NULL

		GROUP BY
			u.id,
			u.email,
			u.password,
			u.is_password_set,
			u.status,
			e.emp_id,
			e.emp_name,
			u.last_login_at,
			u.session_token,
			u.created_at,
			u.updated_at,
			u.deleted_at
	`

	userModel := model.User{}
	userEntity := entity.User{}

	err := r.db.Get(&userModel, query, userID)
	if err != nil {
		log.Error(err.Error(), reqID)

		return userEntity,
			errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
	}

	userEntity = converter.UserModelToUserEntity(userModel)

	log.Info("core>repo>user: GetUserbyID completed for user id "+strconv.Itoa(userID), reqID)
	return userEntity, nil
}

func (r *UserRepoImpl) ListUser(ctx context.Context, filter *filters.ListFilter) (int, []entity.User, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>user: ListUser started", reqID)

	queryStatement := `
		SELECT
			u.id,
			u.email,
			u.password,
			u.is_password_set,
			u.status,

			GROUP_CONCAT(
				DISTINCT ur.role_id
				ORDER BY ur.role_id
			) AS role_ids,

			GROUP_CONCAT(
				DISTINCT r.code
				ORDER BY ur.role_id
			) AS role_codes,

			GROUP_CONCAT(
				DISTINCT r.name
				ORDER BY ur.role_id
			) AS role_names,

			GROUP_CONCAT(
				DISTINCT r.status
				ORDER BY ur.role_id
			) AS role_status,

			e.emp_id AS emp_id,
			COALESCE(e.emp_name, u.empname) AS emp_name,

			u.last_login_at,
			u.session_token,
			u.created_at,
			u.updated_at,
			u.deleted_at

		FROM users u

		LEFT JOIN user_roles ur
			ON ur.user_id = u.id
			AND ur.deleted_at IS NULL

		LEFT JOIN roles r
			ON r.id = ur.role_id
			AND r.deleted_at IS NULL
			AND r.status = 'ACTIVE'

		LEFT JOIN employees e
			ON e.uid = u.id
			AND e.deleted_at IS NULL
	`

	modelmap := model.UserModelMap

	whereStr, args := filterPkg.CreateFilterStr(
		filter.Filters,
		modelmap,
	)

	// Search
	if filter.SearchString != "" {
		search := "%" + strings.TrimSpace(filter.SearchString) + "%"

		whereStr = append(whereStr, "(u.email LIKE ? OR e.emp_name LIKE ? OR e.emp_id LIKE ?)")

		args = append(args, search, search, search)
	}

	// Soft delete
	whereStr = append(whereStr, "u.deleted_at IS NULL")
	whereString := strings.Join(whereStr, " AND ")

	if whereString != "" {
		queryStatement += " WHERE " + whereString
	}

	queryStatement += `
		GROUP BY
			u.id,
			u.email,
			u.password,
			u.is_password_set,
			u.status,
			e.emp_id,
			e.emp_name,
			u.last_login_at,
			u.session_token,
			u.created_at,
			u.updated_at,
			u.deleted_at
	`

	// Sort
	sortStr := filterPkg.CreateSortStr(filter.SortOption, modelmap)
	queryStatement += sortStr

	// Count users
	countQuery := `
		SELECT COUNT(*)
		FROM (
			` + queryStatement + `
		) AS result
	`

	var count int

	err := r.db.Get(&count, countQuery, args...)

	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil,
			errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	// Pagination
	if filter.Page == 0 {
		filter.Page = 1
	}

	offset := commonConstants.NO_OF_RECORDS_PER_PAGE *
		(filter.Page - 1)

	limitQuery := queryStatement +
		" LIMIT ?,?"

	argsWithPagination := append(
		append([]interface{}{}, args...),
		offset,
		commonConstants.NO_OF_RECORDS_PER_PAGE,
	)

	// Fetch users
	usersModel := []model.User{}

	err = r.db.Select(&usersModel, limitQuery, argsWithPagination...)

	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil,
			errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	// Convert to entity
	userEntities := []entity.User{}
	for _, userModel := range usersModel {
		userEntity := converter.UserModelToUserEntity(
			userModel,
		)
		userEntities = append(
			userEntities,
			userEntity,
		)
	}

	log.Info("core>repo>user: ListUser completed", reqID)
	return count, userEntities, nil
}
func (r *UserRepoImpl) GetUserDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.User, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>user: GetUserDetails started", reqID)

	selectStr := strings.Join(selectColumns, ", ")
	whereStr := strings.Join(whereColumn, " AND ")

	userModel := model.User{}

	query := "SELECT " + selectStr + " FROM " + table + " WHERE " + whereStr
	err := r.db.Get(&userModel, query, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	userEntity := converter.UserModelToUserEntity(userModel)

	log.Info("core>repo>user: GetUserDetails completed", reqID)
	return &userEntity, nil
}

func (r *UserRepoImpl) AssignUserRole(ctx context.Context, userID int, roleID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>user: AssignUserRole started for user id "+strconv.Itoa(userID)+" role id "+strconv.Itoa(roleID), reqID)
	query := `
		INSERT INTO user_roles (
			user_id,
			role_id
		)
		VALUES (?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		// Same user + same role already exists
		if mysqlErr, ok := err.(*mysql.MySQLError); ok &&
			mysqlErr.Number == 1062 {
			log.Error("user already has this role", reqID)
			return errors.ResponseBadRequestError("User already has this role")
		}

		log.Error("failed to assign user role: "+err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}
	log.Info("core>repo>user: AssignUserRole completed for user id "+strconv.Itoa(userID), reqID)
	return nil
}

func (r *UserRepoImpl) ChangePassword(ctx context.Context, password string, userID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>user:  ChangePassword started for user id "+strconv.Itoa(userID), reqID)

	// query := "UPDATE users SET password = ? WHERE id = ? AND users.deleted_at IS NULL"
	query := "UPDATE users SET password = ?, is_password_set = 'YES' WHERE id = ? AND users.deleted_at IS NULL"
	_, err := r.db.Exec(query, password, userID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>user: ChangePassword completed for user id "+strconv.Itoa(userID), reqID)
	return nil
}

func (r *UserRepoImpl) UpdateEmployeeName(ctx context.Context, userID int, empName string) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)

	query := `
		UPDATE employees
		SET emp_name = ?
		WHERE uid = ?
		AND deleted_at IS NULL
	`

	_, err := r.db.Exec(query, empName, userID)

	if err != nil {
		log.Error("failed to update employee name: "+err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}
	return nil
}

func (r *UserRepoImpl) UpdateUserPassword(ctx context.Context, userID int, hashedPassword string) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)

	query := `
		UPDATE users
		SET
			password = ?,
			is_password_set = 'NO',
			updated_at = NOW()
		WHERE id = ?
		AND deleted_at IS NULL
	`

	result, err := r.db.Exec(query, hashedPassword, userID)
	if err != nil {
		log.Error("failed to update user password: "+err.Error(), reqID)

		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error("failed to get affected rows: "+err.Error(), reqID)

		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	if rowsAffected == 0 {
		return errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
	}
	return nil
}

func (r *UserRepoImpl) UpdateUserRoles(ctx context.Context, userID int, roleIDs []int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>user: UpdateUserRoles started for user id "+strconv.Itoa(userID), reqID)

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction: "+err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	// Rollback automatically if any error occurs.
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Soft delete all currently active roles
	_, err = tx.ExecContext(
		ctx,
		`
		UPDATE user_roles
		SET deleted_at = NOW()
		WHERE user_id = ?
		AND deleted_at IS NULL
		`,
		userID,
	)

	if err != nil {
		log.Error("failed to remove existing user roles: "+err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	// 2. Add/restore selected roles
	for _, roleID := range roleIDs {
		if roleID <= 0 {
			continue
		}

		// First try to restore an existing soft-deleted role.
		result, err := tx.ExecContext(
			ctx,
			`
			UPDATE user_roles
			SET deleted_at = NULL
			WHERE user_id = ?
			AND role_id = ?
			AND deleted_at IS NOT NULL
			`,
			userID,
			roleID,
		)

		if err != nil {
			log.Error("failed to restore user role: "+err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Error("failed to get affected rows: "+err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}

		// If role did not previously exist, insert it.
		if rowsAffected == 0 {

			_, err = tx.ExecContext(
				ctx,
				`
				INSERT INTO user_roles (
					user_id,
					role_id
				)
				VALUES (?, ?)
				`,
				userID,
				roleID,
			)

			if err != nil {
				log.Error("failed to insert user role: "+err.Error(), reqID)
				return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
			}
		}
	}

	// 3. Commit transaction
	if err = tx.Commit(); err != nil {
		log.Error("failed to commit user roles: "+err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>user: UpdateUserRoles completed for user id "+strconv.Itoa(userID), reqID)
	return nil
}
