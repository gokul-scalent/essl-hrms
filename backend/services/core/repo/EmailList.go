package repo

import (
	"context"
	"fmt"
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

type EmailListRepoImpl struct {
	db *sqlx.DB
}

func NewEmailListRepoImpl(db *sqlx.DB) (*EmailListRepoImpl, error) {
	return &EmailListRepoImpl{
		db: db,
	}, nil
}

func (r *EmailListRepoImpl) CreateEmailList(ctx context.Context, emailList entity.EmailList, userID int) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>emailList: CreateEmailList started", reqID)

	query := "INSERT INTO email_lists (user_id, name, total_records, processed_records, safe_count, risky_count, invalid_count, unknown_count, pending_count, priority) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	result, err := r.db.Exec(query, userID, emailList.Name, emailList.TotalRecords, emailList.ProcessedRecords, emailList.SafeCount, emailList.RiskyCount, emailList.InvalIDCount, emailList.UnknownCount, emailList.PendingCount, emailList.Priority)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	emailListID, err := result.LastInsertId()
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>emailList: CreateEmailList completed & email list id is "+strconv.Itoa(int(emailListID)), reqID)
	return int(emailListID), nil
}

func (r *EmailListRepoImpl) PartialUpdateEmailList(ctx context.Context, emailList entity.EmailList) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>emailList:  PartialUpdateEmailList started for email list id "+strconv.Itoa(emailList.ID), reqID)

	columns := []string{}
	args := []interface{}{}

	if emailList.InvalIDCount != 0 {
		columns = append(columns, "invalid_count=?")
		args = append(args, emailList.InvalIDCount)
	}

	if emailList.Name != "" {
		columns = append(columns, "name=?")
		args = append(args, emailList.Name)
	}

	if emailList.PendingCount != 0 {
		columns = append(columns, "pending_count=?")
		args = append(args, emailList.PendingCount)
	}

	if emailList.Priority != "" {
		columns = append(columns, "priority=?")
		args = append(args, emailList.Priority)
	}

	if emailList.ProcessedRecords != 0 {
		columns = append(columns, "processed_records=?")
		args = append(args, emailList.ProcessedRecords)
	}

	if emailList.RiskyCount != 0 {
		columns = append(columns, "risky_count=?")
		args = append(args, emailList.RiskyCount)
	}

	if emailList.SafeCount != 0 {
		columns = append(columns, "safe_count=?")
		args = append(args, emailList.SafeCount)
	}

	if emailList.TotalRecords != 0 {
		columns = append(columns, "total_records=?")
		args = append(args, emailList.TotalRecords)
	}

	if emailList.UnknownCount != 0 {
		columns = append(columns, "unknown_count=?")
		args = append(args, emailList.UnknownCount)
	}

	if emailList.User.ID != 0 {
		columns = append(columns, "user_id=?")
		args = append(args, emailList.User.ID)
	}

	args = append(args, emailList.ID)

	columnStr := strings.Join(columns, ", ")

	if columnStr != "" {
		query := "UPDATE email_lists SET " + columnStr + " WHERE  id=?  AND email_lists.deleted_at IS NULL"

		_, err := r.db.Exec(query, args...)
		if err != nil {
			log.Error(err.Error(), reqID)
			//check duplicate name
			if strings.Contains(err.Error(), "Error 1062") {
				return errors.ResponseBadRequestError("Email list name already exists")
			}

			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}
	}

	log.Info("core>repo>emailList: PartialUpdateEmailList completed for email list id "+strconv.Itoa(emailList.ID), reqID)
	return nil
}

func (r *EmailListRepoImpl) UpdateEmailList(ctx context.Context, emailList entity.EmailList) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>emailList: UpdateEmailList started for email list id "+strconv.Itoa(emailList.ID), reqID)

	query := "UPDATE email_lists SET name=?, total_records=?, processed_records=?, safe_count=?, risky_count=?, invalid_count=?, unknown_count=?, pending_count=?, priority=? WHERE id=?  	AND deleted_at IS NULL"

	_, err := r.db.Exec(query, emailList.Name, emailList.TotalRecords, emailList.ProcessedRecords, emailList.SafeCount, emailList.RiskyCount, emailList.InvalIDCount, emailList.UnknownCount, emailList.PendingCount, emailList.Priority, emailList.ID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>emailList: UpdateEmailList completed for email list id "+strconv.Itoa(emailList.ID), reqID)
	return nil
}

