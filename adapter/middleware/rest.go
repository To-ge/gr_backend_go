package middleware

import (
	"context"
	"net/http"

	"github.com/To-ge/gr_backend_go/config"
	"github.com/To-ge/gr_backend_go/infrastructure/database"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

// Session middleware to check session
func SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() == "/login" || c.Path() == "/logout" {
			return next(c)
		}
		sess, err := config.SessionStore.Get(c.Request(), "session")
		if err != nil {
			return err
		}
		value := sess.Values[config.SessionKey]
		if value == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}

		key := "session:" + value.(string)
		rdb := database.NewRedisConnector()
		if _, err := rdb.Conn.Get(context.Background(), key).Result(); err == redis.Nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "session expired or invalid"})
		} else if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "session validation failed"})
		}

		return next(c)
	}
}
