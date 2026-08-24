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
	"github.com/scalent.io/scalent-hrms/pkg/utils"
	httpUtils "github.com/scalent.io/scalent-hrms/pkg/utils"
	"github.com/scalent.io/scalent-hrms/pkg/validation"
)

func (h CoreHandlerRegistry) CreateUserHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>user: create user started", reqID)

	userRequest := &coreAPIModel.CreateUser{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, userRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	userEntity := converter.CreateUserAPIToUserEntity(userRequest)

	userID, errResp := h.Options.UserService.CreateUser(c.Request.Context(), userEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>user: create user completed & user id is "+strconv.Itoa(userID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "User created successfully", nil)
}

func (h CoreHandlerRegistry) PartialUpdateUserHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>user: partial update user started", reqID)

	userIDStr := c.Param("id")
	userIDStr = strings.TrimSpace(userIDStr)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Error("invalid user id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid user id"), nil)
		return
	}

	log.Info("core>web>user: partial user update started for user id "+userIDStr, reqID)

	userRequest := &coreAPIModel.UpdateUserRequest{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, userRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	userEntity := converter.UpdateUserAPIRequestToUserEntity(userRequest)
	userEntity.ID = userID

	errResp := h.Options.UserService.PartialUpdateUser(c.Request.Context(), userEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>user: partial update user completed successfully for user id "+userIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "partial update user completed successfully", nil)
}

func (h CoreHandlerRegistry) UpdateUserHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>user: update user started ", reqID)

	userIDStr := c.Param("id")
	userIDStr = strings.TrimSpace(userIDStr)

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Error("invalid user id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid user id"), nil)
		return
	}

	log.Info("core>web>user: user update started for user id "+userIDStr, reqID)

	userRequest := &coreAPIModel.User{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, userRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	userEntity := converter.UserAPIToUserEntity(userRequest)
	userEntity.ID = userID

	errResp := h.Options.UserService.UpdateUser(c.Request.Context(), userEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>user: update user completed successfully for user id "+userIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "web: update user completed successfully", nil)
}

func (h CoreHandlerRegistry) DeleteUserHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>user: delete user started", reqID)

	userIDstr := c.Param("id")
	userIDstr = strings.TrimSpace(userIDstr)

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		log.Error("invalid user id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid user id"), nil)
		return
	}

	log.Info("core>web>user: delete user started for user id "+strconv.Itoa(userID), reqID)

	errResp := h.Options.UserService.DeleteUser(c.Request.Context(), userID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>user: delete user completed for user id "+strconv.Itoa(userID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "user deleted successfully", nil)
}

func (h CoreHandlerRegistry) GetUserbyIDHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>user: get user started", reqID)

	userIDstr := c.Param("id")
	userIDstr = strings.TrimSpace(userIDstr)

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		log.Error("invalid user id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid user id"), nil)
		return
	}

	log.Info("core>web>user:get user started for user id "+userIDstr, reqID)

	user, errResp := h.Options.UserService.GetUserbyID(c.Request.Context(), userID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	userResponse := converter.UserEntityToUserAPIModelResponse(user)

	log.Info("core>web>user: user fetched successfully for user id "+strconv.Itoa(userID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "user fetched successfully", userResponse)
}

func (h CoreHandlerRegistry) ListUserHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>user: user list started", reqID)

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

	totalRecords, usersEntity, errResp := h.Options.UserService.ListUser(c.Request.Context(), &filterEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	usersListResponse := []coreAPIModel.UserResponse{}

	for _, userEntity := range usersEntity {
		userresponse := converter.UserEntityToUserAPIModelResponse(userEntity)
		usersListResponse = append(usersListResponse, userresponse)
	}

	var userResponse coreAPIModel.UserListResponse
	userResponse.TotalRecords = totalRecords
	userResponse.NoOfRecordsPerPage = commonConstants.NO_OF_RECORDS_PER_PAGE
	userResponse.User = usersListResponse

	log.Info("core>web>user: user list completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "users fetched successfully", userResponse)
}

func (h CoreHandlerRegistry) ChangePasswordHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>user: ChangePasswordHandler started", reqID)

	userRequest := &coreAPIModel.ChangePasswordRequest{}

	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, userRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	oldDecodedPassword, errStd := utils.DecodeBase64Password(userRequest.OldPassword)
	if errStd != nil {
		log.Error("Error decoding base64 old password: "+errStd.Error(), reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("Invalid password format"), nil)
		return
	}

	if len(oldDecodedPassword) > 16 {
		log.Error("core>web>login: Old Password length exceeded (max 16 characters allowed)", reqID)
		validationErrors = append(validationErrors, apimodel.InvalidValidationError{
			Field: "password",
			Msg:   "Old Password must not exceed 16 characters.",
		})
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	newDecodedPassword, errStd := utils.DecodeBase64Password(userRequest.NewPassword)
	if errStd != nil {
		log.Error("Error decoding base64 old password: "+errStd.Error(), reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("Invalid password format"), nil)
		return
	}

	if len(newDecodedPassword) > 16 {
		log.Error("core>web>login: New Password length exceeded (max 16 characters allowed)", reqID)
		validationErrors = append(validationErrors, apimodel.InvalidValidationError{
			Field: "password",
			Msg:   "New Password must not exceed 16 characters.",
		})
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	confirmDecodedPassword, errStd := utils.DecodeBase64Password(userRequest.ConfirmPassword)
	if errStd != nil {
		log.Error("Error decoding base64 confirm password: "+errStd.Error(), reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("Invalid password format"), nil)
		return
	}

	if len(confirmDecodedPassword) > 16 {
		log.Error("core>web>login: Confirm Password length exceeded (max 16 characters allowed)", reqID)
		validationErrors = append(validationErrors, apimodel.InvalidValidationError{
			Field: "password",
			Msg:   "Confirm Password must not exceed 16 characters.",
		})
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	if newDecodedPassword != confirmDecodedPassword {
		log.Error("new password and confirm password do not match", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("New Password and Confirm Password do not match"), nil)
		return
	}

	oldPassword := strings.TrimSpace(oldDecodedPassword)
	newPassword := strings.TrimSpace(newDecodedPassword)

	errResp := h.Options.UserService.ChangePassword(c.Request.Context(), oldPassword, newPassword)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	errResp = h.Options.LoginService.LogOut(c.Request.Context())
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>user: ChangePasswordHandler completed successfully ", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "Password changed successfully.", nil)
}
