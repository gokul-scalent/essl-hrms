package casbin

import (
	sqlxadapter "github.com/Blank-Xu/sqlx-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/jmoiron/sqlx"
)

type CasbinImplConfig struct {
	Casbin *casbin.Enforcer
	Db     *sqlx.DB
}

func InitCasbin(Db *sqlx.DB) *casbin.Enforcer {
	if err := Db.Ping(); err != nil {
		panic(err)
	}

	// Initialize the adapter (table name: casbin_policy)
	a, err := sqlxadapter.NewAdapter(Db, "casbin_policy")
	if err != nil {
		panic(err)
	}

	// Define the Casbin model directly in code
	modelText := `
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = r.sub == p.sub && keyMatch4(r.obj, p.obj) && (r.act == p.act || p.act == "*")
	`

	// Create model from string
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		panic(err)
	}

	// Create the enforcer using the model and adapter
	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		panic(err)
	}

	if err := enforcer.LoadPolicy(); err != nil {
		panic(err)
	}

	return enforcer
}
