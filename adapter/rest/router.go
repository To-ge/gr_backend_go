package rest

import (
	"fmt"
	"net/http"

	"github.com/To-ge/gr_backend_go/adapter/rest/handler"
	"github.com/To-ge/gr_backend_go/domain/service"
	"github.com/To-ge/gr_backend_go/infrastructure/database"
	"github.com/To-ge/gr_backend_go/infrastructure/repository"
	"github.com/To-ge/gr_backend_go/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	apiVersion = 1
	rootPath   = fmt.Sprintf("/api/v%d", apiVersion)
)

func InitRouter() (*echo.Echo, error) {
	e := echo.New()

	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)

	e.OPTIONS("/*", func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		return c.NoContent(http.StatusNoContent)
	})

	dbConn, err := database.NewDBConnector()
	if err != nil {
		return nil, err
	}

	authGroup := e.Group(rootPath + "/auth")
	{
		authGroup.POST("/signup", handler.NewUserHandler(usecase.NewUserUsecase(repository.NewUserRepository(dbConn))).CreateUser())
	}

	telemetryLogGroup := e.Group(rootPath + "/telemetry_log")
	{
		telemetryLogGroup.GET("", handler.NewTelemetryLogHandler(usecase.NewTelemetryLogUsecase(repository.NewTelemetryLogRepository(dbConn))).GetTelemetryLogs())
	}

	streamGroup := e.Group(rootPath + "/stream")
	locationGroup := streamGroup.Group("/location")
	{
		locationGroup.GET("/live", handler.NewLocationHandler(usecase.NewLocationUsecase(service.NewLocationService(repository.NewLocationRepository(dbConn)))).StreamLiveLocation())        // 現在受信中の位置情報データの取得
		locationGroup.POST("/archive", handler.NewLocationHandler(usecase.NewLocationUsecase(service.NewLocationService(repository.NewLocationRepository(dbConn)))).StreamArchiveLocation()) // 過去の位置情報データの取得
	}

	return e, nil
}
