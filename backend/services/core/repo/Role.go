package repo

import (
	"context"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/internal/converter"
	"github.com/scalent.io/scalent-hrms/model"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
)

type RoleRepoImpl struct {
	db *sqlx.DB
}

func NewRoleRepoImpl(db *sqlx.DB) (*RoleRepoImpl, error) {
	return &RoleRepoImpl{
		db: db,
	}, nil
}

func (r *RoleRepoImpl) GetRoles(ctx context.Context) ([]entity.Role, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>role: GetRoles started", reqID)

	query := `
		SELECT
			id,
			name,
			code,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM roles
		WHERE deleted_at IS NULL
		ORDER BY id ASC
	`

	roleModels := []model.Role{}
	err := r.db.SelectContext(ctx, &roleModels, query)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	roles := make([]entity.Role, 0, len(roleModels))

	for _, roleModel := range roleModels {
		roles = append(roles, converter.RoleModelToRoleEntity(roleModel))
	}

	log.Info("core>repo>role: GetRoles completed", reqID)
	return roles, nil
}

func (r *RoleRepoImpl) GetRoleByID(ctx context.Context, id int) (*entity.Role, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>role: GetRoleByID started for role id "+strconv.Itoa(id), reqID)

	query := `
		SELECT
			id,
			name,
			code,
			status,
			created_at,
			updated_at,
			deleted_at
		FROM roles
		WHERE id = ?
		  AND deleted_at IS NULL
		LIMIT 1
	`

	roleModel := model.Role{}

	err := r.db.GetContext(ctx, &roleModel, query, id)
	if err != nil {
		log.Error(err.Error(), reqID)
		return nil, errors.ResponseNotFoundError(errors.NOT_FOUND_ERROR)
	}

	roleEntity := converter.RoleModelToRoleEntity(roleModel)

	log.Info("core>repo>role: GetRoleByID completed for role id "+strconv.Itoa(id), reqID)
	return &roleEntity, nil
}
