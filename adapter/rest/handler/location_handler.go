package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/To-ge/gr_backend_go/config"
	"github.com/To-ge/gr_backend_go/pkg"
	"github.com/To-ge/gr_backend_go/usecase"
	"github.com/To-ge/gr_backend_go/usecase/model"
	"github.com/labstack/echo/v4"
)

type locationHandler struct {
	usecase usecase.ILocationUsecase
}

func NewLocationHandler(lu usecase.ILocationUsecase) *locationHandler {
	return &locationHandler{
		usecase: lu,
	}
}

func (lh *locationHandler) StreamLiveLocation() echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("locationHandler.StreamLiveLocation started.")
		defer log.Println("locationHandler.StreamLiveLocation ended.")
		var input *model.StreamLiveLocationInput
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}
		output, err := lh.usecase.StreamLiveLocation(input)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// ヘッダーを設定
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set(echo.HeaderContentType, "application/json; charset=utf-8")
		c.Response().Header().Set("Transfer-Encoding", "chunked")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().WriteHeader(http.StatusOK)

		// ストリームを取得
		writer := c.Response()

		ctx := c.Request().Context()

		var logger *log.Logger
		if config.LoadConfig().Mode == config.Demo {
			fmt.Println("demo mode")
			clientID := fmt.Sprintf("client_%d", time.Now().UnixNano()) // クライアント識別用の一意なIDを生成
			logFilePath := fmt.Sprintf("../gps-out/%s.log", clientID)
			if newLogger, cleanUp, err := pkg.CreateLogFile(logFilePath); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "server error"})
			} else {
				logger = newLogger
				defer cleanUp()
			}
		}

		for {
			select {
			case <-ctx.Done(): // クライアント接続が切れたかどうかを監視する
				log.Println("locationHandler.StreamLiveLocation: Client disconnected or request canceled")
				return nil
			case location, ok := <-output.LocationChannel:
				if !ok {
					return nil
				}

				// JSONエンコード
				locationJSON, err := json.Marshal(location)
				if err != nil {
					log.Printf("Error encoding location: %v\n", err)
					continue
				}

				// 書き込み
				if _, err := writer.Write(append(locationJSON, '\n')); err != nil {
					log.Printf("Error writing to response: %v", err)
					continue
				}

				currentTime := float64(time.Now().UnixMicro()) / math.Pow10(6)
				pkg.OutputLocationLogger.Printf(",%v,%v,%v,%v\n", currentTime, location.Latitude, location.Longitude, location.Altitude)
				if logger != nil {
					logger.Printf(",%f,%v,%v,%v\n", currentTime, location.Latitude, location.Longitude, location.Altitude)
				}
				writer.Flush()
			}
		}
	}
}

func (lh *locationHandler) StreamArchiveLocation() echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("locationHandler.StreamArchiveLocation started.")
		defer log.Println("locationHandler.StreamArchiveLocation ended.")
		var input *model.StreamArchiveLocationInput
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}
		output, err := lh.usecase.StreamArchiveLocation(input)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		// ヘッダーを設定
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set(echo.HeaderContentType, "application/json; charset=utf-8")
		c.Response().Header().Set("Transfer-Encoding", "chunked")
		c.Response().Header().Set("Cache-Control", "no-cache")
		c.Response().WriteHeader(http.StatusOK)

		// ストリームを取得
		writer := c.Response()

		for {
			location, ok := <-output.LocationChannel

			if !ok {
				break
			}

			// JSONエンコード
			locationJSON, err := json.Marshal(location)
			if err != nil {
				log.Printf("Error encoding location: %v\n", err)
				continue
			}

			// 書き込み
			if _, err := writer.Write(append(locationJSON, '\n')); err != nil {
				log.Printf("Error writing to response: %v", err)
				continue
			}

			writer.Flush()
		}

		return nil
	}
}