func (r *EmailListRepoImpl) DeleteEmailList(ctx context.Context, emailListID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>emailList: DeleteEmailList started for email list id "+strconv.Itoa(emailListID), reqID)

	query := "UPDATE email_lists SET deleted_at = ? WHERE id = ?"

	_, err := r.db.Exec(query, time.Now(), emailListID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>emailList: DeleteEmailList completed for email list id "+strconv.Itoa(emailListID), reqID)
	return nil
}

func (r *EmailListRepoImpl) GetEmailListbyID(ctx context.Context, emailListID int, userID int) (entity.EmailList, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>emailList: GetEmailListbyID started for email list id "+strconv.Itoa(emailListID), reqID)

	query := "SELECT * FROM email_lists WHERE email_lists.id=? AND email_lists.user_id = ? AND email_lists.deleted_at IS NULL"

	emailListModel := model.EmailList{}
	emailListEntity := entity.EmailList{}

	err := r.db.Get(&emailListModel, query, emailListID, userID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return emailListEntity, errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
	}

	emailListEntity = converter.EmailListModelToEmailListEntity(emailListModel)

	log.Info("core>repo>emailList: GetEmailListbyID completed for email list id "+strconv.Itoa(emailListID), reqID)
	return emailListEntity, nil
}

func (r *EmailListRepoImpl) ListEmailList(ctx context.Context, filter *filters.ListFilter, userID int) (int, []entity.EmailList, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>emailList: ListEmailList started", reqID)

	queryStatement := "SELECT * FROM email_lists "

	modelmap := model.EmailListModelMap

	whereStr, args := filterPkg.CreateFilterStr(filter.Filters, modelmap)
	// searchSTring
	if filter.SearchString != "" {
		search := "%" + strings.TrimSpace(filter.SearchString) + "%"
		whereStr = append(whereStr, "email_lists.name LIKE ?")
		args = append(args, search)
	}
	// Only return email lists belonging to logged-in user
	whereStr = append(whereStr, "email_lists.user_id = ?")
	args = append(args, userID)
	// Soft delete
	whereStr = append(whereStr, "email_lists.deleted_at IS NULL")

	whereString := strings.Join(whereStr, " AND ")
	whereString = "WHERE " + whereString

	queryStatement += whereString

	sortStr := filterPkg.CreateSortStr(filter.SortOption, modelmap)
	//get the newest added first
	if sortStr == "" {
		sortStr = " ORDER BY email_lists.created_at DESC"
	}

	queryStatement += sortStr

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

	} else {
		if filter.Page == 0 {
			filter.Page = 1
		}

		offset := commonConstants.NO_OF_RECORDS_PER_PAGE * (filter.Page - 1)
		limitQueryStmt = queryStatement + " LIMIT ?,?"
		args = append(args, offset, commonConstants.NO_OF_RECORDS_PER_PAGE)
	}

	emailListsModel := []model.EmailList{}

	err = r.db.Select(&emailListsModel, limitQueryStmt, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	emailListEntities := []entity.EmailList{}

	for _, emailListModel := range emailListsModel {
		emailListEntity := converter.EmailListModelToEmailListEntity(emailListModel)
		emailListEntities = append(emailListEntities, emailListEntity)
	}

	log.Info("core>repo>emailList: ListEmailList completed", reqID)
	return count, emailListEntities, nil
}

func (r *EmailListRepoImpl) GetEmailListDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.EmailList, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>emailList: GetEmailListDetails started", reqID)

	selectStr := strings.Join(selectColumns, ", ")
	whereStr := strings.Join(whereColumn, " AND ")

	emailListModel := model.EmailList{}

	query := "SELECT " + selectStr + " FROM " + table + " WHERE " + whereStr
	err := r.db.Get(&emailListModel, query, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	emailListEntity := converter.EmailListModelToEmailListEntity(emailListModel)

	log.Info("core>repo>emailList: GetEmailListDetails completed", reqID)
	return &emailListEntity, nil
}

// get the pending next email from list
func (r *EmailListRepoImpl) GetNextPendingEmailList(ctx context.Context) (*entity.EmailList, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>emailList: GetNextPendingEmailList started", reqID)

	query := `
		SELECT *
		FROM email_lists
		WHERE deleted_at IS NULL
		  AND pending_count > 0
		ORDER BY priority ASC, created_at ASC
		LIMIT 1
	`

	emailListModel := model.EmailList{}

	err := r.db.Get(&emailListModel, query)
	if err != nil {
		// No pending email list found
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}

		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	emailList := converter.EmailListModelToEmailListEntity(emailListModel)

	log.Info("core>repo>emailList: GetNextPendingEmailList completed", reqID)
	return &emailList, nil
}

func (r *EmailListRepoImpl) IsEmailListNameExists(ctx context.Context, userID int, name string) (bool, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>emailList: IsEmailListNameExists started", reqID)

	var count int

	query := `
		SELECT COUNT(*)
		FROM email_lists
		WHERE user_id = ?
		  AND LOWER(name) = LOWER(?)
		  AND deleted_at IS NULL
	`

	err := r.db.Get(&count, query, userID, name)
	if err != nil {
		log.Error(err.Error(), reqID)
		return false, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	return count > 0, nil
}

func (r *EmailListRepoImpl) GetNextEmailListForVerification(ctx context.Context) (*entity.EmailList, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>emailList: GetNextEmailListForVerification started", reqID)

	query := `
		SELECT DISTINCT  l.id as id
		FROM email_lists el
		INNER JOIN leads l ON l.email_list_id = el.id
		WHERE el.deleted_at IS NULL
		AND l.deleted_at IS NULL
		AND (
				l.verification_status IN ('PENDING', 'TIMEOUT')
				AND l.retry_count < 3
					AND (
						l.next_retry_at IS NULL
						OR l.next_retry_at <= NOW()
					)
				
			)
		ORDER BY el.priority ASC, el.created_at ASC
		LIMIT 1
	`

	emailListModel := model.EmailList{}
	fmt.Println(query)
	err := r.db.Get(&emailListModel, query)
	if err != nil {
		// No pending email list found
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}

		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	emailList := converter.EmailListModelToEmailListEntity(emailListModel)
	log.Info(fmt.Sprintf("**********************************************************%+v", emailList), reqID)
	log.Info("core>repo>emailList: GetNextEmailListForVerification completed", reqID)
	return &emailList, nil
}
