package rest

import (
	"fmt"

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
	dbConn, err := database.NewDBConnector()
	if err != nil {
		return nil, err
	}

	authGroup := e.Group(rootPath + "/auth")
	{
		authGroup.POST("/signup", NewUserHandler(usecase.NewUserUsecase(repository.NewUserRepository(dbConn))).CreateUser())
	}

	streamGroup := e.Group(rootPath + "/stream")
	locationGroup := streamGroup.Group("/location")
	{
		locationGroup.GET("/live", NewLocationHandler(usecase.NewLocationUsecase(service.NewLocationService(repository.NewLocationRepository(dbConn)))).StreamLiveLocation())        // 現在受信中の位置情報データの取得(20分間受信が途絶えると終了)
		locationGroup.POST("/archive", NewLocationHandler(usecase.NewLocationUsecase(service.NewLocationService(repository.NewLocationRepository(dbConn)))).StreamArchiveLocation()) // 過去の位置情報データの取得
	}

	// demonstrationGroup := e.Group(rootPath + "/demonstration")
	{
		// demonstrationGroup.GET("/")
	}

	return e, nil
}
