package repo

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scalent.io/scalent-hrms/entity"
	commonConstants "github.com/scalent.io/scalent-hrms/entity/commonConstants"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	"github.com/scalent.io/scalent-hrms/internal/converter"
	"github.com/scalent.io/scalent-hrms/model"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	filterPkg "github.com/scalent.io/scalent-hrms/pkg/filter"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type UserSettingRepoImpl struct {
	db *sqlx.DB
}

func NewUserSettingRepoImpl(db *sqlx.DB) (*UserSettingRepoImpl, error) {
	return &UserSettingRepoImpl{
		db: db,
	}, nil
}

func (r *UserSettingRepoImpl) CreateUserSetting(ctx context.Context, userSetting entity.UserSetting) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>userSetting: CreateUserSetting started", reqID)

	query := "INSERT INTO user_settings (user_id, verification_interval) VALUES(?, ?)"

	result, err := r.db.Exec(query, userSetting.User.ID, userSetting.VerificationInterval)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	userSettingID, err := result.LastInsertId()
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>userSetting: CreateUserSetting completed & user setting id is "+strconv.Itoa(int(userSettingID)), reqID)
	return int(userSettingID), nil
}

func (r *UserSettingRepoImpl) PartialUpdateUserSetting(ctx context.Context, userSetting entity.UserSetting) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>userSetting:  PartialUpdateUserSetting started for user setting id "+strconv.Itoa(userSetting.ID), reqID)

	columns := []string{}
	args := []interface{}{}

	if userSetting.User.ID != 0 {
		columns = append(columns, "user_id=?")
		args = append(args, userSetting.User.ID)
	}

	if userSetting.VerificationInterval != 0 {
		columns = append(columns, "verification_interval=?")
		args = append(args, userSetting.VerificationInterval)
	}

	args = append(args, userSetting.ID)
	// logged-in user ID
	args = append(args, userSetting.User.ID)
	columnStr := strings.Join(columns, ", ")

	if columnStr != "" {
		query := `
			UPDATE user_settings SET ` + columnStr + ` WHERE id=? AND user_id=? AND deleted_at IS NULL`

		_, err := r.db.Exec(query, args...)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}
	}

	log.Info("core>repo>userSetting: PartialUpdateUserSetting completed for user setting id "+strconv.Itoa(userSetting.ID), reqID)
	return nil
}

func (r *UserSettingRepoImpl) UpdateUserSetting(ctx context.Context, userSetting entity.UserSetting) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>userSetting: UpdateUserSetting started for user setting id "+strconv.Itoa(userSetting.ID), reqID)

	query := "UPDATE user_settings SET user_id=?, verification_interval=? WHERE id=?  	AND deleted_at IS NULL"

	_, err := r.db.Exec(query, userSetting.User.ID, userSetting.VerificationInterval, userSetting.ID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>userSetting: UpdateUserSetting completed for user setting id "+strconv.Itoa(userSetting.ID), reqID)
	return nil
}

