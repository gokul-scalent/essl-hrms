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
)

func (h CoreHandlerRegistry) CreateEmailListHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("Inside CreateEmailListHandler", reqID)
	log.Info("core>web>emailList: CreateEmailList started", reqID)

	emailListRequest := &coreAPIModel.CreateEmailListRequest{}
	if err := c.ShouldBindJSON(emailListRequest); err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError(err.Error()),
			nil,
		)
		return
	}

	emailListEntity := converter.CreateEmailListAPIToEmailListEntity(emailListRequest)

	emailListID, errResp := h.Options.EmailListService.CreateEmailList(
		c.Request.Context(),
		emailListEntity,
		emailListRequest,
	)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>emailList: CreateEmailList completed & email list id is "+strconv.Itoa(emailListID), reqID)

	httpUtils.DataResponse(
		c,
		http.StatusOK,
		"Email list created successfully",
		nil,
	)
}
func (h CoreHandlerRegistry) PartialUpdateEmailListHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>emailList: partial update email list started", reqID)

	emailListIDStr := c.Param("id")
	emailListIDStr = strings.TrimSpace(emailListIDStr)

	emailListID, err := strconv.Atoi(emailListIDStr)
	if err != nil {
		log.Error("invalid emailList id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid emailList id"), nil)
		return
	}

	log.Info("core>web>emailList: partial email list update started for email list id "+emailListIDStr, reqID)

	emailListRequest := &coreAPIModel.UpdateEmailListRequest{}
	// validationErrors, err := validation.DecodeAndValidate(c.Request.Body, emailListRequest, c)
	if err := c.ShouldBind(emailListRequest); err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError(err.Error()),
			nil,
		)
		return
	}

	emailListEntity := converter.UpdateEmailListAPIRequestToEmailListEntity(emailListRequest)
	emailListEntity.ID = emailListID

	errResp := h.Options.EmailListService.PartialUpdateEmailList(c.Request.Context(), emailListEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>emailList: partial update email list completed successfully for email list id "+emailListIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "partial update email list completed successfully", nil)
}

func (h CoreHandlerRegistry) UpdateEmailListHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>emailList: update email list started ", reqID)

	emailListIDStr := c.Param("id")
	emailListIDStr = strings.TrimSpace(emailListIDStr)

	emailListID, err := strconv.Atoi(emailListIDStr)
	if err != nil {
		log.Error("invalid emailList id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid emailList id"), nil)
		return
	}

	log.Info("core>web>emailList: email list update started for email list id "+emailListIDStr, reqID)

	emailListRequest := &coreAPIModel.UpdateEmailListRequest{}

	if err := c.ShouldBind(emailListRequest); err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError(err.Error()),
			nil,
		)
		return
	}

	emailListEntity := converter.UpdateEmailListAPIRequestToEmailListEntity(emailListRequest)
	emailListEntity.ID = emailListID

	errResp := h.Options.EmailListService.UpdateEmailList(c.Request.Context(), emailListEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>emailList: update email list completed successfully for email list id "+emailListIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "web: update email list completed successfully", nil)
}

func (h CoreHandlerRegistry) DeleteEmailListHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>emailList: delete email list started", reqID)

	emailListIDstr := c.Param("id")
	emailListIDstr = strings.TrimSpace(emailListIDstr)

	emailListID, err := strconv.Atoi(emailListIDstr)
	if err != nil {
		log.Error("invalid emailList id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid emailList id"), nil)
		return
	}

	log.Info("core>web>emailList: delete email list started for email list id "+strconv.Itoa(emailListID), reqID)

	errResp := h.Options.EmailListService.DeleteEmailList(c.Request.Context(), emailListID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>emailList: delete email list completed for email list id "+strconv.Itoa(emailListID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "email list deleted successfully", nil)
}

func (h CoreHandlerRegistry) GetEmailListbyIDHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>emailList: get email list started", reqID)

	emailListIDstr := c.Param("id")
	emailListIDstr = strings.TrimSpace(emailListIDstr)

	emailListID, err := strconv.Atoi(emailListIDstr)
	if err != nil {
		log.Error("invalid emailList id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid emailList id"), nil)
		return
	}
	// Get logged-in user/session
	session, err := context.GetSessionFromContext(c.Request.Context())
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseUnauthorizedError(errors.SESSION_ERROR),
			nil,
		)
		return
	}

	userID := session.UserID

	log.Info("core>web>emailList: get email list started for email list id "+emailListIDstr+" and user id "+strconv.Itoa(userID), reqID)

	emailList, errResp := h.Options.EmailListService.GetEmailListbyID(c.Request.Context(), emailListID, userID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	emailListResponse := converter.EmailListEntityToEmailListAPIModelResponse(emailList)

	log.Info("core>web>emailList: email list fetched successfully for email list id "+strconv.Itoa(emailListID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "email list fetched successfully", emailListResponse)
}

func (h CoreHandlerRegistry) ListEmailListHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>emailList: email list list started", reqID)

	listFiltersRequest := apimodel.ListFiltersRequest{}

	queryParams := c.Request.URL.Query()

	paramPageValue := queryParams.Get("page")
	pageNo, _ := strconv.Atoi(paramPageValue)
	paramSearchString := queryParams.Get("searchString") //added search string params

	paramFilterVAlue := queryParams.Get("filtersJSON")
	var filtersArray []apimodel.Filter
	json.Unmarshal([]byte(paramFilterVAlue), &filtersArray)

	paramSortOptionVAlue := queryParams.Get("sortOptionJSON")
	var sortOption apimodel.SortOption
	json.Unmarshal([]byte(paramSortOptionVAlue), &sortOption)

	listFiltersRequest.Page = pageNo
	listFiltersRequest.Filters = filtersArray
	listFiltersRequest.SortOption = sortOption
	listFiltersRequest.SearchString = paramSearchString //added search string params

	filterEntity := converter.FilterAPIRequestToFilterEntity(listFiltersRequest)

	// Get logged-in user/session
	session, err := context.GetSessionFromContext(c.Request.Context())
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseUnauthorizedError(errors.SESSION_ERROR),
			nil,
		)
		return
	}

	userID := session.UserID

	totalRecords, emailListsEntity, errResp := h.Options.EmailListService.ListEmailList(c.Request.Context(), &filterEntity, userID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	emailListsListResponse := []coreAPIModel.EmailListResponse{}

	for _, emailListEntity := range emailListsEntity {
		emailListresponse := converter.EmailListEntityToEmailListAPIModelResponse(emailListEntity)
		emailListsListResponse = append(emailListsListResponse, emailListresponse)
	}

	var emailListResponse coreAPIModel.EmailListListResponse
	emailListResponse.TotalRecords = totalRecords
	emailListResponse.NoOfRecordsPerPage = commonConstants.NO_OF_RECORDS_PER_PAGE
	emailListResponse.EmailList = emailListsListResponse

	log.Info("core>web>emailList: email list list completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "email lists fetched successfully", emailListResponse)
}
