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

func (h CoreHandlerRegistry) CreateUserSettingHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>userSetting: create user setting started", reqID)

	userSettingRequest := &coreAPIModel.UserSetting{}

	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, userSettingRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}
	// Get session
	session, sessionErr := context.GetSessionFromContext(c.Request.Context())
	if sessionErr != nil {
		log.Error(sessionErr.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseUnauthorizedError("Unauthorized"),
			nil,
		)
		return
	}

	userSettingEntity := converter.UserSettingAPIToUserSettingEntity(userSettingRequest)
	// Set logged-in user id
	userSettingEntity.User.ID = session.UserID

	userSettingID, errResp := h.Options.UserSettingService.CreateUserSetting(
		c.Request.Context(),
		userSettingEntity,
	)

	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info(
		"core>web>userSetting: create user setting completed & user setting id is "+strconv.Itoa(userSettingID),
		reqID,
	)

	httpUtils.DataResponse(c, http.StatusOK, "Verification internal saved successfully", nil)
}

func (h CoreHandlerRegistry) PartialUpdateUserSettingHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>userSetting: partial update user setting started", reqID)

	userSettingIDStr := c.Param("id")
	userSettingIDStr = strings.TrimSpace(userSettingIDStr)

	userSettingID, err := strconv.Atoi(userSettingIDStr)
	if err != nil {
		log.Error("invalid userSetting id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid userSetting id"), nil)
		return
	}

	log.Info("core>web>userSetting: partial user setting update started for user setting id "+userSettingIDStr, reqID)

	userSettingRequest := &coreAPIModel.UpdateUserSettingRequest{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, userSettingRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}
	//get the user id from session
	session, sessionErr := context.GetSessionFromContext(c.Request.Context())
	if sessionErr != nil {
		log.Error(sessionErr.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseUnauthorizedError("Unauthorized"),
			nil,
		)
		return
	}

	userSettingEntity := converter.UpdateUserSettingAPIRequestToUserSettingEntity(userSettingRequest)
	userSettingEntity.ID = userSettingID
	userSettingEntity.User.ID = session.UserID

	errResp := h.Options.UserSettingService.PartialUpdateUserSetting(c.Request.Context(), userSettingEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>userSetting: partial update user setting completed successfully for user setting id "+userSettingIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "partial update user setting completed successfully", nil)
}

func (h CoreHandlerRegistry) UpdateUserSettingHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>userSetting: update user setting started ", reqID)

	userSettingIDStr := c.Param("id")
	userSettingIDStr = strings.TrimSpace(userSettingIDStr)

	userSettingID, err := strconv.Atoi(userSettingIDStr)
	if err != nil {
		log.Error("invalid userSetting id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid userSetting id"), nil)
		return
	}

	log.Info("core>web>userSetting: user setting update started for user setting id "+userSettingIDStr, reqID)

	userSettingRequest := &coreAPIModel.UserSetting{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, userSettingRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	userSettingEntity := converter.UserSettingAPIToUserSettingEntity(userSettingRequest)
	userSettingEntity.ID = userSettingID

	errResp := h.Options.UserSettingService.UpdateUserSetting(c.Request.Context(), userSettingEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>userSetting: update user setting completed successfully for user setting id "+userSettingIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "web: update user setting completed successfully", nil)
}

func (h CoreHandlerRegistry) DeleteUserSettingHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>userSetting: delete user setting started", reqID)

	userSettingIDstr := c.Param("id")
	userSettingIDstr = strings.TrimSpace(userSettingIDstr)

	userSettingID, err := strconv.Atoi(userSettingIDstr)
	if err != nil {
		log.Error("invalid userSetting id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid userSetting id"), nil)
		return
	}

	log.Info("core>web>userSetting: delete user setting started for user setting id "+strconv.Itoa(userSettingID), reqID)

	errResp := h.Options.UserSettingService.DeleteUserSetting(c.Request.Context(), userSettingID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>userSetting: delete user setting completed for user setting id "+strconv.Itoa(userSettingID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "user setting deleted successfully", nil)
}

func (h CoreHandlerRegistry) GetUserSettingbyIDHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>userSetting: get user setting started", reqID)

	userSettingIDstr := c.Param("id")
	userSettingIDstr = strings.TrimSpace(userSettingIDstr)

	userSettingID, err := strconv.Atoi(userSettingIDstr)
	if err != nil {
		log.Error("invalid userSetting id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid userSetting id"), nil)
		return
	}

	log.Info("core>web>userSetting:get user setting started for user setting id "+userSettingIDstr, reqID)

	userSetting, errResp := h.Options.UserSettingService.GetUserSettingbyID(c.Request.Context(), userSettingID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	userSettingResponse := converter.UserSettingEntityToUserSettingAPIModelResponse(userSetting)

	log.Info("core>web>userSetting: user setting fetched successfully for user setting id "+strconv.Itoa(userSettingID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "user setting fetched successfully", userSettingResponse)
}

func (h CoreHandlerRegistry) ListUserSettingHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>userSetting: user setting list started", reqID)

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

	totalRecords, userSettingsEntity, errResp := h.Options.UserSettingService.ListUserSetting(c.Request.Context(), &filterEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	userSettingsListResponse := []coreAPIModel.UserSettingResponse{}

	for _, userSettingEntity := range userSettingsEntity {
		userSettingresponse := converter.UserSettingEntityToUserSettingAPIModelResponse(userSettingEntity)
		userSettingsListResponse = append(userSettingsListResponse, userSettingresponse)
	}

	var userSettingResponse coreAPIModel.UserSettingListResponse
	userSettingResponse.TotalRecords = totalRecords
	userSettingResponse.NoOfRecordsPerPage = commonConstants.NO_OF_RECORDS_PER_PAGE
	userSettingResponse.UserSetting = userSettingsListResponse

	log.Info("core>web>userSetting: user setting list completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "user settings fetched successfully", userSettingResponse)
}
