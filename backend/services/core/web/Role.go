package web

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	coreAPIModel "github.com/scalent.io/scalent-hrms/apimodel/core"
	"github.com/scalent.io/scalent-hrms/internal/converter"
	"github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	httpUtils "github.com/scalent.io/scalent-hrms/pkg/utils"
)

func (h CoreHandlerRegistry) GetRolesHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>role: get roles started", reqID)

	roles, errResp := h.Options.RoleService.GetRoles(c.Request.Context())

	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	roleResponses := make([]coreAPIModel.RoleResponse, 0, len(roles))

	for _, role := range roles {
		roleResponses = append(roleResponses, converter.RoleEntityToRoleAPIModelResponse(role))
	}

	log.Info("core>web>role: roles fetched successfully", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "Roles fetched successfully", roleResponses)
}

func (h CoreHandlerRegistry) GetRoleByIDHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>role: get role by id started", reqID)

	roleIDStr := strings.TrimSpace(c.Param("id"))

	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		log.Error("invalid role id", reqID)

		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid role id"), nil)
		return
	}

	log.Info("core>web>role: get role started for role id "+roleIDStr, reqID)
	role, errResp := h.Options.RoleService.GetRoleByID(
		c.Request.Context(),
		roleID,
	)

	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	roleResponse := converter.RoleEntityToRoleAPIModelResponse(*role)
	log.Info("core>web>role: role fetched successfully for role id "+roleIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "Role fetched successfully", roleResponse)
}
