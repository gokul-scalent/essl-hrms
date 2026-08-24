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

	repo.NewEmailListRepoImpl,
	wire.Bind(new(service.EmailListRepo), new(*repo.EmailListRepoImpl)),

	service.NewEmailListServiceImpl,
	wire.Bind(new(service.EmailListService), new(*service.EmailListServiceImpl)),

	repo.NewLeadRepoImpl,
	wire.Bind(new(service.LeadRepo), new(*repo.LeadRepoImpl)),

	service.NewLeadServiceImpl,
	wire.Bind(new(service.LeadService), new(*service.LeadServiceImpl)),

	//login module
	repo.NewLoginRepoImpl,
	wire.Bind(new(service.LoginRepo), new(*repo.LoginRepoImpl)),

	service.NewLoginServiceImpl,
	wire.Bind(new(service.LoginService), new(*service.LoginServiceImpl)),

	repo.NewUserSettingRepoImpl,
	wire.Bind(new(service.UserSettingRepo), new(*repo.UserSettingRepoImpl)),

	service.NewUserSettingServiceImpl,
	wire.Bind(new(service.UserSettingService), new(*service.UserSettingServiceImpl)),

// -----==-----==DO NOT ADD CODE BELOW THIS LINE------
)

func initServer(config *CoreConfig) (*web.CoreHandlerRegistry, error) {
	panic(wire.Build(CoreModuleSet))
}
