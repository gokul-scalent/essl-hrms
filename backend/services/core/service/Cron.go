package service

import (
	"context"
	"fmt"
	"time"

	"github.com/scalent.io/scalent-hrms/entity"
	hrmsContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type CronServiceImpl struct {
	cronRepo CronRepo
}

func NewCronServiceImpl(cronRepo CronRepo) (*CronServiceImpl, error) {
	return &CronServiceImpl{
		cronRepo: cronRepo,
	}, nil
}
func (h *CronServiceImpl) CalculateWorkingHours(ctx context.Context, targetDate time.Time) (
	summaries []entity.WorkingHoursSummary, errResp errors.Response) {

	reqID, _ := hrmsContext.GetRequestIDFromContext(ctx)

	log.Info("core>service>cron: calculate working hours started", reqID)

	// ============================================================
	// TIMEZONE
	// ============================================================

	location, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return nil, errors.ResponseInternalServerError(fmt.Sprintf("failed to load timezone: %v", err))
	}

	// ============================================================
	// TARGET DATE
	// ============================================================

	// Keep hardcoded for now as requested.
	targetDate = time.Date(2026, 9, 2, 0, 0, 0, 0, location)

	fromDate := targetDate

	toDate := targetDate.AddDate(0, 0, 1)

	// Office end.
	officeEnd := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 19, 0, 0, 0, location)

	// ============================================================
	// EMPLOYEES
	// ============================================================

	employees, errResp := h.cronRepo.ListActiveEmployees(ctx)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return nil, errResp
	}

	// ============================================================
	// FETCH ATTENDANCE
	// ============================================================

	attendanceLogs, errResp := h.cronRepo.GetAttendanceLogsByDateRange(ctx, fromDate, toDate)

	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return nil, errResp
	}

	log.Info(fmt.Sprintf("core>service>cron: employees=%d attendanceLogs=%d", len(employees), len(attendanceLogs)), reqID)

	// ============================================================
	// GROUP ORIGINAL ATTENDANCE BY EMPLOYEE
	// ============================================================

	employeeLogs := make(map[string][]entity.AttendanceLog)

	for _, attendanceLog := range attendanceLogs {
		// Very important:
		//
		// Do NOT use previously synthesized records
		// when running synthesis again.
		if attendanceLog.Synthesized {
			continue
		}
		employeeLogs[attendanceLog.EmpID] = append(employeeLogs[attendanceLog.EmpID], attendanceLog)
	}

	// ============================================================
	// SYNTHESIS
	// ============================================================

	allSyntheticLogs := make([]entity.AttendanceLog, 0)
	for _, employee := range employees {
		logs := employeeLogs[employee.EmpID]
		if len(logs) == 0 {
			continue
		}

		syntheticLogs := synthesizeAttendanceLogs(logs, officeEnd)

		if len(syntheticLogs) == 0 {
			continue
		}

		log.Info(fmt.Sprintf("core>service>cron: employee=%s synthetic records=%d", employee.EmpID, len(syntheticLogs)), reqID)

		allSyntheticLogs = append(allSyntheticLogs, syntheticLogs...)
	}

	// ============================================================
	// INSERT SYNTHETIC RECORDS
	// ============================================================

	if len(allSyntheticLogs) > 0 {

		errResp = h.cronRepo.InsertSynthesizedAttendanceLogs(ctx, allSyntheticLogs)

		if errResp != nil {
			log.Error(errResp.Error(), reqID)
			return nil, errResp
		}

		log.Info(fmt.Sprintf("core>service>cron: inserted %d synthetic attendance records", len(allSyntheticLogs)), reqID)
	}

	// ============================================================
	// FETCH AGAIN
	//
	// Now database contains:
	//
	// ORIGINAL + SYNTHESIZED
	// ============================================================

	attendanceLogs, errResp = h.cronRepo.GetAttendanceLogsByDateRange(ctx, fromDate, toDate)

	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return nil, errResp
	}

	// ============================================================
	// GROUP COMPLETE ATTENDANCE
	// ============================================================

	employeeLogs = make(map[string][]entity.AttendanceLog)

	for _, attendanceLog := range attendanceLogs {

		employeeLogs[attendanceLog.EmpID] = append(employeeLogs[attendanceLog.EmpID], attendanceLog)
	}

	// ============================================================
	// CALCULATE
	// ============================================================

	summaries = make([]entity.WorkingHoursSummary, 0, len(employees))

	for _, employee := range employees {

		logs := employeeLogs[employee.EmpID]

		if len(logs) == 0 {

			summaries = append(
				summaries,
				entity.WorkingHoursSummary{
					EmpID:            employee.EmpID,
					EmpName:          employee.EmpName,
					Date:             targetDate.Format("2006-01-02"),
					Status:           "ABSENT",
					WorkingHours:     "00:00",
					OutOfOfficeHours: "00:00",
				},
			)

			continue
		}

		// Validate normalized sequence.
		if validationErr := validateAttendanceSequence(logs); validationErr != nil {

			log.Error(fmt.Sprintf("core>service>cron: invalid attendance sequence for employee %s: %v", employee.EmpID, validationErr), reqID)

			// Don't silently calculate bad attendance.
			// You can alternatively continue here if preferred.
			continue
		}

		summary := calculateWorkingHoursFromLogs(targetDate, employee, logs)

		summaries = append(summaries, summary)
	}

	// ============================================================
	// SAVE SUMMARY
	// ============================================================

	errResp = h.cronRepo.UpsertWorkingHoursSummaries(ctx, summaries)

	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return nil, errResp
	}

	log.Info("core>service>cron: calculate working hours completed", reqID)

	return summaries, nil
}
