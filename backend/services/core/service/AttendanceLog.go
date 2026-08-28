package service

import (
	"context"
	"strconv"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/filters"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
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
