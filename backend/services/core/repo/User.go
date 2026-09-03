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

	columns := []string{}
	args := []interface{}{}

	if user.Email != "" {
		columns = append(columns, "email=?")
		args = append(args, user.Email)
	}

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

	args = append(args, user.ID)

	columnStr := strings.Join(columns, ", ")

	if columnStr != "" {
		query := "UPDATE users SET " + columnStr + " WHERE  id=?  AND users.deleted_at IS NULL"

		_, err := r.db.Exec(query, args...)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}
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
	log.Info("core>repo>user: GetUserbyID started for user id "+strconv.Itoa(userID), reqID)

	query := `
		SELECT
			u.id, u.email, u.password, u.is_password_set, u.status,
			ur.role_id as role_id, r.name as role_name, r.code as role_code, r.status as role_status,
			e.emp_id as emp_id, e.emp_name as emp_name,
			u.last_login_at, u.session_token, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id AND ur.deleted_at IS NULL
		LEFT JOIN roles r ON r.id = ur.role_id
		LEFT JOIN employees e ON e.uid = u.id AND e.deleted_at IS NULL
		WHERE u.id = ?
		AND u.deleted_at IS NULL
	`

	userModel := model.User{}
	userEntity := entity.User{}

	err := r.db.Get(&userModel, query, userID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return userEntity, errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
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
			u.id, u.email, u.password, u.is_password_set, u.status,
			ur.role_id as role_id, r.name as role_name, r.code as role_code, r.status as role_status,
			e.emp_id as emp_id, e.emp_name as emp_name,
			u.last_login_at, u.session_token, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		LEFT JOIN user_roles ur ON ur.user_id = u.id AND ur.deleted_at IS NULL
		LEFT JOIN roles r ON r.id = ur.role_id
		LEFT JOIN employees e ON e.uid = u.id AND e.deleted_at IS NULL
	`

	modelmap := model.UserModelMap

	whereStr, args := filterPkg.CreateFilterStr(filter.Filters, modelmap)

	// Search string
	if filter.SearchString != "" {
		search := "%" + strings.TrimSpace(filter.SearchString) + "%"
		whereStr = append(whereStr, "u.email LIKE ?")
		args = append(args, search)
	}

	// Soft delete
	whereStr = append(whereStr, "u.deleted_at IS NULL")

	whereString := strings.Join(whereStr, " AND ")
	whereString = "WHERE " + whereString

	queryStatement = queryStatement + whereString

	sortStr := filterPkg.CreateSortStr(filter.SortOption, modelmap)
	queryStatement = queryStatement + sortStr

	var limitQueryStmt string

	emptySortOption := filters.SortOption{}

	totalRecordQueryStatement := "SELECT COUNT(id) as totalRecords FROM (" + queryStatement + ") as result"

	var count int
	err := r.db.Get(&count, totalRecordQueryStatement, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	if filter.Page == 0 && len(filter.Filters) == 0 && filter.SortOption == emptySortOption && filter.SearchString == "" {
		limitQueryStmt = queryStatement

	} else {
		if filter.Page == 0 {
			filter.Page = 1
		}

		offset := commonConstants.NO_OF_RECORDS_PER_PAGE * (filter.Page - 1)
		limitQueryStmt = queryStatement + " LIMIT ?,?"
		args = append(args, offset, commonConstants.NO_OF_RECORDS_PER_PAGE)
	}

	usersModel := []model.User{}

	err = r.db.Select(&usersModel, limitQueryStmt, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	userEntities := []entity.User{}

	for _, userModel := range usersModel {
		userEntity := converter.UserModelToUserEntity(userModel)
		userEntities = append(userEntities, userEntity)
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
	query := `
        INSERT INTO user_roles (user_id, role_id)
        VALUES (?, ?)
    `
	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return errors.ResponseInternalServerError(
			errors.INTERNAL_SERVER_ERROR,
		)
	}
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
