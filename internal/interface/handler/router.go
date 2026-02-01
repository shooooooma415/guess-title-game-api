package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shooooooma415/guess-title-game-api/internal/interface/websocket"
)

// NewRouter creates a new HTTP router
func NewRouter(
	userHandler *UserHandler,
	roomHandler *RoomHandler,
	wsHandler *websocket.Handler,
) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://172.16.30.111:3000",
			"http://172.16.30.111:3001",
		},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

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
