package handler

import (
	"fmt"
	"net/http"

	"github.com/To-ge/gr_backend_go/config"
	"github.com/To-ge/gr_backend_go/pkg"
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

func (tlh *telemetryLogHandler) GetTelemetryLogs() echo.HandlerFunc {
	return func(c echo.Context) error {
		isAdmin, ok := c.Get(config.ContextKeyIsAdmin).(bool)
		if !ok {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "authorization failed"})
		}
		var input *model.GetTelemetryLogsInput
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}

		logs, err := tlh.usecase.GetTelemetryLogs(input, isAdmin)

		// c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set(echo.HeaderContentType, "application/json; charset=utf-8")

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, logs)
	}
}

func (tlh *telemetryLogHandler) ToggleTelemetryLogVisibility() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input *model.ToggleTelemetryLogVisibilityInput
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}

		err := tlh.usecase.ToggleTelemetryLogVisibility(input)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		message := fmt.Sprintf("Changing ID:%d %s successful.", input.Id, pkg.Ternary(input.Visible, "public", "private"))
		return c.JSON(http.StatusOK, map[string]string{"message": message})
	}
}
