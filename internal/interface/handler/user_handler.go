package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shooooooma415/guess-title-game-api/internal/usecase/user"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	joinRoomUseCase *user.JoinRoomUseCase
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(joinRoomUseCase *user.JoinRoomUseCase) *UserHandler {
	return &UserHandler{
		joinRoomUseCase: joinRoomUseCase,
	}
}

// JoinRoomRequest represents the request body for joining a room
type JoinRoomRequest struct {
	RoomCode string `json:"room_code"`
	UserName string `json:"user_name"`
}

// JoinRoomResponse represents the response for joining a room
type JoinRoomResponse struct {
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	IsLeader bool   `json:"is_leader"`
}

// JoinRoom handles POST /api/user
func (h *UserHandler) JoinRoom(c echo.Context) error {
	var req JoinRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	input := user.JoinRoomInput{
		RoomCode: req.RoomCode,
		UserName: req.UserName,
	}

	output, err := h.joinRoomUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	response := JoinRoomResponse{
		RoomID:   output.RoomID,
		UserID:   output.UserID,
		IsLeader: output.IsLeader,
	}

	return c.JSON(http.StatusOK, response)
}
