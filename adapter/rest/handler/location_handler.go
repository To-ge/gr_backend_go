package handler

import (
	"encoding/json"
	"log"
	"net/http"

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

func (lh *locationHandler) StreamArchiveLocation() echo.HandlerFunc {
	return func(c echo.Context) error {
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
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextPlainCharsetUTF8)
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
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to encode location"})
			}

			// 書き込み
			if _, err := writer.Write(append(locationJSON, '\n')); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to write to response"})
			}

			writer.Flush()
		}

		return nil
	}
}
