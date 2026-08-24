package repo

import (
	"context"

	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type HomeRepoImpl struct{}

func NewHomeRepoImpl() (*HomeRepoImpl, error) {
	return &HomeRepoImpl{}, nil
}

func (r *HomeRepoImpl) Home(ctx context.Context) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>home: Home started", reqID)

	log.Info("core>repo>home: Home completed", reqID)
	return nil
}
