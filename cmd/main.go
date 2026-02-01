package main

import (
	"log"

	"github.com/shooooooma415/guess-title-game-api/config"
	infrastructureEvent "github.com/shooooooma415/guess-title-game-api/internal/infrastructure/event"
	"github.com/shooooooma415/guess-title-game-api/internal/infrastructure/persistence"
	"github.com/shooooooma415/guess-title-game-api/internal/interface/handler"
	"github.com/shooooooma415/guess-title-game-api/internal/interface/websocket"
	roomUseCase "github.com/shooooooma415/guess-title-game-api/internal/usecase/room"
	userUseCase "github.com/shooooooma415/guess-title-game-api/internal/usecase/user"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	dbCfg := persistence.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}
	db, err := persistence.NewDB(dbCfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	// Initialize event publisher
	eventPublisher := infrastructureEvent.NewInMemoryEventPublisher()

	// Initialize repositories
	userRepo := persistence.NewUserRepository(db)
	roomRepo := persistence.NewRoomRepository(db)
	themeRepo := persistence.NewThemeRepository(db)
	participantRepo := persistence.NewParticipantRepository(db)

	// Initialize use cases
	joinRoomUseCase := userUseCase.NewJoinRoomUseCase(userRepo, roomRepo, participantRepo)
	createRoomUseCase := roomUseCase.NewCreateRoomUseCase(userRepo, roomRepo, themeRepo, participantRepo)
	startGameUseCase := roomUseCase.NewStartGameUseCase(roomRepo, participantRepo, eventPublisher)
	setTopicUseCase := roomUseCase.NewSetTopicUseCase(roomRepo, participantRepo)
	submitAnswerUseCase := roomUseCase.NewSubmitAnswerUseCase(roomRepo, participantRepo, eventPublisher)
	skipDiscussionUseCase := roomUseCase.NewSkipDiscussionUseCase(roomRepo, participantRepo, eventPublisher)
	finishGameUseCase := roomUseCase.NewFinishGameUseCase(roomRepo, participantRepo, eventPublisher)

	// Initialize WebSocket-specific use cases
	fetchRoomUseCase := roomUseCase.NewFetchRoomUseCase(roomRepo)
	fetchParticipantsUseCase := roomUseCase.NewFetchRoomParticipantsUseCase(participantRepo, userRepo)
	startDiscussionUseCase := roomUseCase.NewStartDiscussionUseCase(roomRepo, participantRepo)
	submitFinalAnswerUseCase := roomUseCase.NewSubmitFinalAnswerUseCase(roomRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(joinRoomUseCase)
	roomHandler := handler.NewRoomHandler(
		createRoomUseCase,
		startGameUseCase,
		setTopicUseCase,
		submitAnswerUseCase,
		skipDiscussionUseCase,
		finishGameUseCase,
	)

	// Initialize WebSocket hub and timer
	hub := websocket.NewHub()
	timer := websocket.NewTimer(hub)
	wsHandler := websocket.NewHandler(
		hub,
		timer,
		fetchRoomUseCase,
		fetchParticipantsUseCase,
		startDiscussionUseCase,
		submitFinalAnswerUseCase,
		themeRepo,
	)

	// Start WebSocket hub
	go hub.Run()

	// Setup event handlers for WebSocket
	wsHandler.SetupEventHandlers(eventPublisher)

	// Initialize router
	e := handler.NewRouter(cfg, userHandler, roomHandler, wsHandler)

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := e.Start(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
