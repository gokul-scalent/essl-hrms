package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	httpUtils "github.com/scalent.io/scalent-hrms/pkg/utils"
)

func (h CoreHandlerRegistry) CreateUserOnDeviceHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>User on device: create User on device started ", reqID)

	httpUtils.DataResponse(c, http.StatusOK, "User created successfully", nil)
}
