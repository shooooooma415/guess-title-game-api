package room_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
	roomUseCase "github.com/shooooooma415/guess-title-game-api/internal/usecase/room"
)

func TestStartGameUseCaseExecute(t *testing.T) {
	type fixture struct {
		useCase         *roomUseCase.StartGameUseCase
		roomRepo        *mockRoomRepository
		participantRepo *mockParticipantRepository
		eventPublisher  *mockEventPublisher
	}

	newFixture := func(t *testing.T) *fixture {
		t.Helper()

		roomRepo := &mockRoomRepository{}
		participantRepo := &mockParticipantRepository{}
		eventPublisher := &mockEventPublisher{}

		useCase := roomUseCase.NewStartGameUseCase(
			roomRepo,
			participantRepo,
			eventPublisher,
		)

		return &fixture{
			useCase:         useCase,
			roomRepo:        roomRepo,
			participantRepo: participantRepo,
			eventPublisher:  eventPublisher,
		}
	}

	createTestRoom := func() *room.Room {
		roomID := room.NewRoomID()
		roomCode := room.NewRoomCode()
		themeID, _ := room.NewThemeIDFromString("550e8400-e29b-41d4-a716-446655440000")
		hostUserID, _ := room.NewHostUserIDFromString("550e8400-e29b-41d4-a716-446655440001")
		return room.NewRoom(roomID, roomCode, themeID, hostUserID)
	}

	createHostParticipant := func(roomID, userID string) *participant.Participant {
		participantID := participant.NewParticipantID()
		participantRoomID, _ := participant.NewRoomIDFromString(roomID)
		participantUserID, _ := participant.NewUserIDFromString(userID)
		p := participant.NewParticipant(participantID, participantRoomID, participantUserID, participant.RoleHost)
		return p
	}

	t.Run("正常にゲームが開始されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoom()
		testParticipant := createHostParticipant(testRoom.ID().String(), "550e8400-e29b-41d4-a716-446655440001")

		f.roomRepo.findByIDFunc = func(ctx context.Context, id room.RoomID) (*room.Room, error) {
			return testRoom, nil
		}
		f.participantRepo.findByRoomAndUserFunc = func(ctx context.Context, roomID participant.RoomID, userID participant.UserID) (*participant.Participant, error) {
			return testParticipant, nil
		}

		input := roomUseCase.StartGameInput{
			RoomID: testRoom.ID().String(),
			UserID: "550e8400-e29b-41d4-a716-446655440001",
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
	})

	t.Run("無効なRoomIDの場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)

		input := roomUseCase.StartGameInput{
			RoomID: "invalid-uuid",
			UserID: "550e8400-e29b-41d4-a716-446655440001",
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error for invalid room ID")
		}
	})

	t.Run("Roomが見つからない場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)

		f.roomRepo.findByIDFunc = func(ctx context.Context, id room.RoomID) (*room.Room, error) {
			return nil, errors.New("not found")
		}

		input := roomUseCase.StartGameInput{
			RoomID: "550e8400-e29b-41d4-a716-446655440000",
			UserID: "550e8400-e29b-41d4-a716-446655440001",
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when room not found")
		}
		if err.Error() != "room not found" {
			t.Errorf("Expected 'room not found' error, got: %v", err)
		}
	})

	t.Run("Participantが見つからない場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoom()

		f.roomRepo.findByIDFunc = func(ctx context.Context, id room.RoomID) (*room.Room, error) {
			return testRoom, nil
		}
		f.participantRepo.findByRoomAndUserFunc = func(ctx context.Context, roomID participant.RoomID, userID participant.UserID) (*participant.Participant, error) {
			return nil, errors.New("not found")
		}

		input := roomUseCase.StartGameInput{
			RoomID: testRoom.ID().String(),
			UserID: "550e8400-e29b-41d4-a716-446655440001",
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when participant not found")
		}
		if err.Error() != "participant not found" {
			t.Errorf("Expected 'participant not found' error, got: %v", err)
		}
	})

	t.Run("ホスト以外のユーザーがゲーム開始しようとした場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoom()

		participantID := participant.NewParticipantID()
		participantRoomID, _ := participant.NewRoomIDFromString(testRoom.ID().String())
		participantUserID, _ := participant.NewUserIDFromString("550e8400-e29b-41d4-a716-446655440001")
		testParticipant := participant.NewParticipant(participantID, participantRoomID, participantUserID, participant.RolePlayer)

		f.roomRepo.findByIDFunc = func(ctx context.Context, id room.RoomID) (*room.Room, error) {
			return testRoom, nil
		}
		f.participantRepo.findByRoomAndUserFunc = func(ctx context.Context, roomID participant.RoomID, userID participant.UserID) (*participant.Participant, error) {
			return testParticipant, nil
		}

		input := roomUseCase.StartGameInput{
			RoomID: testRoom.ID().String(),
			UserID: "550e8400-e29b-41d4-a716-446655440001",
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when non-host tries to start game")
		}
		if err.Error() != "only host can start the game" {
			t.Errorf("Expected 'only host can start the game' error, got: %v", err)
		}
	})

	t.Run("Roomの保存に失敗した場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoom()
		testParticipant := createHostParticipant(testRoom.ID().String(), "550e8400-e29b-41d4-a716-446655440001")

		f.roomRepo.findByIDFunc = func(ctx context.Context, id room.RoomID) (*room.Room, error) {
			return testRoom, nil
		}
		f.participantRepo.findByRoomAndUserFunc = func(ctx context.Context, roomID participant.RoomID, userID participant.UserID) (*participant.Participant, error) {
			return testParticipant, nil
		}
		f.roomRepo.saveFunc = func(ctx context.Context, r *room.Room) error {
			return errors.New("save error")
		}

		input := roomUseCase.StartGameInput{
			RoomID: testRoom.ID().String(),
			UserID: "550e8400-e29b-41d4-a716-446655440001",
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when room save fails")
		}
	})
}
