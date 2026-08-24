package web

import (
	"net/http"

	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	httpUtils "github.com/scalent.io/scalent-hrms/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h CoreHandlerRegistry) HomeHandler(c *gin.Context) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>home: HomeHandler started", reqID)

	log.Info("core>web>login: HomeHandler completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "Home successful.", nil)
}
