package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"subscription_service/internal/service"
)

func errorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		switch {
		case errors.Is(err, service.ErrSubscriptionNotFound):
			return c.NoContent(http.StatusNotFound)

		default:
			return c.NoContent(http.StatusInternalServerError)
		}
	}
}
