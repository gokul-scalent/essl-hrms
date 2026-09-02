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

	// Group attendance punches by employee and date.
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

	// Process each employee/date group.
	for _, employeeLogs := range grouped {

		dailyLog := entity.DailyAttendanceLog{
			EmpID:   employeeLogs[0].EmpID,
			EmpName: employeeLogs[0].EmpName,
			Date:    employeeLogs[0].LogDate,
			Punches: []entity.AttendancePunch{},
		}

		var totalWorkingDuration time.Duration
		var previousPunch *model.DailyAttendanceLog

		// Tracks whether the current/last check-out was already paired.
		lastPunchPaired := false
		for i := range employeeLogs {
			current := employeeLogs[i]
			// Ignore invalid punch.
			if !current.Punch.Valid {
				continue
			}
			// Ignore invalid timestamp.
			if !current.Timestamp.Valid {
				continue
			}
			punch := int(current.Punch.Int64)
			timestamp := current.Timestamp.Time

			if previousPunch == nil {

				// CASE 1: Employee forgot the very first check-in. If the first punch is CHECK_OUT, assume CHECK_IN at 10:00 AM.
				if punch == 1 {

					checkIn := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 10, 0, 0, 0, timestamp.Location())
					checkOut := timestamp

					// Only create the pair if the actual check-out happened after the assumed 10:00 AM check-in.
					if checkOut.After(checkIn) {
						dailyLog.Punches = append(
							dailyLog.Punches,
							entity.AttendancePunch{
								CheckIn:  &checkIn,
								CheckOut: &checkOut,
							},
						)

						duration := checkOut.Sub(checkIn)
						if duration > 0 {
							totalWorkingDuration += duration
						}
					}

					previousPunch = nil
					lastPunchPaired = true

					continue
				}
				// First punch is CHECK_IN. Keep it unmatched until a CHECK_OUT is received.
				if punch == 0 {

					currentCopy := current
					previousPunch = &currentCopy
					lastPunchPaired = false
				}

				continue
			}

			previousPunchType := int(previousPunch.Punch.Int64)
			previousTime := previousPunch.Timestamp.Time

			// CASE 2:normal check-in and check-out pair
			if previousPunchType == 0 && punch == 1 {

				checkIn := previousTime
				checkOut := timestamp

				dailyLog.Punches = append(
					dailyLog.Punches,
					entity.AttendancePunch{
						CheckIn:  &checkIn,
						CheckOut: &checkOut,
					},
				)

				// Calculate working duration.
				duration := checkOut.Sub(checkIn)
				if duration > 0 {
					totalWorkingDuration += duration
				}
				// Both punches are paired.
				previousPunch = nil
				lastPunchPaired = true
				continue
			}

			// CASE 3: check-in -> check-in condition ----> First check-in has no check-out. Assume check-out = second check-in - 15 minutes.
			if previousPunchType == 0 && punch == 0 {
				checkIn := previousTime
				checkOut := timestamp.Add(-15 * time.Minute)

				if checkOut.After(checkIn) {

					dailyLog.Punches = append(
						dailyLog.Punches,
						entity.AttendancePunch{
							CheckIn:  &checkIn,
							CheckOut: &checkOut,
						},
					)

					duration := checkOut.Sub(checkIn)
					if duration > 0 {
						totalWorkingDuration += duration
					}
				}
				// Current IN becomes the new unmatched punch.
				currentCopy := current
				previousPunch = &currentCopy
				lastPunchPaired = false

				continue
			}

			// CASE 4: check-out -> check-out Missing check-in before second check-out. Assume check-in = previous check-out + 15 minutes.
			if previousPunchType == 1 && punch == 1 {

				checkIn := previousTime.Add(15 * time.Minute)
				checkOut := timestamp

				if checkIn.Before(checkOut) {
					dailyLog.Punches = append(
						dailyLog.Punches,
						entity.AttendancePunch{
							CheckIn:  &checkIn,
							CheckOut: &checkOut,
						},
					)

					duration := checkOut.Sub(checkIn)

					if duration > 0 {
						totalWorkingDuration += duration
					}
					// Current OUT has been used as checkout.
					lastPunchPaired = true
				} else {
					lastPunchPaired = false
				}

				// Current OUT becomes previous punch.
				currentCopy := current
				previousPunch = &currentCopy

				continue
			}
		}

		if previousPunch != nil {

			// Make sure punch and timestamp are valid.
			if previousPunch.Punch.Valid && previousPunch.Timestamp.Valid {

				punch := int(previousPunch.Punch.Int64)
				timestamp := previousPunch.Timestamp.Time

				// CASE 5: Last punch is check-in. Employee forgot the final check-out. Assume check-out = 7:00 PM.
				if punch == 0 {

					checkIn := timestamp

					dailyLog.Punches = append(
						dailyLog.Punches,
						entity.AttendancePunch{
							CheckIn:  &checkIn,
							CheckOut: nil,
						},
					)

					dailyLog.Status = "INCOMPLETE"
				}

				// Last punch is CHECK_OUT and was not paired. Keep CHECK_IN as NULL.
				if punch == 1 && !lastPunchPaired {

					checkOut := timestamp

					dailyLog.Punches = append(
						dailyLog.Punches,
						entity.AttendancePunch{
							CheckIn:  nil,
							CheckOut: &checkOut,
						},
					)

					dailyLog.Status = "INCOMPLETE"
				}
			}
		}
		if dailyLog.Status == "" {

			if len(dailyLog.Punches) > 0 {
				dailyLog.Status = "PRESENT"
			} else {
				dailyLog.Status = "ABSENT"
			}
		}
		dailyLog.WorkingHours = validation.FormatWorkingHours(totalWorkingDuration)

		result = append(result, dailyLog)
	}
	return result
}
