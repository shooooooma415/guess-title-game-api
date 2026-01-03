package room_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
	roomUseCase "github.com/shooooooma415/guess-title-game-api/internal/usecase/room"
)

func TestSetTopicUseCaseExecute(t *testing.T) {
	type fixture struct {
		useCase         *roomUseCase.SetTopicUseCase
		roomRepo        *mockRoomRepository
		participantRepo *mockParticipantRepository
	}

	newFixture := func(t *testing.T) *fixture {
		t.Helper()

		roomRepo := &mockRoomRepository{}
		participantRepo := &mockParticipantRepository{}

		useCase := roomUseCase.NewSetTopicUseCase(
			roomRepo,
			participantRepo,
		)

		return &fixture{
			useCase:         useCase,
			roomRepo:        roomRepo,
			participantRepo: participantRepo,
		}
	}

	createTestRoom := func() *room.Room {
		roomID := room.NewRoomID()
		roomCode := room.NewRoomCode()
		themeID, _ := room.NewThemeIDFromString("550e8400-e29b-41d4-a716-446655440000")
		hostUserID, _ := room.NewHostUserIDFromString("550e8400-e29b-41d4-a716-446655440001")
		r := room.NewRoom(roomID, roomCode, themeID, hostUserID)
		r.Start() // Set status to setting_topic
		return r
	}

	createHostParticipant := func(roomID, userID string) *participant.Participant {
		participantID := participant.NewParticipantID()
		participantRoomID, _ := participant.NewRoomIDFromString(roomID)
		participantUserID, _ := participant.NewUserIDFromString(userID)
		return participant.NewParticipant(participantID, participantRoomID, participantUserID, participant.RoleHost)
	}

	t.Run("正常にトピックが設定されること", func(t *testing.T) {
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

		input := roomUseCase.SetTopicInput{
			RoomID: testRoom.ID().String(),
			UserID: "550e8400-e29b-41d4-a716-446655440001",
			Topic:  "Test Topic",
			Emojis: []string{"😀", "😁", "😂"},
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

		input := roomUseCase.SetTopicInput{
			RoomID: "invalid-uuid",
			UserID: "550e8400-e29b-41d4-a716-446655440001",
			Topic:  "Test Topic",
			Emojis: []string{"😀"},
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

		input := roomUseCase.SetTopicInput{
			RoomID: "550e8400-e29b-41d4-a716-446655440000",
			UserID: "550e8400-e29b-41d4-a716-446655440001",
			Topic:  "Test Topic",
			Emojis: []string{"😀"},
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when room not found")
		}
	})

	t.Run("ホスト以外のユーザーがトピック設定しようとした場合はエラーが返されること", func(t *testing.T) {
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

		input := roomUseCase.SetTopicInput{
			RoomID: testRoom.ID().String(),
			UserID: "550e8400-e29b-41d4-a716-446655440001",
			Topic:  "Test Topic",
			Emojis: []string{"😀"},
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when non-host tries to set topic")
		}
		if err.Error() != "only host can set topic" {
			t.Errorf("Expected 'only host can set topic' error, got: %v", err)
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

		input := roomUseCase.SetTopicInput{
			RoomID: testRoom.ID().String(),
			UserID: "550e8400-e29b-41d4-a716-446655440001",
			Topic:  "Test Topic",
			Emojis: []string{"😀"},
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when room save fails")
		}
	})
}
