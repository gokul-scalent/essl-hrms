package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/scalent.io/scalent-hrms/internal/converter"
	hrmsContext "github.com/scalent.io/scalent-hrms/pkg/context"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/model"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type CronRepoImpl struct {
	db *sqlx.DB
}

func NewCronRepoImpl(db *sqlx.DB) (*CronRepoImpl, error) {
	return &CronRepoImpl{
		db: db,
	}, nil
}

func (r *CronRepoImpl) ListActiveEmployees(ctx context.Context) ([]entity.Employee, errors.Response) {
	reqID, _ := hrmsContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>cron: list active employees started", reqID)
	defer log.Info("core>repo>cron: list active employees completed", reqID)

	query := "SELECT * FROM employees WHERE deleted_at IS NULL"
	employeesModel := []model.Employee{}

	err := r.db.Select(&employeesModel, query)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	employees := make([]entity.Employee, 0, len(employeesModel))
	for _, employeeModel := range employeesModel {
		employees = append(employees, converter.EmployeeModelToEmployeeEntity(employeeModel))
	}

	return employees, nil
}

func (r *CronRepoImpl) GetAttendanceLogsByDateRange(ctx context.Context, fromDate time.Time, toDate time.Time) ([]entity.AttendanceLog, errors.Response) {
	reqID, _ := hrmsContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>cron: get attendance logs by date range started", reqID)
	defer log.Info("core>repo>cron: get attendance logs by date range completed", reqID)

	//query := "SELECT * FROM attendance_logs WHERE DATE(attendance_logs.timestamp)= ?  AND( emp_id=? OR emp_id=? ) order by emp_id, timestamp asc"
	query := "SELECT * FROM attendance_logs WHERE DATE(attendance_logs.timestamp)= ? order by emp_id, timestamp asc"
	attendanceLogModel := []model.AttendanceLog{}
	dateStr := fromDate.Format("2006-01-02")

	log.Info(fmt.Sprintf("Executing query: %s | args: %s", query, dateStr), reqID)

	err := r.db.Select(&attendanceLogModel, query, dateStr)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	attendanceLogEntity := converter.AllAttendanceLogModelToAttendanceLogEntity(attendanceLogModel)
	log.Info(fmt.Sprintf("Query result count: %d", len(attendanceLogEntity)), reqID)
	return attendanceLogEntity, nil

}

func (r *CronRepoImpl) UpsertWorkingHoursSummaries(ctx context.Context, summaries []entity.WorkingHoursSummary) errors.Response {
	reqID, _ := hrmsContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>cron: upsert working hours summaries started", reqID)
	defer log.Info("core>repo>cron: upsert working hours summaries completed", reqID)

	if len(summaries) == 0 {
		return nil
	}

	// createTableQuery := `
	// 	CREATE TABLE IF NOT EXISTS attendance_working_hours (
	// 		id BIGINT NOT NULL AUTO_INCREMENT,
	// 		emp_id VARCHAR(64) NOT NULL,
	// 		log_date DATE NOT NULL,
	// 		check_in_time DATETIME NULL,
	// 		check_out_time DATETIME NULL,
	// 		working_hours VARCHAR(8) NOT NULL,
	// 		out_of_office_hours VARCHAR(8) NOT NULL,
	// 		status ENUM('PRESENT','ABSENT') NOT NULL DEFAULT 'ABSENT',
	// 		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	// 		updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	// 		PRIMARY KEY (id),
	// 		UNIQUE KEY uq_attendance_working_hours_emp_date (emp_id, log_date),
	// 		KEY idx_attendance_working_hours_date (log_date)
	// 	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
	// `

	// if _, err := r.db.Exec(createTableQuery); err != nil {
	// 	log.Error(err.Error(), reqID)
	// 	return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	// }

	query := `
		INSERT INTO attendance_working_hours
			(emp_id, log_date, check_in_time, check_out_time, working_hours, out_of_office_hours, status)
		VALUES
			(?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			check_in_time = VALUES(check_in_time),
			check_out_time = VALUES(check_out_time),
			working_hours = VALUES(working_hours),
			out_of_office_hours = VALUES(out_of_office_hours),
			status = VALUES(status),
			updated_at = CURRENT_TIMESTAMP
	`

	for _, summary := range summaries {
		var checkIn interface{}
		if summary.CheckInTime != nil {
			checkIn = *summary.CheckInTime
		}

		var checkOut interface{}
		if summary.CheckOutTime != nil {
			checkOut = *summary.CheckOutTime
		}

		if _, err := r.db.Exec(
			query,
			summary.EmpID,
			summary.Date,
			checkIn,
			checkOut,
			summary.WorkingHours,
			summary.OutOfOfficeHours,
			summary.Status,
		); err != nil {
			log.Error(err.Error(), reqID)
			return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
		}
	}

	return nil

}
func (r *CronRepoImpl) InsertSynthesizedAttendanceLogs(
	ctx context.Context,
	logs []entity.AttendanceLog,
) errors.Response {

	if len(logs) == 0 {
		return nil
	}

	query := `
		INSERT IGNORE INTO attendance_logs (
			uid,
			emp_id,
			timestamp,
			status,
			punch,
			attendance_state,
			device_name,
			created_at,
			synthesized
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.ResponseInternalServerError(err.Error())
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return errors.ResponseInternalServerError(err.Error())
	}
	defer stmt.Close()

	for _, attendanceLog := range logs {

		_, err = stmt.ExecContext(
			ctx,
			attendanceLog.UID,
			attendanceLog.EmpID,
			attendanceLog.Timestamp,
			attendanceLog.Status,
			attendanceLog.Punch,
			attendanceLog.AttendanceState,
			attendanceLog.DeviceName,
			attendanceLog.CreatedAt,
			attendanceLog.Synthesized,
		)

		if err != nil {
			_ = tx.Rollback()

			return errors.ResponseInternalServerError(
				fmt.Sprintf(
					"failed to insert synthesized attendance log for employee %s: %v",
					attendanceLog.EmpID,
					err,
				),
			)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.ResponseInternalServerError(err.Error())
	}

	return nil
}
