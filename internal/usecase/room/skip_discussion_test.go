package room_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
	roomUseCase "github.com/shooooooma415/guess-title-game-api/internal/usecase/room"
)

func TestSkipDiscussionUseCaseExecute(t *testing.T) {
	type fixture struct {
		useCase         *roomUseCase.SkipDiscussionUseCase
		roomRepo        *mockRoomRepository
		participantRepo *mockParticipantRepository
		eventPublisher  *mockEventPublisher
	}

	newFixture := func(t *testing.T) *fixture {
		t.Helper()

		roomRepo := &mockRoomRepository{}
		participantRepo := &mockParticipantRepository{}
		eventPublisher := &mockEventPublisher{}

		useCase := roomUseCase.NewSkipDiscussionUseCase(
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

	createTestRoomWithDummyData := func() *room.Room {
		roomID := room.NewRoomID()
		roomCode := room.NewRoomCode()
		themeID, _ := room.NewThemeIDFromString("550e8400-e29b-41d4-a716-446655440000")
		hostUserID, _ := room.NewHostUserIDFromString("550e8400-e29b-41d4-a716-446655440001")
		r := room.NewRoom(roomID, roomCode, themeID, hostUserID)
		r.Start()
		r.ChangeStatus(room.StatusDiscussing)

		// Set dummy data
		origEmojis := room.NewEmojiList([]string{"😀", "😁", "😂"})
		dispEmojis := room.NewEmojiList([]string{"😀", "😁", "😂", "😃"})
		dummyIdx, _ := room.NewDummyIndex(3)
		dummyEmoji, _ := room.NewDummyEmoji("😃")
		r.SetGameData(origEmojis, dispEmojis, dummyIdx, dummyEmoji)

		return r
	}

	createHostParticipant := func(roomID, userID string) *participant.Participant {
		participantID := participant.NewParticipantID()
		participantRoomID, _ := participant.NewRoomIDFromString(roomID)
		participantUserID, _ := participant.NewUserIDFromString(userID)
		return participant.NewParticipant(participantID, participantRoomID, participantUserID, participant.RoleHost)
	}

	t.Run("正常にディスカッションがスキップされること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoomWithDummyData()
		testParticipant := createHostParticipant(testRoom.ID().String(), "550e8400-e29b-41d4-a716-446655440001")

		f.roomRepo.findByIDFunc = func(ctx context.Context, id room.RoomID) (*room.Room, error) {
			return testRoom, nil
		}
		f.participantRepo.findByRoomAndUserFunc = func(ctx context.Context, roomID participant.RoomID, userID participant.UserID) (*participant.Participant, error) {
			return testParticipant, nil
		}

		input := roomUseCase.SkipDiscussionInput{
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

		input := roomUseCase.SkipDiscussionInput{
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

		input := roomUseCase.SkipDiscussionInput{
			RoomID: "550e8400-e29b-41d4-a716-446655440000",
			UserID: "550e8400-e29b-41d4-a716-446655440001",
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when room not found")
		}
	})

	t.Run("ホスト以外のユーザーがスキップしようとした場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoomWithDummyData()

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

		input := roomUseCase.SkipDiscussionInput{
			RoomID: testRoom.ID().String(),
			UserID: "550e8400-e29b-41d4-a716-446655440001",
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when non-host tries to skip discussion")
		}
		if err.Error() != "only host or leader can skip discussion" {
			t.Errorf("Expected 'only host or leader can skip discussion' error, got: %v", err)
		}
	})

	t.Run("ダミーデータが設定されていない場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		roomID := room.NewRoomID()
		roomCode := room.NewRoomCode()
		themeID, _ := room.NewThemeIDFromString("550e8400-e29b-41d4-a716-446655440000")
		hostUserID, _ := room.NewHostUserIDFromString("550e8400-e29b-41d4-a716-446655440001")
		testRoom := room.NewRoom(roomID, roomCode, themeID, hostUserID)
		testRoom.Start()
		testRoom.ChangeStatus(room.StatusDiscussing)
		// Dummy data not set

		testParticipant := createHostParticipant(testRoom.ID().String(), "550e8400-e29b-41d4-a716-446655440001")

		f.roomRepo.findByIDFunc = func(ctx context.Context, id room.RoomID) (*room.Room, error) {
			return testRoom, nil
		}
		f.participantRepo.findByRoomAndUserFunc = func(ctx context.Context, roomID participant.RoomID, userID participant.UserID) (*participant.Participant, error) {
			return testParticipant, nil
		}

		input := roomUseCase.SkipDiscussionInput{
			RoomID: testRoom.ID().String(),
			UserID: "550e8400-e29b-41d4-a716-446655440001",
		}

		// act
		err := f.useCase.Execute(context.Background(), input)

		// assert
		if err == nil {
			t.Error("Expected error when dummy data is not set")
		} else if err.Error() != "dummy data is required before skipping discussion" {
			t.Errorf("Expected 'dummy data is required before skipping discussion' error, got: %v", err)
		}
	})

	t.Run("Roomの保存に失敗した場合はエラーが返されること", func(t *testing.T) {
		// arrange
		f := newFixture(t)
		testRoom := createTestRoomWithDummyData()
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

		input := roomUseCase.SkipDiscussionInput{
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
