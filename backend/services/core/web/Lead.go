package web

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

func (h CoreHandlerRegistry) CreateLeadHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>lead: create lead started", reqID)

	leadRequest := &coreAPIModel.CreateLeadRequest{}
	if err := c.ShouldBindWith(leadRequest, binding.FormMultipart); err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError(err.Error()),
			nil,
		)
		return
	}

	switch leadRequest.Type {

	case commonConstants.UPLOAD_TYPE_SINGLE:
		if strings.TrimSpace(leadRequest.Email) == "" {
			httpUtils.ErrorResponse(
				c,
				errors.ResponseBadRequestError("email is required"),
				nil,
			)
			return
		}

	case commonConstants.UPLOAD_TYPE_CSV:
		if leadRequest.File == nil {
			httpUtils.ErrorResponse(
				c,
				errors.ResponseBadRequestError("csv file is required"),
				nil,
			)
			return
		}

	default:
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError("invalid upload type"),
			nil,
		)
		return
	}

	leadEntity := converter.CreateLeadAPIToLeadEntity(leadRequest)

	_, errResp := h.Options.LeadService.CreateLead(
		c.Request.Context(),
		leadEntity,
		leadRequest,
	)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>lead: create lead completed", reqID)

	httpUtils.DataResponse(
		c,
		http.StatusOK,
		"Lead(s) created successfully",
		nil,
	)
}

func (h CoreHandlerRegistry) PartialUpdateLeadHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>lead: partial update lead started", reqID)

	leadIDStr := c.Param("id")
	leadIDStr = strings.TrimSpace(leadIDStr)

	leadID, err := strconv.Atoi(leadIDStr)
	if err != nil {
		log.Error("invalid lead id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid lead id"), nil)
		return
	}

	log.Info("core>web>lead: partial lead update started for lead id "+leadIDStr, reqID)

	leadRequest := &coreAPIModel.UpdateLeadRequest{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, leadRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	leadEntity := converter.UpdateLeadAPIRequestToLeadEntity(leadRequest)
	leadEntity.ID = leadID

	errResp := h.Options.LeadService.PartialUpdateLead(c.Request.Context(), leadEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>lead: partial update lead completed successfully for lead id "+leadIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "partial update lead completed successfully", nil)
}

func (h CoreHandlerRegistry) UpdateLeadHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>lead: update lead started ", reqID)

	leadIDStr := c.Param("id")
	leadIDStr = strings.TrimSpace(leadIDStr)

	leadID, err := strconv.Atoi(leadIDStr)
	if err != nil {
		log.Error("invalid lead id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid lead id"), nil)
		return
	}

	log.Info("core>web>lead: lead update started for lead id "+leadIDStr, reqID)

	leadRequest := &coreAPIModel.Lead{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, leadRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	leadEntity := converter.LeadAPIToLeadEntity(leadRequest)
	leadEntity.ID = leadID

	errResp := h.Options.LeadService.UpdateLead(c.Request.Context(), leadEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>lead: update lead completed successfully for lead id "+leadIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "web: update lead completed successfully", nil)
}

func (h CoreHandlerRegistry) DeleteLeadHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>lead: delete lead started", reqID)

	leadIDstr := c.Param("id")
	leadIDstr = strings.TrimSpace(leadIDstr)

	leadID, err := strconv.Atoi(leadIDstr)
	if err != nil {
		log.Error("invalid lead id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid lead id"), nil)
		return
	}

	log.Info("core>web>lead: delete lead started for lead id "+strconv.Itoa(leadID), reqID)

	errResp := h.Options.LeadService.DeleteLead(c.Request.Context(), leadID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>lead: delete lead completed for lead id "+strconv.Itoa(leadID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "lead deleted successfully", nil)
}

