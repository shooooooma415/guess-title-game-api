package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
	domainUser "github.com/shooooooma415/guess-title-game-api/internal/domain/user"
	userUseCase "github.com/shooooooma415/guess-title-game-api/internal/usecase/user"
)

func TestJoinRoomUseCaseExecute(t *testing.T) {
	type fixture struct {
		useCase         *userUseCase.JoinRoomUseCase
		userRepo        *mockUserRepository
		roomRepo        *mockRoomRepository
		participantRepo *mockParticipantRepository
	}

	newFixture := func(t *testing.T) *fixture {
		t.Helper()

		userRepo := &mockUserRepository{}
		roomRepo := &mockRoomRepository{}
		participantRepo := &mockParticipantRepository{}

		useCase := userUseCase.NewJoinRoomUseCase(
			userRepo,
			roomRepo,
			participantRepo,
		)

		return &fixture{
			useCase:         useCase,
			userRepo:        userRepo,
			roomRepo:        roomRepo,
			participantRepo: participantRepo,
		}
	}

	createTestRoom := func() *room.Room {
		roomID := room.NewRoomID()
		roomCode := room.NewRoomCode()
		themeID, _ := room.NewThemeIDFromString("550e8400-e29b-41d4-a716-446655440000")
		hostUserID, _ := room.NewHostUserIDFromString("550e8400-e29b-41d4-a716-446655440001")
		return room.NewRoom(roomID, roomCode, themeID, hostUserID)
	}

	t.Run("正常にルームに参加できること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoom()

		f.roomRepo.findByCodeFunc = func(ctx context.Context, code room.RoomCode) (*room.Room, error) {
			return testRoom, nil
		}
		f.participantRepo.findByRoomIDFunc = func(ctx context.Context, roomID participant.RoomID) ([]*participant.Participant, error) {
			return []*participant.Participant{}, nil
		}

		input := userUseCase.JoinRoomInput{
			RoomCode: testRoom.Code().String(),
			UserName: "Test User",
		}

		// act
		output, err := f.useCase.Execute(context.Background(), input)

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

		if !output.IsLeader {
			t.Error("Expected first participant to be leader")
		}
	})

	t.Run("RoomCodeが空の場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)

		input := userUseCase.JoinRoomInput{
			RoomCode: "",
			UserName: "Test User",
		}

		// act
		output, err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when room code is empty")
		}
		if err.Error() != "room code is required" {
			t.Errorf("Expected 'room code is required' error, got: %v", err)
		}
		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})

	t.Run("UserNameが空の場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)

		input := userUseCase.JoinRoomInput{
			RoomCode: "123456",
			UserName: "",
		}

		// act
		output, err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when user name is empty")
		}
		if err.Error() != "user name is required" {
			t.Errorf("Expected 'user name is required' error, got: %v", err)
		}
		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})

	t.Run("無効なRoomCodeの場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)

		input := userUseCase.JoinRoomInput{
			RoomCode: "12345", // Invalid: must be 6 characters
			UserName: "Test User",
		}

		// act
		output, err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error for invalid room code")
		}
		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})

	t.Run("Roomが見つからない場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)

		f.roomRepo.findByCodeFunc = func(ctx context.Context, code room.RoomCode) (*room.Room, error) {
			return nil, errors.New("not found")
		}

		input := userUseCase.JoinRoomInput{
			RoomCode: "123456",
			UserName: "Test User",
		}

		// act
		output, err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when room not found")
		}
		if err.Error() != "room not found" {
			t.Errorf("Expected 'room not found' error, got: %v", err)
		}
		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})

	t.Run("Userの保存に失敗した場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoom()

		f.roomRepo.findByCodeFunc = func(ctx context.Context, code room.RoomCode) (*room.Room, error) {
			return testRoom, nil
		}
		f.userRepo.saveFunc = func(ctx context.Context, u *domainUser.User) error {
			return errors.New("save error")
		}

		input := userUseCase.JoinRoomInput{
			RoomCode: testRoom.Code().String(),
			UserName: "Test User",
		}

		// act
		output, err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when user save fails")
		}
		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})

	t.Run("Participantの保存に失敗した場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoom()

		f.roomRepo.findByCodeFunc = func(ctx context.Context, code room.RoomCode) (*room.Room, error) {
			return testRoom, nil
		}
		f.participantRepo.findByRoomIDFunc = func(ctx context.Context, roomID participant.RoomID) ([]*participant.Participant, error) {
			return []*participant.Participant{}, nil
		}
		f.participantRepo.saveFunc = func(ctx context.Context, p *participant.Participant) error {
			return errors.New("save error")
		}

		input := userUseCase.JoinRoomInput{
			RoomCode: testRoom.Code().String(),
			UserName: "Test User",
		}

		// act
		output, err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when participant save fails")
		}
		if output != nil {
			t.Error("Expected nil output when error occurs")
		}
	})

	t.Run("2人目以降の参加者はリーダーにならないこと", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoom()

		// Existing participant
		existingParticipantID := participant.NewParticipantID()
		existingRoomID, _ := participant.NewRoomIDFromString(testRoom.ID().String())
		existingUserID, _ := participant.NewUserIDFromString("550e8400-e29b-41d4-a716-446655440002")
		existingParticipant := participant.NewParticipant(existingParticipantID, existingRoomID, existingUserID, participant.RolePlayer)
		existingParticipant.SetAsLeader()

		f.roomRepo.findByCodeFunc = func(ctx context.Context, code room.RoomCode) (*room.Room, error) {
			return testRoom, nil
		}
		f.participantRepo.findByRoomIDFunc = func(ctx context.Context, roomID participant.RoomID) ([]*participant.Participant, error) {
			return []*participant.Participant{existingParticipant}, nil
		}

		input := userUseCase.JoinRoomInput{
			RoomCode: testRoom.Code().String(),
			UserName: "Test User 2",
		}

		// act
		output, err := f.useCase.Execute(context.Background(), input)

		// assert
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if output == nil {
			t.Fatal("Expected output, got nil")
		}

		if output.IsLeader {
			t.Error("Expected second participant not to be leader")
		}
	})
}
