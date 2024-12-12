package handler

import (
	"net/http"

	"github.com/To-ge/gr_backend_go/usecase"
	"github.com/To-ge/gr_backend_go/usecase/model"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	usecase usecase.IUserUsecase
}

func NewUserHandler(uu usecase.IUserUsecase) *userHandler {
	return &userHandler{
		usecase: uu,
	}
}

func (uh *userHandler) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input *model.CreateUserInput
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}
		user, err := uh.usecase.CreateUser(input)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, user)
	}
}
