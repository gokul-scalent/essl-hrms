//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/scalent.io/scalent-hrms/internal/auth"
	"github.com/scalent.io/scalent-hrms/internal/middleware"
	"github.com/scalent.io/scalent-hrms/pkg/casbin"
	"github.com/scalent.io/scalent-hrms/pkg/db/cache"
	"github.com/scalent.io/scalent-hrms/pkg/db/sqlx"
	"github.com/scalent.io/scalent-hrms/pkg/reacher"
	"github.com/scalent.io/scalent-hrms/services/core/repo"
	"github.com/scalent.io/scalent-hrms/services/core/service"
	"github.com/scalent.io/scalent-hrms/services/core/web"
)

var CoreModuleSet = wire.NewSet(
	wire.FieldsOf(new(*CoreConfig), "server", "db", "Reacher"),
	NewCacheConfig,
	cache.NewRedisInstance,
	casbin.InitCasbin,
	sqlx.NewDBConn,
	NewServiceConfig,
	reacher.NewReacherClient,

	wire.Struct(new(web.CoreHandlerRegistryOptions), "*"),
	web.NewCoreHandlerRegistry,

	auth.NewAuthImpl,

	middleware.NewMiddlewareImpl,
	wire.Bind(new(middleware.Middleware), new(*middleware.MiddlewareImpl)),

	repo.NewHomeRepoImpl,
	wire.Bind(new(service.HomeRepo), new(*repo.HomeRepoImpl)),

	service.NewHomeServiceImpl,
	wire.Bind(new(service.HomeService), new(*service.HomeServiceImpl)),

	repo.NewUserRepoImpl,
	wire.Bind(new(service.UserRepo), new(*repo.UserRepoImpl)),

	service.NewUserServiceImpl,
	wire.Bind(new(service.UserService), new(*service.UserServiceImpl)),

	//login module
	repo.NewLoginRepoImpl,
	wire.Bind(new(service.LoginRepo), new(*repo.LoginRepoImpl)),

	service.NewLoginServiceImpl,
	wire.Bind(new(service.LoginService), new(*service.LoginServiceImpl)),

	repo.NewEmployeeRepoImpl,
	wire.Bind(new(service.EmployeeRepo), new(*repo.EmployeeRepoImpl)),

	service.NewEmployeeServiceImpl,
	wire.Bind(new(service.EmployeeService), new(*service.EmployeeServiceImpl)),

// -----==-----==DO NOT ADD CODE BELOW THIS LINE------
)

func initServer(config *CoreConfig) (*web.CoreHandlerRegistry, error) {
	panic(wire.Build(CoreModuleSet))
}
