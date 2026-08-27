package repo

import (
	"context"
	"strconv"
	"strings"

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

	if len(whereStr) > 0 {
		whereString := strings.Join(whereStr, " AND ")
		queryStatement += " WHERE " + whereString
	}

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
