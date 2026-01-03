package room_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/theme"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/user"
	roomUseCase "github.com/shooooooma415/guess-title-game-api/internal/usecase/room"
)

// Helper function to create a test theme
func createTestTheme() *theme.Theme {
	themeID := theme.NewThemeID()
	themeTitle, _ := theme.NewThemeTitle("Test Theme")
	hint := theme.NewHint("Test Hint")
	return theme.NewTheme(themeID, themeTitle, hint)
}
func TestCreateRoomUseCase_Execute_Success(t *testing.T) {
	// Setup
	testTheme := createTestTheme()

	userRepo := &mockUserRepository{}
	roomRepo := &mockRoomRepository{}
	themeRepo := &mockThemeRepository{
		findAllFunc: func(ctx context.Context) ([]*theme.Theme, error) {
			return []*theme.Theme{testTheme}, nil
		},
	}
	participantRepo := &mockParticipantRepository{}

	useCase := roomUseCase.NewCreateRoomUseCase(userRepo, roomRepo, themeRepo, participantRepo)

	// Execute
	output, err := useCase.Execute(context.Background())

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if output == nil {
		t.Fatal("Expected output, got nil")
	}

	if output.RoomID == "" {
		t.Error("Expected RoomID to be set")
	}

	if output.UserID == "" {
		t.Error("Expected UserID to be set")
	}

	if output.RoomCode == "" {
		t.Error("Expected RoomCode to be set")
	}

	if output.Theme != "Test Theme" {
		t.Errorf("Expected Theme to be 'Test Theme', got: %s", output.Theme)
	}

	if output.Hint != "Test Hint" {
		t.Errorf("Expected Hint to be 'Test Hint', got: %s", output.Hint)
	}
}

func TestCreateRoomUseCase_Execute_NoThemes(t *testing.T) {
	// Setup
	themeRepo := &mockThemeRepository{
		findAllFunc: func(ctx context.Context) ([]*theme.Theme, error) {
			return []*theme.Theme{}, nil
		},
	}

	useCase := roomUseCase.NewCreateRoomUseCase(
		&mockUserRepository{},
		&mockRoomRepository{},
		themeRepo,
		&mockParticipantRepository{},
	)

	// Execute
	output, err := useCase.Execute(context.Background())

	// Assert
	if err == nil {
		t.Error("Expected error when no themes available")
	}

	if output != nil {
		t.Error("Expected nil output when error occurs")
	}
}

func TestCreateRoomUseCase_Execute_ThemeRepoError(t *testing.T) {
	// Setup
	themeRepo := &mockThemeRepository{
		findAllFunc: func(ctx context.Context) ([]*theme.Theme, error) {
			return nil, errors.New("database error")
		},
	}

	useCase := roomUseCase.NewCreateRoomUseCase(
		&mockUserRepository{},
		&mockRoomRepository{},
		themeRepo,
		&mockParticipantRepository{},
	)

	// Execute
	output, err := useCase.Execute(context.Background())

	// Assert
	if err == nil {
		t.Error("Expected error when theme repository fails")
	}

	if output != nil {
		t.Error("Expected nil output when error occurs")
	}
}

func TestCreateRoomUseCase_Execute_UserSaveError(t *testing.T) {
	// Setup
	testTheme := createTestTheme()

	userRepo := &mockUserRepository{
		saveFunc: func(ctx context.Context, _ *user.User) error {
			return errors.New("user save error")
		},
	}
	roomRepo := &mockRoomRepository{}
	themeRepo := &mockThemeRepository{
		findAllFunc: func(ctx context.Context) ([]*theme.Theme, error) {
			return []*theme.Theme{testTheme}, nil
		},
	}
	participantRepo := &mockParticipantRepository{}

	useCase := roomUseCase.NewCreateRoomUseCase(userRepo, roomRepo, themeRepo, participantRepo)

	// Execute
	output, err := useCase.Execute(context.Background())

	// Assert
	if err == nil {
		t.Error("Expected error when user save fails")
	}

	if output != nil {
		t.Error("Expected nil output when error occurs")
	}
}

func TestCreateRoomUseCase_Execute_RoomSaveError(t *testing.T) {
	// Setup
	testTheme := createTestTheme()

	roomRepo := &mockRoomRepository{
		saveFunc: func(ctx context.Context, _ *room.Room) error {
			return errors.New("room save error")
		},
	}
	themeRepo := &mockThemeRepository{
		findAllFunc: func(ctx context.Context) ([]*theme.Theme, error) {
			return []*theme.Theme{testTheme}, nil
		},
	}

	useCase := roomUseCase.NewCreateRoomUseCase(
		&mockUserRepository{},
		roomRepo,
		themeRepo,
		&mockParticipantRepository{},
	)

	// Execute
	output, err := useCase.Execute(context.Background())

	// Assert
	if err == nil {
		t.Error("Expected error when room save fails")
	}

	if output != nil {
		t.Error("Expected nil output when error occurs")
	}
}

func TestCreateRoomUseCase_Execute_ParticipantSaveError(t *testing.T) {
	// Setup
	testTheme := createTestTheme()

	themeRepo := &mockThemeRepository{
		findAllFunc: func(ctx context.Context) ([]*theme.Theme, error) {
			return []*theme.Theme{testTheme}, nil
		},
	}
	participantRepo := &mockParticipantRepository{
		saveFunc: func(ctx context.Context, _ *participant.Participant) error {
			return errors.New("participant save error")
		},
	}

	useCase := roomUseCase.NewCreateRoomUseCase(
		&mockUserRepository{},
		&mockRoomRepository{},
		themeRepo,
		participantRepo,
	)

	// Execute
	output, err := useCase.Execute(context.Background())

	// Assert
	if err == nil {
		t.Error("Expected error when participant save fails")
	}

	if output != nil {
		t.Error("Expected nil output when error occurs")
	}
}
