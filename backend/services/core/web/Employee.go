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

func (h CoreHandlerRegistry) CreateEmployeeHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>employee: create employee started", reqID)

	employeeRequest := &coreAPIModel.Employee{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, employeeRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	employeeEntity := converter.CreateEmployeeAPIRequestToEmployeeEntity(employeeRequest)

	employeeID, errResp := h.Options.EmployeeService.CreateEmployee(c.Request.Context(), employeeEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>employee: create employee completed & employee id is "+strconv.Itoa(employeeID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "Employee created successfully", nil)
}

func (h CoreHandlerRegistry) PartialUpdateEmployeeHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>employee: partial update employee started", reqID)

	employeeIDStr := c.Param("id")
	employeeIDStr = strings.TrimSpace(employeeIDStr)

	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		log.Error("invalid employee id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid employee id"), nil)
		return
	}

	log.Info("core>web>employee: partial employee update started for employee id "+employeeIDStr, reqID)

	employeeRequest := &coreAPIModel.UpdateEmployeeRequest{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, employeeRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	employeeEntity := converter.UpdateEmployeeAPIRequestToEmployeeEntity(employeeRequest)
	employeeEntity.ID = employeeID

	errResp := h.Options.EmployeeService.PartialUpdateEmployee(c.Request.Context(), employeeEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>employee: partial update employee completed successfully for employee id "+employeeIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "partial update employee completed successfully", nil)
}

func (h CoreHandlerRegistry) UpdateEmployeeHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>employee: update employee started ", reqID)

	employeeIDStr := c.Param("id")
	employeeIDStr = strings.TrimSpace(employeeIDStr)

	employeeID, err := strconv.Atoi(employeeIDStr)
	if err != nil {
		log.Error("invalid employee id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid employee id"), nil)
		return
	}

	log.Info("core>web>employee: employee update started for employee id "+employeeIDStr, reqID)

	employeeRequest := &coreAPIModel.UpdateEmployeeRequest{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, employeeRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	employeeEntity := converter.UpdateEmployeeAPIRequestToEmployeeEntity(employeeRequest)
	employeeEntity.ID = employeeID

	errResp := h.Options.EmployeeService.UpdateEmployee(c.Request.Context(), employeeEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>employee: update employee completed successfully for employee id "+employeeIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "web: update employee completed successfully", nil)
}

func (h CoreHandlerRegistry) DeleteEmployeeHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>employee: delete employee started", reqID)

	employeeIDstr := c.Param("id")
	employeeIDstr = strings.TrimSpace(employeeIDstr)

	employeeID, err := strconv.Atoi(employeeIDstr)
	if err != nil {
		log.Error("invalid employee id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid employee id"), nil)
		return
	}

	log.Info("core>web>employee: delete employee started for employee id "+strconv.Itoa(employeeID), reqID)

	errResp := h.Options.EmployeeService.DeleteEmployee(c.Request.Context(), employeeID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>employee: delete employee completed for employee id "+strconv.Itoa(employeeID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "employee deleted successfully", nil)
}

func (h CoreHandlerRegistry) GetEmployeebyIDHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>employee: get employee started", reqID)

	employeeIDstr := c.Param("id")
	employeeIDstr = strings.TrimSpace(employeeIDstr)

	employeeID, err := strconv.Atoi(employeeIDstr)
	if err != nil {
		log.Error("invalid employee id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid employee id"), nil)
		return
	}

	log.Info("core>web>employee:get employee started for employee id "+employeeIDstr, reqID)

	employee, errResp := h.Options.EmployeeService.GetEmployeebyID(c.Request.Context(), employeeID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	employeeResponse := converter.EmployeeEntityToUserAPIModelResponse(employee)

	log.Info("core>web>employee: employee fetched successfully for employee id "+strconv.Itoa(employeeID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "employee fetched successfully", employeeResponse)
}

func (h CoreHandlerRegistry) ListEmployeeHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>employee: employee list started", reqID)

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

	totalRecords, employeesEntity, errResp := h.Options.EmployeeService.ListEmployee(c.Request.Context(), &filterEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	employeesListResponse := []coreAPIModel.EmployeeResponse{}

	for _, employeeEntity := range employeesEntity {
		employeeresponse := converter.EmployeeEntityToUserAPIModelResponse(employeeEntity)
		employeesListResponse = append(employeesListResponse, employeeresponse)
	}

	var employeeResponse coreAPIModel.EmployeeListResponse
	employeeResponse.TotalRecords = totalRecords
	employeeResponse.NoOfRecordsPerPage = commonConstants.NO_OF_RECORDS_PER_PAGE
	employeeResponse.Employee = employeesListResponse

	log.Info("core>web>employee: employee list completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "employees fetched successfully", employeeResponse)
}