func (h CoreHandlerRegistry) GetLeadbyIDHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>lead: get lead started", reqID)

	leadIDstr := c.Param("id")
	leadIDstr = strings.TrimSpace(leadIDstr)

	leadID, err := strconv.Atoi(leadIDstr)
	if err != nil {
		log.Error("invalid lead id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid lead id"), nil)
		return
	}

	log.Info("core>web>lead:get lead started for lead id "+leadIDstr, reqID)

	lead, errResp := h.Options.LeadService.GetLeadbyID(c.Request.Context(), leadID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	leadResponse := converter.LeadEntityToLeadAPIModelResponse(lead)

	log.Info("core>web>lead: lead fetched successfully for lead id "+strconv.Itoa(leadID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "lead fetched successfully", leadResponse)
}

func (h CoreHandlerRegistry) ListLeadHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>lead: lead list started", reqID)

	listFiltersRequest := apimodel.ListFiltersRequest{}

	queryParams := c.Request.URL.Query()

	paramPageValue := queryParams.Get("page")
	pageNo, _ := strconv.Atoi(paramPageValue)
	paramSearchString := queryParams.Get("searchString") //added search string params

	paramEmailListID := queryParams.Get("emailListID")
	emailListID, err := strconv.Atoi(paramEmailListID)
	if err != nil || emailListID <= 0 {
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError("invalid emailListID"),
			nil,
		)
		return
	}

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

	totalRecords, leadsEntity, statusCount, errResp := h.Options.LeadService.ListLead(c.Request.Context(), emailListID, &filterEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	leadsListResponse := []coreAPIModel.LeadResponse{}

	for _, leadEntity := range leadsEntity {
		leadresponse := converter.LeadEntityToLeadAPIModelResponse(leadEntity)
		leadsListResponse = append(leadsListResponse, leadresponse)
	}

	var leadResponse coreAPIModel.LeadListResponse
	leadResponse.TotalRecords = totalRecords
	leadResponse.NoOfRecordsPerPage = commonConstants.NO_OF_RECORDS_PER_PAGE
	leadResponse.Lead = leadsListResponse
	//to get the count of all status summary in the list
	leadResponse.SafeCount = statusCount.SafeCount
	leadResponse.RiskyCount = statusCount.RiskyCount
	leadResponse.InvalidCount = statusCount.InvalidCount
	leadResponse.UnknownCount = statusCount.UnknownCount
	leadResponse.PendingCount = statusCount.PendingCount
	leadResponse.TimeoutCount = statusCount.TimeoutCount

	log.Info("core>web>lead: lead list completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "leads fetched successfully", leadResponse)
}

func (h CoreHandlerRegistry) DownloadLeadsCSVHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>lead>download csv reached", reqID)
	//get the status query or params
	status := strings.TrimSpace(c.Query("status"))
	emailListID := strings.TrimSpace(c.Query("emailListID"))

	log.Info("core>lead>download csv status: "+status, reqID)

	// Validate required email list id
	if emailListID == "" {
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError("emailListID is required"),
			nil,
		)
		return
	}

	log.Info("core>lead>download csv started", reqID)
	// Convert emailListID from string to integer
	emailListIDInt, err := strconv.Atoi(emailListID)
	if err != nil {
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError("invalid emailListID"),
			nil,
		)
		return
	}

	// If status is empty -> all leads are returned.
	// Otherwise only leads matching the given status are returned.
	leads, errResp := h.Options.LeadService.GetSafeLeads(
		c.Request.Context(),
		emailListIDInt,
		status,
	)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	// Set response headers so browser downloads the response as a CSV file.
	c.Header("Content-Type", "text/csv")
	c.Header("Cache-Control", "no-cache")

	// Default filename when downloading all records.
	filename := "all-leads.csv"
	// give file name according to status: safe-leads.csv, risky-leads.csv
	if status != "" {
		filename = status + "-leads.csv"
	}
	c.Header("Content-Disposition", "attachment; filename="+filename)
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	if err := writer.Write([]string{
		"email",
		"first_name",
		"last_name",
		"job_title",
		"company",
		"city",
		"country",
		"industry",
	}); err != nil {
		log.Error("csv header write error: "+err.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseInternalServerError("failed to generate csv"),
			nil,
		)
		return
	}

	// Write each lead as a CSV row
	for _, lead := range leads {
		if err := writer.Write([]string{
			lead.Email,
			lead.FirstName,
			lead.LastName,
			lead.JobTitle,
			lead.Company,
			lead.City,
			lead.Country,
			lead.Industry,
		}); err != nil {
			log.Error("csv row write error: "+err.Error(), reqID)
			httpUtils.ErrorResponse(
				c,
				errors.ResponseInternalServerError("failed to write csv row"),
				nil,
			)
			return
		}
	}

	// Flush any buffered data to the response
	writer.Flush()

	// Check if any error occurred while flushing
	if err := writer.Error(); err != nil {
		log.Error("csv flush error: "+err.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseInternalServerError("failed to generate csv"),
			nil,
		)
		return
	}

	log.Info("core>lead>download csv completed", reqID)
}

func (h CoreHandlerRegistry) ReverifyLeadHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>lead: reverify lead started", reqID)

	leadIDStr := strings.TrimSpace(c.Param("id"))

	leadID, err := strconv.Atoi(leadIDStr)
	if err != nil {
		log.Error("invalid lead id", reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError("invalid lead id"),
			nil,
		)
		return
	}

	log.Info("core>web>lead: reverify started for lead id "+leadIDStr, reqID)

	errResp := h.Options.LeadService.ReverifyLead(c.Request.Context(), leadID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>lead: reverify completed for lead id "+leadIDStr, reqID)

	httpUtils.DataResponse(
		c,
		http.StatusOK,
		"Lead queued for re-verification successfully",
		nil,
	)
}
