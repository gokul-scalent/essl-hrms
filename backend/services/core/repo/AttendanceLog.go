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

type AttendanceLogRepoImpl struct {
	db *sqlx.DB
}

func NewAttendanceLogRepoImpl(db *sqlx.DB) (*AttendanceLogRepoImpl, error) {
	return &AttendanceLogRepoImpl{
		db: db,
	}, nil
}

func (r *AttendanceLogRepoImpl) CreateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>attendanceLog: CreateAttendanceLog started", reqID)

	query := "INSERT INTO attendance_logs (uid, emp_id, timestamp, status, punch, attendance_state, device_name) VALUES(?, ?, ?, ?, ?, ?, ?)"

	result, err := r.db.Exec(query, attendanceLog.UID, attendanceLog.EmpID, attendanceLog.Timestamp, attendanceLog.Status, attendanceLog.Punch, attendanceLog.AttendanceState, attendanceLog.DeviceName)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	attendanceLogID, err := result.LastInsertId()
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>attendanceLog: CreateAttendanceLog completed & attendance log id is "+strconv.Itoa(int(attendanceLogID)), reqID)
	return int(attendanceLogID), nil
}

func (r *AttendanceLogRepoImpl) PartialUpdateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>attendanceLog:  PartialUpdateAttendanceLog started for attendance log id "+strconv.Itoa(attendanceLog.ID), reqID)

	columns := []string{}
	args := []interface{}{}

	if attendanceLog.AttendanceState != "" {
		columns = append(columns, "attendance_state=?")
		args = append(args, attendanceLog.AttendanceState)
	}

	if attendanceLog.DeviceName != "" {
		columns = append(columns, "device_name=?")
		args = append(args, attendanceLog.DeviceName)
	}

	if attendanceLog.EmpID != "" {
		columns = append(columns, "emp_id=?")
		args = append(args, attendanceLog.EmpID)
	}

	if attendanceLog.Punch != 0 {
		columns = append(columns, "punch=?")
		args = append(args, attendanceLog.Punch)
	}

	if attendanceLog.Status != 0 {
		columns = append(columns, "status=?")
		args = append(args, attendanceLog.Status)
	}

	if attendanceLog.Timestamp.IsZero() {
		columns = append(columns, "timestamp=?")
		args = append(args, attendanceLog.Timestamp)
	}

	if attendanceLog.UID != 0 {
		columns = append(columns, "uid=?")
		args = append(args, attendanceLog.UID)
	}

	args = append(args, attendanceLog.ID)

	columnStr := strings.Join(columns, ", ")

	if columnStr != "" {
		query := "UPDATE attendance_logs SET " + columnStr + " WHERE  id=?  "

		_, err := r.db.Exec(query, args...)
		if err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}
	}

	log.Info("core>repo>attendanceLog: PartialUpdateAttendanceLog completed for attendance log id "+strconv.Itoa(attendanceLog.ID), reqID)
	return nil
}

func (r *AttendanceLogRepoImpl) UpdateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>attendanceLog: UpdateAttendanceLog started for attendance log id "+strconv.Itoa(attendanceLog.ID), reqID)

	query := "UPDATE attendance_logs SET uid=?, emp_id=?, timestamp=?, status=?, punch=?, attendance_state=?, device_name=? WHERE id=?  "

	_, err := r.db.Exec(query, attendanceLog.UID, attendanceLog.EmpID, attendanceLog.Timestamp, attendanceLog.Status, attendanceLog.Punch, attendanceLog.AttendanceState, attendanceLog.DeviceName, attendanceLog.ID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>attendanceLog: UpdateAttendanceLog completed for attendance log id "+strconv.Itoa(attendanceLog.ID), reqID)
	return nil
}

func (r *AttendanceLogRepoImpl) GetAttendanceLogbyID(ctx context.Context, attendanceLogID int) (entity.AttendanceLog, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>attendanceLog: GetAttendanceLogbyID started for attendance log id "+strconv.Itoa(attendanceLogID), reqID)

	query := "SELECT * FROM attendance_logs WHERE attendance_logs.id=? "

	attendanceLogModel := model.AttendanceLog{}
	attendanceLogEntity := entity.AttendanceLog{}

	err := r.db.Get(&attendanceLogModel, query, attendanceLogID)
	if err != nil {
		log.Error(err.Error(), reqID)
		return attendanceLogEntity, errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
	}

	attendanceLogEntity = converter.AttendanceLogModelToAttendanceLogEntity(attendanceLogModel)

	log.Info("core>repo>attendanceLog: GetAttendanceLogbyID completed for attendance log id "+strconv.Itoa(attendanceLogID), reqID)
	return attendanceLogEntity, nil
}

