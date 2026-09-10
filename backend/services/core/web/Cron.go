package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	httpUtils "github.com/scalent.io/scalent-hrms/pkg/utils"
)

func (h CoreHandlerRegistry) CalculateWorkingHoursHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>cron: calculate working hours started", reqID)

	var targetDate time.Time
	dateParam := c.Query("date")
	if dateParam != "" {
		parsedDate, err := time.ParseInLocation("2006-01-02", dateParam, time.Local)
		if err != nil {
			httpUtils.DataResponse(c, http.StatusBadRequest, "Invalid date format. Use YYYY-MM-DD", nil)
			return
		}
		targetDate = parsedDate
	}

	workingHoursSummary, errResp := h.Options.CronService.CalculateWorkingHours(c.Request.Context(), targetDate)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>cron: calculate working hours completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "Working hours calculated successfully", workingHoursSummary)
}
