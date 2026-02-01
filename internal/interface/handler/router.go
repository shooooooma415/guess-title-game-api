package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shooooooma415/guess-title-game-api/config"
	customMiddleware "github.com/shooooooma415/guess-title-game-api/internal/interface/middleware"
	"github.com/shooooooma415/guess-title-game-api/internal/interface/websocket"
)

// NewRouter creates a new HTTP router
func NewRouter(
	cfg *config.Config,
	userHandler *UserHandler,
	roomHandler *RoomHandler,
	wsHandler *websocket.Handler,
) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(customMiddleware.CORSConfig(cfg))

	// Health check
	e.GET("/health", healthCheck)

	// WebSocket endpoint
	e.GET("/ws", wsHandler.HandleWebSocket)

	// API routes
	api := e.Group("/api")
	{
		// User routes
		api.POST("/user", userHandler.JoinRoom)

		// Room routes
		api.POST("/rooms", roomHandler.CreateRoom)
		api.POST("/rooms/:room_id/start", roomHandler.StartGame)
		api.POST("/rooms/:room_id/topic", roomHandler.SetTopic)
		api.POST("/rooms/:room_id/answer", roomHandler.SubmitAnswer)
		api.POST("/rooms/:room_id/skip-discussion", roomHandler.SkipDiscussion)
		api.POST("/rooms/:room_id/finish", roomHandler.FinishGame)
	}

	return e
}

func healthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "ok",
	})
}
