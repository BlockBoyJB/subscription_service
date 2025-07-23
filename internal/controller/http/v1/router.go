package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "subscription_service/docs"
	"subscription_service/internal/service"
)

func NewRouter(g *echo.Echo, services *service.Services) {
	g.Use(middleware.Recover())
	g.Use(errorMiddleware)

	g.GET("/ping", ping)
	g.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := g.Group("/api/v1")

	newSubscriptionRouter(v1.Group("/subscription"), services.Subscription)
}

func ping(c echo.Context) error {
	return c.NoContent(200)
}
