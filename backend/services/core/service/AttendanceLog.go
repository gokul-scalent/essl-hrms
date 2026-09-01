package service

import (
	"context"
	"strconv"
	"time"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	"github.com/scalent.io/scalent-hrms/model"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	"github.com/scalent.io/scalent-hrms/pkg/validation"
)

type AttendanceLogServiceImpl struct {
	attendanceLogRepo AttendanceLogRepo
}

func NewAttendanceLogServiceImpl(attendanceLogRepo AttendanceLogRepo) (*AttendanceLogServiceImpl, error) {
	return &AttendanceLogServiceImpl{
		attendanceLogRepo: attendanceLogRepo,
	}, nil
}

func (s *AttendanceLogServiceImpl) CreateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) (int, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>attendanceLog: create attendance log started", reqID)

	attendanceLogID, errResp := s.attendanceLogRepo.CreateAttendanceLog(ctx, attendanceLog)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}

	log.Info("core>service>attendanceLog: create attendance log completed & attendance log id is "+strconv.Itoa(attendanceLogID), reqID)
	return attendanceLogID, nil
}

func (s *AttendanceLogServiceImpl) PartialUpdateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>attendanceLog: partila update attendance log started for attendance log id "+strconv.Itoa(attendanceLog.ID), reqID)

	errResp := s.attendanceLogRepo.PartialUpdateAttendanceLog(ctx, attendanceLog)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>attendanceLog: update attendance log completed for attendance log id "+strconv.Itoa(attendanceLog.ID), reqID)
	return nil
}

func (s *AttendanceLogServiceImpl) UpdateAttendanceLog(ctx context.Context, attendanceLog entity.AttendanceLog) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>attendanceLog: update attendance log started for attendance log id "+strconv.Itoa(attendanceLog.ID), reqID)

	errResp := s.attendanceLogRepo.UpdateAttendanceLog(ctx, attendanceLog)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>attendanceLog: update attendance log completed for attendance log id "+strconv.Itoa(attendanceLog.ID), reqID)
	return nil
}

func (s *AttendanceLogServiceImpl) GetAttendanceLogbyID(ctx context.Context, attendanceLogID int) (entity.AttendanceLog, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>attendanceLog: get attendance log started for attendance log id "+strconv.Itoa(attendanceLogID), reqID)

	attendanceLog, errResp := s.attendanceLogRepo.GetAttendanceLogbyID(ctx, attendanceLogID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return entity.AttendanceLog{}, errResp
	}

	log.Info("core>service>attendanceLog: attendance log fetched successfully for attendance log id "+strconv.Itoa(attendanceLogID), reqID)
	return attendanceLog, nil
}

func (s *AttendanceLogServiceImpl) ListAttendanceLog(ctx context.Context, filter *filters.ListFilter) (int, []entity.AttendanceLog, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>attendanceLog: attendance log list started", reqID)

	totalRecords, attendanceLogsEntity, errResp := s.attendanceLogRepo.ListAttendanceLog(ctx, filter)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, nil, errResp
	}

	log.Info("core>service>attendanceLog: attendance log list completed", reqID)
	return totalRecords, attendanceLogsEntity, nil
}

func (s *AttendanceLogServiceImpl) ListDailyAttendanceLog(ctx context.Context, filter *filters.ListFilter, empID, fromDate, toDate string) (int, []entity.DailyAttendanceLog, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)

	log.Info("core>service>attendanceLog: ListDailyAttendanceLog started", reqID)

	count, attendanceLogs, errResp :=
		s.attendanceLogRepo.ListDailyAttendanceLog(ctx, filter, empID, fromDate, toDate)

	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, nil, errResp
	}

	dailyAttendance := s.CalculateDailyAttendance(attendanceLogs)

	log.Info("core>service>attendanceLog: ListDailyAttendanceLog completed", reqID)
	return count, dailyAttendance, nil
}

func (s *AttendanceLogServiceImpl) CalculateDailyAttendance(logs []model.DailyAttendanceLog) []entity.DailyAttendanceLog {

	// If there are no attendance punches, return an empty result.
	if len(logs) == 0 {
		return []entity.DailyAttendanceLog{}
	}

	// Group attendance punches by employee and date. because the database returns individual punches, daily attendance record per employee.
	grouped := make(map[string][]model.DailyAttendanceLog)

	// Loop through all attendance punches received from the repository.
	for _, log := range logs {

		// Create a unique key using employee ID and attendance date.
		key := log.EmpID + "_" + log.LogDate.Format("2006-01-02")

		// Add the punch to the employee/date group.
		grouped[key] = append(grouped[key], log)
	}

	// This will contain the final daily attendance records.
	result := []entity.DailyAttendanceLog{}

	// Process each employee/date group separately.
	for _, employeeLogs := range grouped {

		dailyLog := entity.DailyAttendanceLog{
			EmpID:   employeeLogs[0].EmpID,
			EmpName: employeeLogs[0].EmpName,
			Date:    employeeLogs[0].LogDate,
			Punches: []entity.AttendancePunch{},
		}

		// Stores the current Check-In time and When we find a Check-Out, we use this value to create a Check-In / Check-Out pair.
		var checkIn *time.Time

		// store the total working hours If an employee has multiple Check-In / Check-Out pairs, all completed durations are added together.
		var totalWorkingDuration time.Duration

		for _, log := range employeeLogs {

			// Punch 0 means Check-In.
			if log.Punch == 0 {

				// Store the Check-In timestamp.
				checkInTime := log.Timestamp
				checkIn = &checkInTime
				// Wait for the next Check-Out punch.
				continue
			}

			// Punch 1 means Check-Out.
			// We only process it when a Check-In already exists in db
			if log.Punch == 1 && checkIn != nil {

				// Store the Check-Out timestamp.
				checkOutTime := log.Timestamp

				// Create a Check-In / Check-Out pair.
				dailyLog.Punches = append(
					dailyLog.Punches,
					entity.AttendancePunch{
						CheckIn:  checkIn,
						CheckOut: &checkOutTime,
					},
				)
				totalWorkingDuration += checkOutTime.Sub(*checkIn)
				// Reset Check-In so that the next Check-In can start a new attendance pair.
				checkIn = nil
			}
		}

		// If Check-In exists here, it means the employee has checked in but has not checked out yet.
		if checkIn != nil {

			// Add the incomplete Check-In to the punches list.
			dailyLog.Punches = append(
				dailyLog.Punches,
				entity.AttendancePunch{
					CheckIn:  checkIn,
					CheckOut: nil,
				},
			)

			// Employee has an incomplete attendance record.
			dailyLog.Status = "INCOMPLETE"

		} else if len(dailyLog.Punches) > 0 {
			// At least one complete Check-In / Check-Out pair exists.
			dailyLog.Status = "PRESENT"

		} else {
			// No valid attendance pair was found.
			dailyLog.Status = "ABSENT"
		}

		dailyLog.WorkingHours = validation.FormatWorkingHours(totalWorkingDuration)
		result = append(result, dailyLog)
	}
	return result
}
