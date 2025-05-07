package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/To-ge/gr_backend_go/config"
	"github.com/To-ge/gr_backend_go/domain/service"
	"github.com/To-ge/gr_backend_go/infrastructure/database"
	"github.com/To-ge/gr_backend_go/infrastructure/repository"
	"github.com/To-ge/gr_backend_go/usecase"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type restMiddleware struct {
	rdbc        *database.RedisConnector
	authUsecase usecase.IAuthenticationUsecase
}

func NewRestMiddleware(dbc *database.DBConnector, rdbc *database.RedisConnector) *restMiddleware {
	return &restMiddleware{
		rdbc:        rdbc,
		authUsecase: usecase.NewAuthenticationUsecase(service.NewAuthenticationService(repository.NewUserRepository(dbc), repository.NewAuthenticationRepository(dbc, rdbc))),
	}
}

// Session middleware to check session
func (rm *restMiddleware) SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := config.SessionStore.Get(c.Request(), "session")
		if err != nil {
			log.Printf("session-check is failed. %v\n", err.Error())
			return err
		}
		value := sess.Values[config.SessionKey]
		if value == nil {
			log.Println("session-check is failed.")
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		key := value.(string)
		if _, err := rm.rdbc.Conn.Get(context.Background(), key).Result(); err == redis.Nil {
			log.Printf("session-check is failed. %v\n", err.Error())
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "session expired or invalid"})
		} else if err != nil {
			log.Printf("session-check is failed. %v\n", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "session validation failed"})
		}
		rm.authUsecase.RefreshSessionExpiration(key)

		log.Println("session-check is succesful.")
		return next(c)
	}
}
