package rest

import (
	"fmt"
	"net/http"

	"github.com/To-ge/gr_backend_go/adapter/middleware"
	"github.com/To-ge/gr_backend_go/adapter/rest/handler"
	"github.com/To-ge/gr_backend_go/config"
	"github.com/To-ge/gr_backend_go/domain/service"
	"github.com/To-ge/gr_backend_go/infrastructure/database"
	"github.com/To-ge/gr_backend_go/infrastructure/repository"
	"github.com/To-ge/gr_backend_go/usecase"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

var (
	apiVersion = 1
	rootPath   = fmt.Sprintf("/api/v%d", apiVersion)
)

func InitRouter() (*echo.Echo, error) {
	e := echo.New()

	var allowedOrigins = []string{
		config.LoadConfig().FEUrl,
	}
	e.Use(
		// echoMiddleware.Logger(),
		echoMiddleware.Recover(),
		echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowOrigins:     allowedOrigins,
			AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
		}),
	)

	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	dbConn, err := database.NewDBConnector()
	if err != nil {
		return nil, err
	}
	rdbConn := database.NewRedisConnector()

	restMiddleware := middleware.NewRestMiddleware(dbConn, rdbConn)

	authGroup := e.Group(rootPath + "/auth")
	{
		authGroup.POST("/signup", handler.NewUserHandler(usecase.NewUserUsecase(repository.NewUserRepository(dbConn))).CreateUser())
		authGroup.POST("/login", handler.NewAuthenticationHandler(usecase.NewAuthenticationUsecase(service.NewAuthenticationService(repository.NewUserRepository(dbConn), repository.NewAuthenticationRepository(dbConn, rdbConn)))).SignIn())
		authGroup.GET("/logout", handler.NewAuthenticationHandler(usecase.NewAuthenticationUsecase(service.NewAuthenticationService(repository.NewUserRepository(dbConn), repository.NewAuthenticationRepository(dbConn, rdbConn)))).SignOut())
		authGroup.GET("/session-check", handler.NewAuthenticationHandler(usecase.NewAuthenticationUsecase(service.NewAuthenticationService(repository.NewUserRepository(dbConn), repository.NewAuthenticationRepository(dbConn, rdbConn)))).SessionCheck(), restMiddleware.SessionMiddleware)
	}

	telemetryLogGroup := e.Group(rootPath+"/telemetry_log", restMiddleware.SessionMiddleware)
	{
		telemetryLogGroup.GET("", handler.NewTelemetryLogHandler(usecase.NewTelemetryLogUsecase(repository.NewTelemetryLogRepository(dbConn))).GetTelemetryLogs())
	}

	streamGroup := e.Group(rootPath + "/stream")
	locationGroup := streamGroup.Group("/location", restMiddleware.SessionMiddleware)
	{
		locationGroup.GET("/live", handler.NewLocationHandler(usecase.NewLocationUsecase(service.NewLocationService(repository.NewLocationRepository(dbConn)))).StreamLiveLocation())        // 現在受信中の位置情報データの取得
		locationGroup.POST("/archive", handler.NewLocationHandler(usecase.NewLocationUsecase(service.NewLocationService(repository.NewLocationRepository(dbConn)))).StreamArchiveLocation()) // 過去の位置情報データの取得
	}

	return e, nil
}
