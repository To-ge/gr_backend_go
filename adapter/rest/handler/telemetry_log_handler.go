package handler

import (
	"net/http"

	"github.com/To-ge/gr_backend_go/usecase"
	"github.com/To-ge/gr_backend_go/usecase/model"
	"github.com/labstack/echo/v4"
)

type telemetryLogHandler struct {
	usecase usecase.ITelemetryLogUsecase
}

func NewTelemetryLogHandler(tlu usecase.ITelemetryLogUsecase) *telemetryLogHandler {
	return &telemetryLogHandler{
		usecase: tlu,
	}
}

func (uh *telemetryLogHandler) GetTelemetryLogs() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input *model.GetTelemetryLogsInput
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}
		logs, err := uh.usecase.GetTelemetryLogs(input)

		// c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set(echo.HeaderContentType, "application/json; charset=utf-8")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, logs)
	}
}