func (r *UserSettingRepoImpl) DeleteUserSetting(ctx context.Context, userSettingID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>userSetting: DeleteUserSetting started for user setting id "+strconv.Itoa(userSettingID), reqID)

	query := "UPDATE user_settings SET deleted_at = ? WHERE id = ?"

	_, err := r.db.Exec(query, time.Now(), userSettingID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>userSetting: DeleteUserSetting completed for user setting id "+strconv.Itoa(userSettingID), reqID)
	return nil
}

func (r *UserSettingRepoImpl) GetUserSettingbyID(ctx context.Context, userSettingID int) (entity.UserSetting, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>userSetting: GetUserSettingbyID started for user setting id "+strconv.Itoa(userSettingID), reqID)

	query := "SELECT * FROM user_settings WHERE user_settings.id=? AND user_settings.deleted_at IS NULL"

	userSettingModel := model.UserSetting{}
	userSettingEntity := entity.UserSetting{}

	err := r.db.Get(&userSettingModel, query, userSettingID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return userSettingEntity, errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
	}

	userSettingEntity = converter.UserSettingModelToUserSettingEntity(userSettingModel)

	log.Info("core>repo>userSetting: GetUserSettingbyID completed for user setting id "+strconv.Itoa(userSettingID), reqID)
	return userSettingEntity, nil
}

func (r *UserSettingRepoImpl) ListUserSetting(ctx context.Context, filter *filters.ListFilter) (int, []entity.UserSetting, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>userSetting: ListUserSetting started", reqID)

	queryStatement := "SELECT * FROM user_settings "

	modelmap := model.UserSettingModelMap

	whereStr, args := filterPkg.CreateFilterStr(filter.Filters, modelmap)

	nullString := " AND user_settings.deleted_at IS NULL "

	if len(whereStr) == 0 {
		nullString = "  user_settings.deleted_at IS NULL "
	}

	whereString := strings.Join(whereStr, " AND ")
	whereString = "WHERE " + whereString

	queryStatement = queryStatement + whereString + nullString

	sortStr := filterPkg.CreateSortStr(filter.SortOption, modelmap)
	queryStatement = queryStatement + sortStr

	var limitQueryStmt string

	emptySortOption := filters.SortOption{}

	totalRecordQueryStatement := "SELECT COUNT(id) as totalRecords FROM (" + queryStatement + ") as result  "
	var count int
	err := r.db.Get(&count, totalRecordQueryStatement, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	if filter.Page == 0 && len(filter.Filters) == 0 && filter.SortOption == emptySortOption {
		limitQueryStmt = queryStatement

	} else if filter.Page == 0 || filter.Page != 0 || len(filter.Filters) > 0 || filter.SortOption != emptySortOption {
		if filter.Page == 0 {
			filter.Page = 1
		}

		offset := commonConstants.NO_OF_RECORDS_PER_PAGE * (filter.Page - 1)
		limitQueryStmt = queryStatement + " LIMIT ?,?"
		args = append(args, offset, commonConstants.NO_OF_RECORDS_PER_PAGE)
	}

	userSettingsModel := []model.UserSetting{}

	err = r.db.Select(&userSettingsModel, limitQueryStmt, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	userSettingEntities := []entity.UserSetting{}

	for _, userSettingModel := range userSettingsModel {
		userSettingEntity := converter.UserSettingModelToUserSettingEntity(userSettingModel)
		userSettingEntities = append(userSettingEntities, userSettingEntity)
	}

	log.Info("core>repo>userSetting: ListUserSetting completed", reqID)
	return count, userSettingEntities, nil
}

func (r *UserSettingRepoImpl) GetUserSettingDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.UserSetting, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>userSetting: GetUserSettingDetails started", reqID)

	selectStr := strings.Join(selectColumns, ", ")
	whereStr := strings.Join(whereColumn, " AND ")

	userSettingModel := model.UserSetting{}

	query := "SELECT " + selectStr + " FROM " + table + " WHERE " + whereStr
	err := r.db.Get(&userSettingModel, query, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	userSettingEntity := converter.UserSettingModelToUserSettingEntity(userSettingModel)

	log.Info("core>repo>userSetting: GetUserSettingDetails completed", reqID)
	return &userSettingEntity, nil
}

// get the verification interval time
func (r *UserSettingRepoImpl) GetVerificationInterval(ctx context.Context, userID int) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	var verificationInterval int
	query := `
		SELECT verification_interval
		FROM user_settings
		WHERE user_id = ?
		  AND deleted_at IS NULL
		LIMIT 1
	`

	err := r.db.Get(&verificationInterval, query, userID)
	if err != nil {
		// No setting found -> use default which we have set
		if strings.Contains(strings.ToLower(err.Error()), "no rows") {
			log.Info(
				"User setting not found, using default verification interval: 15 seconds",
				reqID,
			)
			return 15, nil
		}

		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(
			errors.INTERNAL_SERVER_ERROR,
		)
	}

	return verificationInterval, nil
}
