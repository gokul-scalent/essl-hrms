package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scalent.io/scalent-hrms/apimodel"
	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
	commonConstants "github.com/scalent.io/scalent-hrms/entity/commonConstants"
	"github.com/scalent.io/scalent-hrms/internal/converter"
	"github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	httpUtils "github.com/scalent.io/scalent-hrms/pkg/utils"
	"github.com/scalent.io/scalent-hrms/pkg/validation"
)

func (h CoreHandlerRegistry) CreateAttendanceLogHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>attendanceLog: create attendance log started", reqID)

	attendanceLogRequest := &coreAPIModel.AttendanceLog{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, attendanceLogRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	attendanceLogEntity := converter.CreateAttendanceLogAPIToAttendanceLogEntity(attendanceLogRequest)

	attendanceLogID, errResp := h.Options.AttendanceLogService.CreateAttendanceLog(c.Request.Context(), attendanceLogEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>attendanceLog: create attendance log completed & attendance log id is "+strconv.Itoa(attendanceLogID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "AttendanceLog created successfully", nil)
}

func (h CoreHandlerRegistry) PartialUpdateAttendanceLogHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>attendanceLog: partial update attendance log started", reqID)

	attendanceLogIDStr := c.Param("id")
	attendanceLogIDStr = strings.TrimSpace(attendanceLogIDStr)

	attendanceLogID, err := strconv.Atoi(attendanceLogIDStr)
	if err != nil {
		log.Error("invalid attendanceLog id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid attendanceLog id"), nil)
		return
	}

	log.Info("core>web>attendanceLog: partial attendance log update started for attendance log id "+attendanceLogIDStr, reqID)

	attendanceLogRequest := &coreAPIModel.UpdateAttendanceLogRequest{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, attendanceLogRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	attendanceLogEntity := converter.UpdateAttendanceLogAPIRequestToAttendanceLogEntity(attendanceLogRequest)
	attendanceLogEntity.ID = attendanceLogID

	errResp := h.Options.AttendanceLogService.PartialUpdateAttendanceLog(c.Request.Context(), attendanceLogEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>attendanceLog: partial update attendance log completed successfully for attendance log id "+attendanceLogIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "partial update attendance log completed successfully", nil)
}

func (h CoreHandlerRegistry) UpdateAttendanceLogHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>attendanceLog: update attendance log started ", reqID)

	attendanceLogIDStr := c.Param("id")
	attendanceLogIDStr = strings.TrimSpace(attendanceLogIDStr)

	attendanceLogID, err := strconv.Atoi(attendanceLogIDStr)
	if err != nil {
		log.Error("invalid attendanceLog id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid attendanceLog id"), nil)
		return
	}

	log.Info("core>web>attendanceLog: attendance log update started for attendance log id "+attendanceLogIDStr, reqID)

	attendanceLogRequest := &coreAPIModel.AttendanceLog{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, attendanceLogRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	attendanceLogEntity := converter.CreateAttendanceLogAPIToAttendanceLogEntity(attendanceLogRequest)
	attendanceLogEntity.ID = attendanceLogID

	errResp := h.Options.AttendanceLogService.UpdateAttendanceLog(c.Request.Context(), attendanceLogEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>attendanceLog: update attendance log completed successfully for attendance log id "+attendanceLogIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "web: update attendance log completed successfully", nil)
}

func (h CoreHandlerRegistry) GetAttendanceLogbyIDHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>attendanceLog: get attendance log started", reqID)

	attendanceLogIDstr := c.Param("id")
	attendanceLogIDstr = strings.TrimSpace(attendanceLogIDstr)

	attendanceLogID, err := strconv.Atoi(attendanceLogIDstr)
	if err != nil {
		log.Error("invalid attendanceLog id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid attendanceLog id"), nil)
		return
	}

	log.Info("core>web>attendanceLog:get attendance log started for attendance log id "+attendanceLogIDstr, reqID)

	attendanceLog, errResp := h.Options.AttendanceLogService.GetAttendanceLogbyID(c.Request.Context(), attendanceLogID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	attendanceLogResponse := converter.AttendanceLogEntityToAttendanceLogAPIModelResponse(attendanceLog)

	log.Info("core>web>attendanceLog: attendance log fetched successfully for attendance log id "+strconv.Itoa(attendanceLogID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "attendance log fetched successfully", attendanceLogResponse)
}

func (h CoreHandlerRegistry) ListAttendanceLogHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>attendanceLog: attendance log list started", reqID)

	listFiltersRequest := apimodel.ListFiltersRequest{}

	queryParams := c.Request.URL.Query()

	paramPageValue := queryParams.Get("page")
	pageNo, _ := strconv.Atoi(paramPageValue)

	paramFilterVAlue := queryParams.Get("filtersJSON")
	var filtersArray []apimodel.Filter
	json.Unmarshal([]byte(paramFilterVAlue), &filtersArray)

	paramSortOptionVAlue := queryParams.Get("sortOptionJSON")
	var sortOption apimodel.SortOption
	json.Unmarshal([]byte(paramSortOptionVAlue), &sortOption)

	listFiltersRequest.Page = pageNo
	listFiltersRequest.Filters = filtersArray
	listFiltersRequest.SortOption = sortOption

	filterEntity := converter.FilterAPIRequestToFilterEntity(listFiltersRequest)

	totalRecords, attendanceLogsEntity, errResp := h.Options.AttendanceLogService.ListAttendanceLog(c.Request.Context(), &filterEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	attendanceLogsListResponse := []coreAPIModel.AttendanceLogResponse{}

	for _, attendanceLogEntity := range attendanceLogsEntity {
		attendanceLogresponse := converter.AttendanceLogEntityToAttendanceLogAPIModelResponse(attendanceLogEntity)
		attendanceLogsListResponse = append(attendanceLogsListResponse, attendanceLogresponse)
	}

	var attendanceLogResponse coreAPIModel.AttendanceLogListResponse
	attendanceLogResponse.TotalRecords = totalRecords
	attendanceLogResponse.NoOfRecordsPerPage = commonConstants.NO_OF_RECORDS_PER_PAGE
	attendanceLogResponse.AttendanceLog = attendanceLogsListResponse

	log.Info("core>web>attendanceLog: attendance log list completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "attendance logs fetched successfully", attendanceLogResponse)
}
