package web

import (
	"io"

	"github.com/scalent.io/scalent-hrms/internal/middleware"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	"github.com/scalent.io/scalent-hrms/pkg/server"
	coreService "github.com/scalent.io/scalent-hrms/services/core/service"

	"github.com/gin-gonic/gin"
)

type CoreHandlerRegistryOptions struct {
	Config          *server.Config
	Middleware      middleware.Middleware
	HomeService     coreService.HomeService
	UserService     coreService.UserService
	LoginService    coreService.LoginService
	EmployeeService coreService.EmployeeService
}

type CoreHandlerRegistry struct {
	Options CoreHandlerRegistryOptions
}

func NewCoreHandlerRegistry(options CoreHandlerRegistryOptions) *CoreHandlerRegistry {
	return &CoreHandlerRegistry{
		Options: options,
	}
}

func (h *CoreHandlerRegistry) StartServer() error {

	router, err := h.registerRoutes()
	if err != nil {
		log.Print(err)
	}

	log.Info("Server Started Successfully", "")
	router.Run(h.Options.Config.Port)
	return nil
}

func (h CoreHandlerRegistry) registerRoutes() (*gin.Engine, error) {

	router := gin.Default()
	//redirect gin logs to io.discard which goes bydefault to terminal
	gin.DefaultWriter = io.Discard
	router.Use(h.Options.Middleware.Cors())

	middleware.InitRateLimiterRedis("localhost:8081")

	coreRouter := router.Group("/scalent-hrms")
	coreRouter.POST("/login", h.LoginHandler)
	coreRouter.GET("/home", h.HomeHandler)

	coreRouter.Use(h.Options.Middleware.Access())
	coreRouter.POST("/logout", h.LogOutHandler)

	userRouter := coreRouter.Group("/user")
	userRouter.POST("/", h.CreateUserHandler)
	userRouter.PATCH("/:id", h.PartialUpdateUserHandler)
	userRouter.PUT("/:id", h.UpdateUserHandler)
	userRouter.DELETE("/:id", h.DeleteUserHandler)
	userRouter.GET("/:id", h.GetUserbyIDHandler)
	userRouter.GET("/list", h.ListUserHandler)
	userRouter.PATCH("/change-password", h.ChangePasswordHandler)

	//---add the following line above in CoreHandlerRegistryOptions struct
	//---UserService  coreService.UserService
	//---delete these lines after copy

	employeeRouter := coreRouter.Group("/employee")
	employeeRouter.POST("/", h.CreateEmployeeHandler)
	employeeRouter.PATCH("/:id", h.PartialUpdateEmployeeHandler)
	employeeRouter.PUT("/:id", h.UpdateEmployeeHandler)
	employeeRouter.DELETE("/:id", h.DeleteEmployeeHandler)
	employeeRouter.GET("/:id", h.GetEmployeebyIDHandler)
	employeeRouter.GET("/list", h.ListEmployeeHandler)

	//---add the following line above in CoreHandlerRegistryOptions struct
	//---EmployeeService  coreService.EmployeeService
	//---delete these lines after copy

	//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
	return router, nil
}
