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
	"github.com/scalent.io/scalent-hrms/pkg/reacher"
	"github.com/scalent.io/scalent-hrms/pkg/validation"
)

type LeadRepoImpl struct {
	db *sqlx.DB
}

func NewLeadRepoImpl(db *sqlx.DB) (*LeadRepoImpl, error) {
	return &LeadRepoImpl{
		db: db,
	}, nil
}

func (r *LeadRepoImpl) CreateLead(ctx context.Context, lead entity.Lead) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: CreateLead started", reqID)

	query := "INSERT INTO leads (email_list_id, priority, email, email_provider, first_name, last_name, job_title, company, city, country, industry, is_safe, final_status, is_reachable, is_disposable, is_role_account, verification_status, verified_on) VALUES(?, ?, ?, ?, ?, ?, ?, ?,  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	var verifiedOn interface{}

	if lead.VerifiedOn == nil {
		verifiedOn = nil
	} else {
		verifiedOn = *lead.VerifiedOn
	}
	result, err := r.db.Exec(query, lead.EmailList.ID, lead.Priority, lead.Email, lead.EmailProvIDer, lead.FirstName, lead.LastName, lead.JobTitle, lead.Company, lead.City, lead.Country, lead.Industry, lead.IsSafe, lead.FinalStatus, lead.IsReachable, lead.IsDisposable, lead.IsRoleAccount, lead.VerificationStatus, verifiedOn)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	leadID, err := result.LastInsertId()
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>lead: CreateLead completed & lead id is "+strconv.Itoa(int(leadID)), reqID)
	return int(leadID), nil
}

func (r *LeadRepoImpl) PartialUpdateLead(ctx context.Context, lead entity.Lead, emailChanged bool) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info(
		"core>repo>lead: PartialUpdateLead started for lead id "+strconv.Itoa(lead.ID),
		reqID,
	)

	columns := []string{}
	args := []interface{}{}

	if lead.Email != "" {
		columns = append(columns, "email=?")
		args = append(args, lead.Email)
	}

	if lead.City != "" {
		columns = append(columns, "city=?")
		args = append(args, lead.City)
	}

	if lead.Company != "" {
		columns = append(columns, "company=?")
		args = append(args, lead.Company)
	}

	if lead.Country != "" {
		columns = append(columns, "country=?")
		args = append(args, lead.Country)
	}

	if lead.EmailList.ID != 0 {
		columns = append(columns, "email_list_id=?")
		args = append(args, lead.EmailList.ID)
	}

	if lead.FinalStatus != "" {
		columns = append(columns, "final_status=?")
		args = append(args, lead.FinalStatus)
	}

	if lead.FirstName != "" {
		columns = append(columns, "first_name=?")
		args = append(args, lead.FirstName)
	}

	if lead.Industry != "" {
		columns = append(columns, "industry=?")
		args = append(args, lead.Industry)
	}

	if lead.JobTitle != "" {
		columns = append(columns, "job_title=?")
		args = append(args, lead.JobTitle)
	}

	if lead.LastName != "" {
		columns = append(columns, "last_name=?")
		args = append(args, lead.LastName)
	}

	if lead.Priority != "" {
		columns = append(columns, "priority=?")
		args = append(args, lead.Priority)
	}
	// If email changed, reset verification state.
	if emailChanged {
		columns = append(
			columns,
			"verification_status=?",
			"is_safe=?",
			"final_status=?",
			"is_reachable=?",
			"is_disposable=?",
			"is_role_account=?",
			"email_provider=?",
			"verified_on=?",
		)

		args = append(
			args,
			commonConstants.VERIFICATION_STATUS_NOT_VERIFIED,
			commonConstants.LEAD_STATUS_UNKNOWN,
			commonConstants.FINAL_STATUS_UNKNOWN,
			commonConstants.LEAD_STATUS_UNKNOWN,
			false,
			false,
			"",
			nil,
		)
	} else {
		if lead.VerificationStatus != "" {
			columns = append(columns, "verification_status=?")
			args = append(args, lead.VerificationStatus)
		}

		if lead.VerifiedOn != nil {
			columns = append(columns, "verified_on=?")
			args = append(args, *lead.VerifiedOn)
		}
	}

	if len(columns) == 0 {
		return nil
	}
	args = append(args, lead.ID)

	columnStr := strings.Join(columns, ", ")

	if columnStr != "" {
		query := "UPDATE leads SET " + columnStr +
			" WHERE id=? AND deleted_at IS NULL"

		_, err := r.db.ExecContext(ctx, query, args...)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseInternalServerError(
				errors.INTERNAL_SERVER_ERROR,
			)
		}
	}

	log.Info("core>repo>lead: PartialUpdateLead completed for lead id "+strconv.Itoa(lead.ID), reqID)
	return nil
}

