package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
)

// ParticipantRepository implements the participant.Repository interface
type ParticipantRepository struct {
	db *sql.DB
}

// NewParticipantRepository creates a new ParticipantRepository
func NewParticipantRepository(db *sql.DB) *ParticipantRepository {
	return &ParticipantRepository{db: db}
}

// Save persists a participant
func (r *ParticipantRepository) Save(ctx context.Context, p *participant.Participant) error {
	query := `
		INSERT INTO participants (id, room_id, user_id, role, is_leader, joined_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (room_id, user_id) DO UPDATE
		SET role = EXCLUDED.role,
			is_leader = EXCLUDED.is_leader
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		p.ID().String(),
		p.RoomID().String(),
		p.UserID().String(),
		p.Role().String(),
		p.IsLeader(),
		p.JoinedAt(),
	)

	return err
}

// FindByID retrieves a participant by ID
func (r *ParticipantRepository) FindByID(ctx context.Context, id participant.ParticipantID) (*participant.Participant, error) {
	query := `
		SELECT id, room_id, user_id, role, is_leader, joined_at
		FROM participants
		WHERE id = $1
	`

	return r.scanParticipant(ctx, query, id.String())
}

// FindByRoomID retrieves all participants in a room
func (r *ParticipantRepository) FindByRoomID(ctx context.Context, roomID participant.RoomID) ([]*participant.Participant, error) {
	query := `
		SELECT id, room_id, user_id, role, is_leader, joined_at
		FROM participants
		WHERE room_id = $1
		ORDER BY joined_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, roomID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*participant.Participant
	for rows.Next() {
		p, err := r.scanParticipantFromRows(rows)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, rows.Err()
}

// FindByRoomAndUser retrieves a specific participant by room and user
func (r *ParticipantRepository) FindByRoomAndUser(ctx context.Context, roomID participant.RoomID, userID participant.UserID) (*participant.Participant, error) {
	query := `
		SELECT id, room_id, user_id, role, is_leader, joined_at
		FROM participants
		WHERE room_id = $1 AND user_id = $2
	`

	var (
		id        string
		roomIDStr string
		userIDStr string
		role      string
		isLeader  bool
		joinedAt  interface{}
	)

	err := r.db.QueryRowContext(ctx, query, roomID.String(), userID.String()).Scan(
		&id, &roomIDStr, &userIDStr, &role, &isLeader, &joinedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("participant not found")
		}
		return nil, err
	}

	participantID, _ := participant.NewParticipantIDFromString(id)
	participantRoomID, _ := participant.NewRoomIDFromString(roomIDStr)
	participantUserID, _ := participant.NewUserIDFromString(userIDStr)
	participantRole, _ := participant.NewParticipantRoleFromString(role)

	p := participant.NewParticipant(participantID, participantRoomID, participantUserID, participantRole)
	if isLeader {
		p.SetAsLeader()
	}

	return p, nil
}

// scanParticipant scans a participant from a query result
func (r *ParticipantRepository) scanParticipant(ctx context.Context, query string, arg interface{}) (*participant.Participant, error) {
	var (
		id       string
		roomID   string
		userID   string
		role     string
		isLeader bool
		joinedAt interface{}
	)

	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&id, &roomID, &userID, &role, &isLeader, &joinedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("participant not found")
		}
		return nil, err
	}

	participantID, _ := participant.NewParticipantIDFromString(id)
	participantRoomID, _ := participant.NewRoomIDFromString(roomID)
	participantUserID, _ := participant.NewUserIDFromString(userID)
	participantRole, _ := participant.NewParticipantRoleFromString(role)

	p := participant.NewParticipant(participantID, participantRoomID, participantUserID, participantRole)
	if isLeader {
		p.SetAsLeader()
	}

	return p, nil
}

// scanParticipantFromRows scans a participant from rows
func (r *ParticipantRepository) scanParticipantFromRows(rows *sql.Rows) (*participant.Participant, error) {
	var (
		id       string
		roomID   string
		userID   string
		role     string
		isLeader bool
		joinedAt interface{}
	)

	err := rows.Scan(&id, &roomID, &userID, &role, &isLeader, &joinedAt)
	if err != nil {
		return nil, err
	}

	participantID, _ := participant.NewParticipantIDFromString(id)
	participantRoomID, _ := participant.NewRoomIDFromString(roomID)
	participantUserID, _ := participant.NewUserIDFromString(userID)
	participantRole, _ := participant.NewParticipantRoleFromString(role)

	p := participant.NewParticipant(participantID, participantRoomID, participantUserID, participantRole)
	if isLeader {
		p.SetAsLeader()
	}

	return p, nil
}

// Delete removes a participant
func (r *ParticipantRepository) Delete(ctx context.Context, roomID participant.RoomID, userID participant.UserID) error {
	query := `DELETE FROM participants WHERE room_id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, roomID.String(), userID.String())
	return err
}
