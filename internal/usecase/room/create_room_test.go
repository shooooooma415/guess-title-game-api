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

func TestCreateRoomUseCaseExecute(t *testing.T) {
	type fixture struct {
		useCase         *roomUseCase.CreateRoomUseCase
		userRepo        *mockUserRepository
		roomRepo        *mockRoomRepository
		themeRepo       *mockThemeRepository
		participantRepo *mockParticipantRepository
	}

	newFixture := func(t *testing.T) *fixture {
		t.Helper()

		userRepo := &mockUserRepository{}
		roomRepo := &mockRoomRepository{}
		themeRepo := &mockThemeRepository{}
		participantRepo := &mockParticipantRepository{}

		useCase := roomUseCase.NewCreateRoomUseCase(
			userRepo,
			roomRepo,
			themeRepo,
			participantRepo,
		)

		return &fixture{
			useCase:         useCase,
			userRepo:        userRepo,
			roomRepo:        roomRepo,
			themeRepo:       themeRepo,
			participantRepo: participantRepo,
		}
	}

	createTestTheme := func() *theme.Theme {
		themeID := theme.NewThemeID()
		themeTitle, _ := theme.NewThemeTitle("Test Theme")
		hint := theme.NewHint("Test Hint")
		return theme.NewTheme(themeID, themeTitle, hint)
	}

	t.Run("正常にルームが作成されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testTheme := createTestTheme()
		f.themeRepo.findAllFunc = func(ctx context.Context) ([]*theme.Theme, error) {
			return []*theme.Theme{testTheme}, nil
		}

		// act
		output, err := f.useCase.Execute(context.Background())

		// assert
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
	})

	t.Run("テーマが存在しない場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		f.themeRepo.findAllFunc = func(ctx context.Context) ([]*theme.Theme, error) {
			return []*theme.Theme{}, nil
		}

		// act
		output, err := f.useCase.Execute(context.Background())

		// assert
		if err == nil {
			t.Error("Expected error when no themes available")
		}

		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})

	t.Run("ThemeRepositoryでエラーが発生した場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		f.themeRepo.findAllFunc = func(ctx context.Context) ([]*theme.Theme, error) {
			return nil, errors.New("database error")
		}

		// act
		output, err := f.useCase.Execute(context.Background())

		// assert
		if err == nil {
			t.Error("Expected error when theme repository fails")
		}

		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})

	t.Run("Userの保存に失敗した場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testTheme := createTestTheme()
		f.themeRepo.findAllFunc = func(ctx context.Context) ([]*theme.Theme, error) {
			return []*theme.Theme{testTheme}, nil
		}
		f.userRepo.saveFunc = func(ctx context.Context, _ *user.User) error {
			return errors.New("user save error")
		}

		// act
		output, err := f.useCase.Execute(context.Background())

		// assert
		if err == nil {
			t.Error("Expected error when user save fails")
		}

		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})

	t.Run("Roomの保存に失敗した場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testTheme := createTestTheme()
		f.themeRepo.findAllFunc = func(ctx context.Context) ([]*theme.Theme, error) {
			return []*theme.Theme{testTheme}, nil
		}
		f.roomRepo.saveFunc = func(ctx context.Context, _ *room.Room) error {
			return errors.New("room save error")
		}

		// act
		output, err := f.useCase.Execute(context.Background())

		// assert
		if err == nil {
			t.Error("Expected error when room save fails")
		}

		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})

	t.Run("Participantの保存に失敗した場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testTheme := createTestTheme()
		f.themeRepo.findAllFunc = func(ctx context.Context) ([]*theme.Theme, error) {
			return []*theme.Theme{testTheme}, nil
		}
		f.participantRepo.saveFunc = func(ctx context.Context, _ *participant.Participant) error {
			return errors.New("participant save error")
		}

		// act
		output, err := f.useCase.Execute(context.Background())

		// assert
		if err == nil {
			t.Error("Expected error when participant save fails")
		}

		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})
}