func (r *LeadRepoImpl) UpdateLead(ctx context.Context, lead entity.Lead) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: UpdateLead started for lead id "+strconv.Itoa(lead.ID), reqID)

	query := "UPDATE leads SET email_list_id=?, priority=?, email=?,  first_name=?, last_name=?, job_title=?, company=?, city=?, country=?, industry=?, is_safe=?, final_status=?, is_reachable=?, is_disposable=?, is_role_account=?, verification_status=?, verified_on=? WHERE id=?  	AND deleted_at IS NULL"

	_, err := r.db.Exec(query, lead.EmailList.ID, lead.Priority, lead.Email, lead.FirstName, lead.LastName, lead.JobTitle, lead.Company, lead.City, lead.Country, lead.Industry, lead.IsSafe, lead.FinalStatus, lead.IsReachable, lead.IsDisposable, lead.IsRoleAccount, lead.VerificationStatus, lead.VerifiedOn, lead.ID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>lead: UpdateLead completed for lead id "+strconv.Itoa(lead.ID), reqID)
	return nil
}

func (r *LeadRepoImpl) DeleteLead(ctx context.Context, leadID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: DeleteLead started for lead id "+strconv.Itoa(leadID), reqID)

	query := "UPDATE leads SET deleted_at = ? WHERE id = ?"

	_, err := r.db.Exec(query, time.Now(), leadID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>lead: DeleteLead completed for lead id "+strconv.Itoa(leadID), reqID)
	return nil
}

func (r *LeadRepoImpl) GetLeadbyID(ctx context.Context, leadID int) (entity.Lead, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: GetLeadbyID started for lead id "+strconv.Itoa(leadID), reqID)

	query := "SELECT * FROM leads WHERE leads.id=? AND leads.deleted_at IS NULL"

	leadModel := model.Lead{}
	leadEntity := entity.Lead{}

	err := r.db.Get(&leadModel, query, leadID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return leadEntity, errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
	}

	leadEntity = converter.LeadModelToLeadEntity(leadModel)

	log.Info("core>repo>lead: GetLeadbyID completed for lead id "+strconv.Itoa(leadID), reqID)
	return leadEntity, nil
}

func (r *LeadRepoImpl) ListLead(ctx context.Context, emailListID int, filter *filters.ListFilter) (int, []entity.Lead, entity.LeadStatusCount, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: ListLead started", reqID)

	queryStatement := "SELECT * FROM leads "

	modelmap := model.LeadModelMap

	whereStr, args := filterPkg.CreateFilterStr(filter.Filters, modelmap)

	log.Info(fmt.Sprintf("whereStr: %+v", whereStr), reqID)
	log.Info(fmt.Sprintf("args: %+v", args), reqID)
	// Search
	if filter.SearchString != "" {
		search := "%" + filter.SearchString + "%"

		whereStr = append(whereStr, `(email LIKE ? OR
		first_name LIKE ? OR
		last_name LIKE ? OR
		company LIKE ? )`)

		args = append(args,
			search,
			search,
			search,
			search,
		)
	}

	// Always filter by selected email list
	whereStr = append(whereStr, "email_list_id = ?")
	args = append(args, emailListID)

	// Exclude deleted records
	whereStr = append(whereStr, "leads.deleted_at IS NULL")

	queryStatement += " WHERE " + strings.Join(whereStr, " AND ")
	// Log the generated query and arguments
	log.Info("Final Query: "+queryStatement, reqID)
	log.Info(fmt.Sprintf("Final Args: %+v", args), reqID)

	sortStr := filterPkg.CreateSortStr(filter.SortOption, modelmap)
	queryStatement += sortStr

	var limitQueryStmt string

	emptySortOption := filters.SortOption{}

	totalRecordQueryStatement := "SELECT COUNT(id) as totalRecords FROM (" + queryStatement + ") as result  "
	var count int
	err := r.db.Get(&count, totalRecordQueryStatement, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, entity.LeadStatusCount{}, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	// Status counts for the email list
	var statusCount entity.LeadStatusCount

	countQuery := `
	SELECT
		COUNT(CASE WHEN is_safe = 'SAFE' THEN 1 END) AS safe_count,
		COUNT(CASE WHEN is_safe = 'RISKY' THEN 1 END) AS risky_count,
		COUNT(CASE WHEN is_safe = 'INVALID' THEN 1 END) AS invalid_count,
		COUNT(CASE WHEN is_safe = 'UNKNOWN' THEN 1 END) AS unknown_count,
		COUNT(CASE WHEN is_safe = 'PENDING' THEN 1 END) AS pending_count,
		COUNT(CASE WHEN is_safe = 'TIMEOUT' THEN 1 END) AS timeout_count
	FROM leads
	WHERE email_list_id = ?
	AND deleted_at IS NULL
`

	err = r.db.GetContext(ctx, &statusCount, countQuery, emailListID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, entity.LeadStatusCount{}, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
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

	leadsModel := []model.Lead{}

	err = r.db.Select(&leadsModel, limitQueryStmt, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, entity.LeadStatusCount{}, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	leadEntities := []entity.Lead{}

	for _, leadModel := range leadsModel {
		leadEntity := converter.LeadModelToLeadEntity(leadModel)
		leadEntities = append(leadEntities, leadEntity)
	}

	log.Info("core>repo>lead: ListLead completed", reqID)
	return count, leadEntities, statusCount, nil
}

func (r *LeadRepoImpl) GetLeadDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.Lead, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: GetLeadDetails started", reqID)

	selectStr := strings.Join(selectColumns, ", ")
	whereStr := strings.Join(whereColumn, " AND ")

	leadModel := model.Lead{}

	query := "SELECT " + selectStr + " FROM " + table + " WHERE " + whereStr
	err := r.db.Get(&leadModel, query, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	leadEntity := converter.LeadModelToLeadEntity(leadModel)

	log.Info("core>repo>lead: GetLeadDetails completed", reqID)
	return &leadEntity, nil
}

// get the pending leads by email list id
func (r *LeadRepoImpl) GetPendingLeadsByEmailListID(ctx context.Context, emailListID int, limit int) ([]entity.Lead, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: GetPendingLeadsByEmailListID started", reqID)

	// Fetch pending leads whose retry time has arrived.
	query := `
		SELECT *
		FROM leads
		WHERE id = ?
		  AND verification_status = 'PENDING'
		  AND deleted_at IS NULL
		  AND (next_retry_at IS NULL OR next_retry_at <= NOW()) 
		ORDER BY id ASC
		LIMIT ?
	`

	leadModels := []model.Lead{}
	log.Info(fmt.Sprintf("Executing query: %s with emailListID: %d and limit: %d", query, emailListID, limit), reqID)
	err := r.db.Select(&leadModels, query, emailListID, limit)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	leadEntities := make([]entity.Lead, 0, len(leadModels))
	for _, leadModel := range leadModels {
		leadEntities = append(leadEntities, converter.LeadModelToLeadEntity(leadModel))
	}

	log.Info("core>repo>lead: GetPendingLeadsByEmailListID completed", reqID)
	return leadEntities, nil
}

// check duplicate email
func (r *LeadRepoImpl) IsLeadExists(ctx context.Context, emailListID int, email string) (bool, errors.Response) {

	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: IsLeadExists started", reqID)

	query := `
		SELECT COUNT(1)
		FROM leads
		WHERE email_list_id = ?
		AND email = ?
		AND deleted_at IS NULL
	`

	var count int

	err := r.db.Get(&count, query, emailListID, email)
	if err != nil {
		log.Error(err.Error(), reqID)
		return false, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	return count > 0, nil
}

func (r *LeadRepoImpl) UpdateEmailListCounts(ctx context.Context, emailListID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: UpdateEmailListCounts started", reqID)

	query := `
		UPDATE email_lists
		SET
			total_records = (
				SELECT COUNT(*)
				FROM leads
				WHERE email_list_id = ? AND deleted_at IS NULL
			),
			processed_records = (
				SELECT COUNT(*)
				FROM leads
				WHERE email_list_id = ?
				AND verification_status <> ?
				AND deleted_at IS NULL
			),
			safe_count = (
				SELECT COUNT(*)
				FROM leads
				WHERE email_list_id = ?
				AND is_safe = ?
				AND deleted_at IS NULL
			),
			risky_count = (
				SELECT COUNT(*)
				FROM leads
				WHERE email_list_id = ?
				AND is_safe = ?
				AND deleted_at IS NULL
			),
			invalid_count = (
				SELECT COUNT(*)
				FROM leads
				WHERE email_list_id = ?
				AND is_safe = ?
				AND deleted_at IS NULL
			),
			unknown_count = (
				SELECT COUNT(*)
				FROM leads
				WHERE email_list_id = ?
				AND is_safe = ?
				AND deleted_at IS NULL
			),
			pending_count = (
				SELECT COUNT(*)
				FROM leads
				WHERE email_list_id = ?
				AND verification_status = ?
				AND deleted_at IS NULL
			)
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		emailListID,
		emailListID, commonConstants.VERIFICATION_STATUS_PENDING,
		emailListID, commonConstants.LEAD_STATUS_SAFE,
		emailListID, commonConstants.LEAD_STATUS_RISKY,
		emailListID, commonConstants.LEAD_STATUS_INVALID,
		emailListID, commonConstants.LEAD_STATUS_UNKNOWN,
		emailListID, commonConstants.VERIFICATION_STATUS_PENDING,
		emailListID,
	)

	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>lead: UpdateEmailListCounts completed", reqID)
	return nil
}

func (r *LeadRepoImpl) UpdateLeadVerification(ctx context.Context, leadID int, resp *reacher.VerifyEmailResponse) errors.Response {

	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: UpdateLeadVerification started", reqID)

	maxAttempts := 3

	//max 3 attempts and next verification at 30 sec for now
	if strings.EqualFold(resp.IsReachable, "timeout") {
		return r.HandleLeadTimeout(ctx, leadID, maxAttempts, 30*time.Second)
	}

	isSafe := commonConstants.LEAD_STATUS_UNKNOWN
	finalStatus := commonConstants.FINAL_STATUS_UNKNOWN
	verificationStatus := commonConstants.VERIFICATION_STATUS_COMPLETED

	switch strings.ToLower(resp.IsReachable) {
	case "safe":
		isSafe = commonConstants.LEAD_STATUS_SAFE
		finalStatus = commonConstants.FINAL_STATUS_GOOD

	case "risky":
		isSafe = commonConstants.LEAD_STATUS_RISKY
		finalStatus = commonConstants.FINAL_STATUS_RISKY

	case "invalid":
		isSafe = commonConstants.LEAD_STATUS_INVALID
		finalStatus = commonConstants.FINAL_STATUS_BAD

	default:
		isSafe = commonConstants.LEAD_STATUS_UNKNOWN
		finalStatus = commonConstants.FINAL_STATUS_UNKNOWN
	}

	emailProvider := "Unknown"

	if len(resp.MX.Records) > 0 {
		emailProvider = validation.GetEmailProviderFromMX(
			resp.MX.Records[0],
		)
	}

	query := `
		UPDATE leads
		SET
			is_reachable = ?,
			is_safe = ?,
			final_status = ?,
			is_disposable = ?,
			is_role_account = ?,
			email_provider = ?,
			verification_status = ?,
			verified_on = NOW(),
			next_retry_at = NULL
		WHERE id = ?
		AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		strings.ToUpper(resp.IsReachable),
		isSafe,
		finalStatus,
		resp.Misc.IsDisposable,
		resp.Misc.IsRoleAccount,
		emailProvider,
		verificationStatus,
		leadID,
	)

	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(
			errors.INTERNAL_SERVER_ERROR,
		)
	}

	log.Info("core>repo>lead: UpdateLeadVerification completed", reqID)
	return nil
}

func (r *LeadRepoImpl) GetSafeLeads(ctx context.Context, emailListID int, status string) ([]entity.Lead, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: get leads started", reqID)

	query := `
		SELECT
			email,
			first_name,
			last_name,
			job_title,
			company,
			city,
			country,
			industry
		FROM leads
		WHERE email_list_id = ?
	`

	args := []interface{}{emailListID}

	if status != "" {
		query += " AND is_safe = ?"
		args = append(args, status)
	}

	query += " ORDER BY id ASC"

	leadModels := []model.Lead{}

	err := r.db.SelectContext(ctx, &leadModels, query, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError("failed to fetch leads")
	}

	leadEntities := make([]entity.Lead, 0, len(leadModels))
	for _, leadModel := range leadModels {
		leadEntities = append(leadEntities, converter.LeadModelToLeadEntity(leadModel))
	}

	return leadEntities, nil
}

// re verify the lead
func (r *LeadRepoImpl) ReverifyLead(ctx context.Context, leadID int) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>lead: ReverifyLead started", reqID)

	query := `
	UPDATE leads
	SET
		verification_status = ?,
		is_safe = ?,
		final_status = ?,
		is_reachable = ?,
		is_disposable = ?,
		is_role_account = ?,
		email_provider = ?,
		retry_count = 0,
        next_retry_at = NULL,
		verified_on = NULL,
		updated_at = NOW()
	WHERE id = ?
	AND deleted_at IS NULL
	`

	_, err := r.db.Exec(
		query,
		"PENDING",
		"PENDING",
		"UNKNOWN",
		"UNKNOWN",
		false,
		false,
		"",
		leadID,
	)

	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>lead: ReverifyLead completed", reqID)

	return nil
}

func (r *LeadRepoImpl) GetLeadEmailListID(ctx context.Context, leadID int) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info(
		"core>repo>lead: GetLeadEmailListID started for lead id "+
			strconv.Itoa(leadID),
		reqID,
	)

	var emailListID int

	query := `SELECT email_list_id FROM leads WHERE id = ? AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &emailListID, query, leadID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
	}

	log.Info(
		"core>repo>lead: GetLeadEmailListID completed",
		reqID,
	)

	return emailListID, nil
}

func (r *LeadRepoImpl) HandleLeadTimeout(ctx context.Context, leadID int, maxRetries int, retryAfter time.Duration) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)

	log.Info("core>repo>lead: HandleLeadTimeout started", reqID)

	query := `
		UPDATE leads
		SET
			retry_count = retry_count + 1,
			next_retry_at = CASE
				WHEN retry_count + 1 < ? 
				THEN DATE_ADD(NOW(), INTERVAL ? SECOND)
				ELSE NULL
			END,
			verification_status = CASE
				WHEN retry_count + 1 >= ?
				THEN 'FAILED'
				ELSE 'PENDING'
			END,
			is_safe = CASE
				WHEN retry_count + 1 >= ?
				THEN 'TIMEOUT'
				ELSE 'PENDING'
			END,
			final_status = CASE
				WHEN retry_count + 1 >= ?
				THEN 'TIMEOUT'
				ELSE 'UNKNOWN'
			END,
			is_reachable = 'TIMEOUT',
			verified_on = CASE
				WHEN retry_count + 1 >= ?
				THEN NOW()
				ELSE NULL
			END,
			updated_at = NOW()
		WHERE id = ?
		AND deleted_at IS NULL
	`

	retrySeconds := int(retryAfter.Seconds())

	_, err := r.db.ExecContext(
		ctx,
		query,
		maxRetries,
		retrySeconds,
		maxRetries,
		maxRetries,
		maxRetries,
		maxRetries,
		leadID,
	)

	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(
			errors.INTERNAL_SERVER_ERROR,
		)
	}

	log.Info("core>repo>lead: HandleLeadTimeout completed", reqID)
	return nil
}
