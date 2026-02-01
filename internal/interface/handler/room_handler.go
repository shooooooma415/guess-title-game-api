package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	roomUseCase "github.com/shooooooma415/guess-title-game-api/internal/usecase/room"
)

// RoomHandler handles room-related HTTP requests
type RoomHandler struct {
	createRoomUseCase     *roomUseCase.CreateRoomUseCase
	startGameUseCase      *roomUseCase.StartGameUseCase
	setTopicUseCase       *roomUseCase.SetTopicUseCase
	submitAnswerUseCase   *roomUseCase.SubmitAnswerUseCase
	skipDiscussionUseCase *roomUseCase.SkipDiscussionUseCase
	finishGameUseCase     *roomUseCase.FinishGameUseCase
}

// NewRoomHandler creates a new RoomHandler
func NewRoomHandler(
	createRoomUseCase *roomUseCase.CreateRoomUseCase,
	startGameUseCase *roomUseCase.StartGameUseCase,
	setTopicUseCase *roomUseCase.SetTopicUseCase,
	submitAnswerUseCase *roomUseCase.SubmitAnswerUseCase,
	skipDiscussionUseCase *roomUseCase.SkipDiscussionUseCase,
	finishGameUseCase *roomUseCase.FinishGameUseCase,
) *RoomHandler {
	return &RoomHandler{
		createRoomUseCase:     createRoomUseCase,
		startGameUseCase:      startGameUseCase,
		setTopicUseCase:       setTopicUseCase,
		submitAnswerUseCase:   submitAnswerUseCase,
		skipDiscussionUseCase: skipDiscussionUseCase,
		finishGameUseCase:     finishGameUseCase,
	}
}

// CreateRoomResponse represents the response for creating a room
type CreateRoomResponse struct {
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	RoomCode string `json:"room_code"`
	Theme    string `json:"theme"`
	Hint     string `json:"hint"`
}

// CreateRoom handles POST /api/rooms
func (h *RoomHandler) CreateRoom(c echo.Context) error {
	output, err := h.createRoomUseCase.Execute(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	response := CreateRoomResponse{
		RoomID:   output.RoomID,
		UserID:   output.UserID,
		RoomCode: output.RoomCode,
		Theme:    output.Theme,
		Hint:     output.Hint,
	}

	return c.JSON(http.StatusOK, response)
}

// StartGameRequest represents the request body for starting a game
type StartGameRequest struct {
	UserID string `json:"user_id"`
}

// StartGame handles POST /api/rooms/:room_id/start
func (h *RoomHandler) StartGame(c echo.Context) error {
	roomID := c.Param("room_id")

	var req StartGameRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	input := roomUseCase.StartGameInput{
		RoomID: roomID,
		UserID: req.UserID,
	}

	if err := h.startGameUseCase.Execute(c.Request().Context(), input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "started",
	})
}

// SetTopicRequest represents the request body for setting a topic
type SetTopicRequest struct {
	UserID          string   `json:"user_id"`
	Topic           string   `json:"topic"`
	Emojis          []string `json:"emojis"`
	DisplayedEmojis []string `json:"displayed_emojis"`
	OriginalEmojis  []string `json:"original_emojis"`
	DummyIndex      int      `json:"dummy_index"`
	DummyEmoji      string   `json:"dummy_emoji"`
}

// SetTopic handles POST /api/rooms/:room_id/topic
func (h *RoomHandler) SetTopic(c echo.Context) error {
	roomID := c.Param("room_id")

	var req SetTopicRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	input := roomUseCase.SetTopicInput{
		RoomID:          roomID,
		UserID:          req.UserID,
		Topic:           req.Topic,
		Emojis:          req.Emojis,
		DisplayedEmojis: req.DisplayedEmojis,
		OriginalEmojis:  req.OriginalEmojis,
		DummyIndex:      req.DummyIndex,
		DummyEmoji:      req.DummyEmoji,
	}

	if err := h.setTopicUseCase.Execute(c.Request().Context(), input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "topic_set",
	})
}

// SubmitAnswerRequest represents the request body for submitting an answer
type SubmitAnswerRequest struct {
	UserID string `json:"user_id"`
	Answer string `json:"answer"`
}

// SubmitAnswer handles POST /api/rooms/:room_id/answer
func (h *RoomHandler) SubmitAnswer(c echo.Context) error {
	roomID := c.Param("room_id")

	var req SubmitAnswerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	input := roomUseCase.SubmitAnswerInput{
		RoomID: roomID,
		UserID: req.UserID,
		Answer: req.Answer,
	}

	if err := h.submitAnswerUseCase.Execute(c.Request().Context(), input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "answer_submitted",
	})
}

// SkipDiscussionRequest represents the request body for skipping discussion
type SkipDiscussionRequest struct {
	UserID string `json:"user_id"`
}

// SkipDiscussion handles POST /api/rooms/:room_id/skip-discussion
func (h *RoomHandler) SkipDiscussion(c echo.Context) error {
	roomID := c.Param("room_id")

	var req SkipDiscussionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	input := roomUseCase.SkipDiscussionInput{
		RoomID: roomID,
		UserID: req.UserID,
	}

	if err := h.skipDiscussionUseCase.Execute(c.Request().Context(), input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "discussion_skipped",
	})
}

// FinishGameRequest represents the request body for finishing a game
type FinishGameRequest struct {
	UserID string `json:"user_id"`
}

// FinishGame handles POST /api/rooms/:room_id/finish
func (h *RoomHandler) FinishGame(c echo.Context) error {
	roomID := c.Param("room_id")

	var req FinishGameRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	input := roomUseCase.FinishGameInput{
		RoomID: roomID,
		UserID: req.UserID,
	}

	if err := h.finishGameUseCase.Execute(c.Request().Context(), input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "game_finished",
	})
}
