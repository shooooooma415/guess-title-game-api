package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewRouter creates a new HTTP router
func NewRouter() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check
	e.GET("/health", healthCheck)

	// API routes
	api := e.Group("/api/v1")
	{
		// User routes
		// users := api.Group("/users")
		// {
		// 	users.POST("", userHandler.Create)
		// 	users.GET("/:id", userHandler.GetByID)
		// }

		// Room routes
		// rooms := api.Group("/rooms")
		// {
		// 	rooms.POST("", roomHandler.Create)
		// 	rooms.GET("/:id", roomHandler.GetByID)
		// 	rooms.POST("/:code/join", roomHandler.Join)
		// }
		_ = api
	}

	return e
}

func healthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}
