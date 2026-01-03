package user_test

import (
	"context"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/user"
)

// Mock User Repository
type mockUserRepository struct {
	saveFunc     func(context.Context, *user.User) error
	findByIDFunc func(context.Context, user.UserID) (*user.User, error)
	deleteFunc   func(context.Context, user.UserID) error
}

func (m *mockUserRepository) Save(ctx context.Context, u *user.User) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, u)
	}
	return nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) Delete(ctx context.Context, id user.UserID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

// Mock Room Repository
type mockRoomRepository struct {
	saveFunc       func(context.Context, *room.Room) error
	findByIDFunc   func(context.Context, room.RoomID) (*room.Room, error)
	findByCodeFunc func(context.Context, room.RoomCode) (*room.Room, error)
	deleteFunc     func(context.Context, room.RoomID) error
}

func (m *mockRoomRepository) Save(ctx context.Context, r *room.Room) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, r)
	}
	return nil
}

func (m *mockRoomRepository) FindByID(ctx context.Context, id room.RoomID) (*room.Room, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRoomRepository) FindByCode(ctx context.Context, code room.RoomCode) (*room.Room, error) {
	if m.findByCodeFunc != nil {
		return m.findByCodeFunc(ctx, code)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRoomRepository) Delete(ctx context.Context, id room.RoomID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

// Mock Participant Repository
type mockParticipantRepository struct {
	saveFunc              func(context.Context, *participant.Participant) error
	findByIDFunc          func(context.Context, participant.ParticipantID) (*participant.Participant, error)
	findByRoomIDFunc      func(context.Context, participant.RoomID) ([]*participant.Participant, error)
	findByRoomAndUserFunc func(context.Context, participant.RoomID, participant.UserID) (*participant.Participant, error)
	deleteFunc            func(context.Context, participant.RoomID, participant.UserID) error
}

func (m *mockParticipantRepository) Save(ctx context.Context, p *participant.Participant) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, p)
	}
	return nil
}

func (m *mockParticipantRepository) FindByID(ctx context.Context, id participant.ParticipantID) (*participant.Participant, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockParticipantRepository) FindByRoomID(ctx context.Context, roomID participant.RoomID) ([]*participant.Participant, error) {
	if m.findByRoomIDFunc != nil {
		return m.findByRoomIDFunc(ctx, roomID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockParticipantRepository) FindByRoomAndUser(ctx context.Context, roomID participant.RoomID, userID participant.UserID) (*participant.Participant, error) {
	if m.findByRoomAndUserFunc != nil {
		return m.findByRoomAndUserFunc(ctx, roomID, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockParticipantRepository) Delete(ctx context.Context, roomID participant.RoomID, userID participant.UserID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, roomID, userID)
	}
	return errors.New("not implemented")
}
