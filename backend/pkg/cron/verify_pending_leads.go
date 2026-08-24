package cron

import (
	"context"
	"time"

	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	coreService "github.com/scalent.io/scalent-hrms/services/core/service"
)

var leadService coreService.LeadService

func Init(ls coreService.LeadService) {
	leadService = ls
}

func VerifyPendingLeads() (time.Duration, errors.Response) {
	if leadService == nil {
		log.Error("LeadService not initialized", "")
		return 0, errors.ResponseInternalServerError(
			errors.INTERNAL_SERVER_ERROR,
		)
	}

	return leadService.VerifyPendingLeads(context.Background())
}
