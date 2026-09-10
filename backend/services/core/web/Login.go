package web

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scalent.io/scalent-hrms/apimodel"
	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
	scalentContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	"github.com/scalent.io/scalent-hrms/pkg/utils"
	httpUtils "github.com/scalent.io/scalent-hrms/pkg/utils"
	"github.com/scalent.io/scalent-hrms/pkg/validation"
)

func (h CoreHandlerRegistry) LoginHandler(c *gin.Context) {
	reqID, _ := scalentContext.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>login: LoginHandler started", reqID)

	loginRequest := &coreAPIModel.LoginRequest{}

	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, loginRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	// Normalize email
	loginRequest.Email = strings.ToLower(strings.TrimSpace(loginRequest.Email))

	if loginRequest.Email == "" {
		validationErrors := []apimodel.InvalidValidationError{
			{
				Field: "email",
				Msg:   "Email is required",
			},
		}
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	if loginRequest.Password == "" {
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError("Password is required"),
			nil,
		)
		return
	}

	decodedPassword, decodeErr := utils.DecodeBase64Password(loginRequest.Password)
	if decodeErr != nil {
		log.Error("password decode failed: "+decodeErr.Error(), reqID)
		httpUtils.ErrorResponse(
			c,
			errors.ResponseBadRequestError("Invalid password format"),
			nil,
		)
		return
	}

	if len(decodedPassword) > 16 {
		validationErrors := []apimodel.InvalidValidationError{
			{
				Field: "password",
				Msg:   "Password must not exceed 16 characters",
			},
		}
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	// Login
	userEntity, token, errResp :=
		h.Options.LoginService.Login(
			c.Request.Context(),
			loginRequest.Email,
			decodedPassword,
		)

	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	// Convert roles to role codes
	roleCodes := make([]string, 0, len(userEntity.Roles))

	for _, role := range userEntity.Roles {
		if role.Code != "" {
			roleCodes = append(roleCodes, role.Code)
		}
	}

	// Make sure the user has at least one role
	if len(roleCodes) == 0 {
		log.Error("login successful but user has no active roles", reqID)

		httpUtils.ErrorResponse(c, errors.ResponseUnauthorizedError("User has no active role assigned"), nil)
		return
	}

	response := coreAPIModel.LoginResponse{
		Email:         userEntity.Email,
		Roles:         roleCodes,
		Token:         token,
		IsPasswordSet: userEntity.IsPasswordSet,
	}

	log.Info("core>web>login: LoginHandler completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "Login successful.", response)
}

func (h *CoreHandlerRegistry) LogOutHandler(c *gin.Context) {
	reqID, _ := scalentContext.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>login: LogOutHandler started", reqID)

	errResp := h.Options.LoginService.LogOut(c.Request.Context())
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>login: LogOutHandler completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "Logout successful.", nil)
}