func (r *AttendanceLogRepoImpl) ListAttendanceLog(ctx context.Context, filter *filters.ListFilter) (int, []entity.AttendanceLog, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>attendanceLog: ListAttendanceLog started", reqID)

	queryStatement := "SELECT * FROM attendance_logs "

	modelmap := model.AttendanceLogModelMap

	whereStr, args := filterPkg.CreateFilterStr(filter.Filters, modelmap)

	// Search attendance logs by employee ID, device name
	if filter.SearchString != "" {
		search := "%" + strings.TrimSpace(filter.SearchString) + "%"

		whereStr = append(
			whereStr,
			"(attendance_logs.emp_id LIKE ? OR attendance_logs.device_name LIKE ?)",
		)
		args = append(args, search, search)
	}

	if len(whereStr) > 0 {
		queryStatement += " WHERE " + strings.Join(whereStr, " AND ")
	}

	// Latest attendance record first
	queryStatement += " ORDER BY timestamp DESC"

	var limitQueryStmt string

	emptySortOption := filters.SortOption{}

	totalRecordQueryStatement := "SELECT COUNT(id) FROM (" + queryStatement + ") AS result"

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

	attendanceLogsModel := []model.AttendanceLog{}

	err = r.db.Select(&attendanceLogsModel, limitQueryStmt, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	attendanceLogEntities := []entity.AttendanceLog{}

	for _, attendanceLogModel := range attendanceLogsModel {
		attendanceLogEntity := converter.AttendanceLogModelToAttendanceLogEntity(attendanceLogModel)
		attendanceLogEntities = append(attendanceLogEntities, attendanceLogEntity)
	}

	log.Info("core>repo>attendanceLog: ListAttendanceLog completed", reqID)
	return count, attendanceLogEntities, nil
}

func (r *AttendanceLogRepoImpl) GetAttendanceLogDetails(ctx context.Context, selectColumns []string, table string, whereColumn []string, args []interface{}) (*entity.AttendanceLog, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>attendanceLog: GetAttendanceLogDetails started", reqID)

	selectStr := strings.Join(selectColumns, ", ")
	whereStr := strings.Join(whereColumn, " AND ")

	attendanceLogModel := model.AttendanceLog{}

	query := "SELECT " + selectStr + " FROM " + table + " WHERE " + whereStr
	err := r.db.Get(&attendanceLogModel, query, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	attendanceLogEntity := converter.AttendanceLogModelToAttendanceLogEntity(attendanceLogModel)

	log.Info("core>repo>attendanceLog: GetAttendanceLogDetails completed", reqID)
	return &attendanceLogEntity, nil
}

func (r *AttendanceLogRepoImpl) ListDailyAttendanceLog(ctx context.Context, filter *filters.ListFilter, empID, fromDate, toDate string) (int, []model.DailyAttendanceLog, errors.Response) {

	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>attendanceLog: ListDailyAttendanceLog started", reqID)

	if fromDate == "" {
		fromDate = time.Now().Format("2006-01-02")
	}

	if toDate == "" {
		toDate = fromDate
	}

	countQuery := `
		SELECT COUNT(*)
		FROM employees e
	`

	countWhere := []string{}
	countArgs := []interface{}{}

	if empID != "" {
		countWhere = append(countWhere, "e.emp_id = ?")
		countArgs = append(countArgs, empID)
	}

	if filter.SearchString != "" {
		search := "%" + strings.TrimSpace(filter.SearchString) + "%"
		countWhere = append(
			countWhere,
			"(e.emp_id LIKE ? OR e.emp_name LIKE ?)",
		)

		countArgs = append(countArgs, search, search)
	}

	if len(countWhere) > 0 {
		countQuery += " WHERE " + strings.Join(countWhere, " AND ")
	}

	var count int

	err := r.db.Get(&count, countQuery, countArgs...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	if filter.Page == 0 {
		filter.Page = 1
	}
	offset := commonConstants.NO_OF_RECORDS_PER_PAGE * (filter.Page - 1)
	query := `
		SELECT
			e.emp_id,
			e.emp_name,
			CAST(? AS DATE) AS log_date,
			al.timestamp,
			al.punch
		FROM (
			SELECT
				emp_id,
				emp_name
			FROM employees
	`

	args := []interface{}{
		fromDate,
	}

	employeeWhere := []string{}

	if empID != "" {
		employeeWhere = append(employeeWhere, "emp_id = ?")
		args = append(args, empID)
	}

	if filter.SearchString != "" {
		search := "%" + strings.TrimSpace(filter.SearchString) + "%"

		employeeWhere = append(
			employeeWhere,
			"(emp_id LIKE ? OR emp_name LIKE ?)",
		)

		args = append(args, search, search)
	}

	if len(employeeWhere) > 0 {
		query += " WHERE " + strings.Join(employeeWhere, " AND ")
	}

	query += `
			ORDER BY emp_id ASC
			LIMIT ?, ?
		) e

		LEFT JOIN attendance_logs al
			ON al.emp_id = e.emp_id
			AND DATE(al.timestamp) BETWEEN ? AND ?
			AND al.device_name IN ('Front Entry', 'Back Entry')

		ORDER BY
			e.emp_id ASC,
			al.timestamp ASC
	`

	args = append(args, offset, commonConstants.NO_OF_RECORDS_PER_PAGE, fromDate, toDate)

	var attendanceLogs []model.DailyAttendanceLog

	err = r.db.Select(&attendanceLogs, query, args...)
	if err != nil {
		log.Error(err.Error(), reqID)
		return 0, nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	log.Info("core>repo>attendanceLog: ListDailyAttendanceLog completed", reqID)
	return count, attendanceLogs, nil
}
