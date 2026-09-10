package web

import (
	"context"
	"io"
	"time"

	"github.com/scalent.io/scalent-hrms/internal/middleware"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	"github.com/scalent.io/scalent-hrms/pkg/server"
	coreService "github.com/scalent.io/scalent-hrms/services/core/service"

	"github.com/gin-gonic/gin"
)

type CoreHandlerRegistryOptions struct {
	Config               *server.Config
	Middleware           middleware.Middleware
	HomeService          coreService.HomeService
	UserService          coreService.UserService
	LoginService         coreService.LoginService
	EmployeeService      coreService.EmployeeService
	AttendanceLogService coreService.AttendanceLogService
	CronService          coreService.CronService
	RoleService          coreService.RoleService
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

	h.startDailyWorkingHoursCron()

	log.Info("Server Started Successfully", "")
	router.Run(h.Options.Config.Port)
	return nil
}

func (h *CoreHandlerRegistry) startDailyWorkingHoursCron() {
	go func() {
		for {
			now := time.Now()
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

			time.Sleep(time.Until(nextMidnight))

			targetDate := nextMidnight.AddDate(0, 0, -1)
			_, errResp := h.Options.CronService.CalculateWorkingHours(context.Background(), targetDate)
			if errResp != nil {
				log.Error("core>web>cron scheduler: "+errResp.Error(), "")
				continue
			}

			log.Info("core>web>cron scheduler: working-hours summary calculated for "+targetDate.Format("2006-01-02"), "")
		}
	}()
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

	// Keep cron endpoint public (no authorization middleware).
	cronRouter := coreRouter.Group("/cron")
	cronRouter.POST("/calculate-working-hours", h.CalculateWorkingHoursHandler)
	coreRouter.Use(h.Options.Middleware.Access())

	roleRouter := coreRouter.Group("/role")
	roleRouter.GET("/list", h.GetRolesHandler)
	roleRouter.GET("/:id", h.GetRoleByIDHandler)

	coreRouter.POST("/logout", h.LogOutHandler)

	protectedRouter := coreRouter.Group("")
	protectedRouter.Use(h.Options.Middleware.Access())
	protectedRouter.POST("/logout", h.LogOutHandler)

	userRouter := protectedRouter.Group("/user")
	userRouter.POST("/", h.CreateUserHandler)
	userRouter.PATCH("/:id", h.PartialUpdateUserHandler)
	userRouter.PUT("/:id", h.UpdateUserHandler)
	userRouter.DELETE("/:id", h.DeleteUserHandler)
	userRouter.GET("/:id", h.GetUserbyIDHandler)
	userRouter.GET("/list", h.ListUserHandler)
	userRouter.PATCH("/change-password", h.ChangePasswordHandler)
	userRouter.POST("/:id/send-mail", h.SendUserMailHandler)

	//---add the following line above in CoreHandlerRegistryOptions struct
	//---UserService  coreService.UserService
	//---delete these lines after copy

	employeeRouter := protectedRouter.Group("/employee")
	employeeRouter.POST("/", h.CreateEmployeeHandler)
	employeeRouter.PATCH("/:id", h.PartialUpdateEmployeeHandler)
	employeeRouter.PUT("/:id", h.UpdateEmployeeHandler)
	employeeRouter.DELETE("/:id", h.DeleteEmployeeHandler)
	employeeRouter.GET("/:id", h.GetEmployeebyIDHandler)
	employeeRouter.GET("/list", h.ListEmployeeHandler)

	//---add the following line above in CoreHandlerRegistryOptions struct
	//---EmployeeService  coreService.EmployeeService
	//---delete these lines after copy

	attendanceLogRouter := protectedRouter.Group("/attendance-log")
	attendanceLogRouter.POST("/", h.CreateAttendanceLogHandler)
	attendanceLogRouter.PATCH("/:id", h.PartialUpdateAttendanceLogHandler)
	attendanceLogRouter.PUT("/:id", h.UpdateAttendanceLogHandler)
	attendanceLogRouter.GET("/:id", h.GetAttendanceLogbyIDHandler)
	attendanceLogRouter.GET("/list", h.ListAttendanceLogHandler)
	//attendanceLogRouter.GET("/daily", h.ListDailyAttendanceLogHandler)

	attendanceLogRouter.GET("/daily", h.GetDailyAttendanceHandler)

	biometricDeviceRouter := protectedRouter.Group("/device")
	biometricDeviceRouter.POST("/add-user", h.CreateUserOnDeviceHandler)
	//---add the following line above in CoreHandlerRegistryOptions struct
	//---AttendanceLogService  coreService.AttendanceLogService
	//---delete these lines after copy

	//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
	return router, nil
}
