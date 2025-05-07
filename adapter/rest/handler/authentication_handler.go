package handler

import (
	"net/http"

	"github.com/To-ge/gr_backend_go/usecase"
	"github.com/To-ge/gr_backend_go/usecase/model"
	"github.com/labstack/echo/v4"
)

type authenticationHandler struct {
	usecase usecase.IAuthenticationUsecase
}

func NewAuthenticationHandler(au usecase.IAuthenticationUsecase) *authenticationHandler {
	return &authenticationHandler{
		usecase: au,
	}
}

func (ah *authenticationHandler) SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		var input *model.SignInInput
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
		}
		input.Request = c.Request()
		input.ResponseWriter = c.Response()

		err := ah.usecase.SignIn(input)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "signed in successfully"})
	}
}

func (ah *authenticationHandler) SignOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &model.SignOutInput{
			Request:        c.Request(),
			ResponseWriter: c.Response(),
		}

		err := ah.usecase.SignOut(input)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "signed out successfully"})
	}
}

func (ah *authenticationHandler) SessionCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "session is valid"})
	}
}
