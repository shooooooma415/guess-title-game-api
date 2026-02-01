package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// RoomRepository implements the room.Repository interface
type RoomRepository struct {
	db *sql.DB
}

// NewRoomRepository creates a new RoomRepository
func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// Save persists a room
func (r *RoomRepository) Save(ctx context.Context, rm *room.Room) error {
	query := `
		INSERT INTO rooms (
			id, code, theme_id, topic, answer, status, host_user_id,
			created_at, started_at, original_emojis, displayed_emojis,
			dummy_index, dummy_emoji, assignments
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (id) DO UPDATE
		SET topic = EXCLUDED.topic,
			answer = EXCLUDED.answer,
			status = EXCLUDED.status,
			started_at = EXCLUDED.started_at,
			original_emojis = EXCLUDED.original_emojis,
			displayed_emojis = EXCLUDED.displayed_emojis,
			dummy_index = EXCLUDED.dummy_index,
			dummy_emoji = EXCLUDED.dummy_emoji,
			assignments = EXCLUDED.assignments
	`

	// Convert VOs to primitive values
	var topicStr, answerStr interface{}
	if rm.Topic() != nil {
		topicStr = rm.Topic().String()
		fmt.Printf("[RoomRepository.Save] Topic is not nil: '%s'\n", topicStr)
	} else {
		fmt.Printf("[RoomRepository.Save] Topic is nil\n")
	}
	if rm.Answer() != nil {
		answerStr = rm.Answer().String()
	}

	var originalEmojis, displayedEmojis []string
	if rm.OriginalEmojis() != nil {
		originalEmojis = rm.OriginalEmojis().Values()
	}
	if rm.DisplayedEmojis() != nil {
		displayedEmojis = rm.DisplayedEmojis().Values()
	}

	var dummyIndex interface{}
	if rm.DummyIndex() != nil {
		dummyIndex = rm.DummyIndex().Value()
	}

	var dummyEmoji interface{}
	if rm.DummyEmoji() != nil {
		dummyEmoji = rm.DummyEmoji().String()
	}

	var assignments []string
	if rm.Assignments() != nil {
		assignments = rm.Assignments().Values()
	}

	fmt.Printf("[RoomRepository.Save] Executing SQL with params:\n")
	fmt.Printf("  ID: %s\n", rm.ID().String())
	fmt.Printf("  Code: %s\n", rm.Code().String())
	fmt.Printf("  ThemeID: %s\n", rm.ThemeID().String())
	fmt.Printf("  Topic: %v (type: %T)\n", topicStr, topicStr)
	fmt.Printf("  Answer: %v (type: %T)\n", answerStr, answerStr)
	fmt.Printf("  Status: %s\n", rm.Status().String())
	fmt.Printf("  HostUserID: %s\n", rm.HostUserID().String())

	_, err := r.db.ExecContext(
		ctx,
		query,
		rm.ID().String(),
		rm.Code().String(),
		rm.ThemeID().String(),
		topicStr,
		answerStr,
		rm.Status().String(),
		rm.HostUserID().String(),
		rm.CreatedAt(),
		rm.StartedAt(),
		pq.Array(originalEmojis),
		pq.Array(displayedEmojis),
		dummyIndex,
		dummyEmoji,
		pq.Array(assignments),
	)

	if err != nil {
		fmt.Printf("[RoomRepository.Save] SQL execution failed: %v\n", err)
	} else {
		fmt.Printf("[RoomRepository.Save] SQL execution succeeded\n")
	}

	return err
}

// FindByID retrieves a room by ID
func (r *RoomRepository) FindByID(ctx context.Context, id room.RoomID) (*room.Room, error) {
	query := `
		SELECT id, code, theme_id, topic, answer, status, host_user_id,
			created_at, started_at, original_emojis, displayed_emojis,
			dummy_index, dummy_emoji, assignments
		FROM rooms
		WHERE id = $1
	`

	return r.scanRoom(ctx, query, id.String())
}

// FindByCode retrieves a room by code
func (r *RoomRepository) FindByCode(ctx context.Context, code room.RoomCode) (*room.Room, error) {
	query := `
		SELECT id, code, theme_id, topic, answer, status, host_user_id,
			created_at, started_at, original_emojis, displayed_emojis,
			dummy_index, dummy_emoji, assignments
		FROM rooms
		WHERE code = $1
	`

	return r.scanRoom(ctx, query, code.String())
}

// scanRoom scans a room from a query result
func (r *RoomRepository) scanRoom(ctx context.Context, query string, arg interface{}) (*room.Room, error) {
	var (
		id              string
		code            string
		themeID         string
		topic           sql.NullString
		answer          sql.NullString
		status          string
		hostUserID      string
		createdAt       interface{}
		startedAt       interface{}
		originalEmojis  []string
		displayedEmojis []string
		dummyIndex      sql.NullInt64
		dummyEmoji      sql.NullString
		assignments     []string
	)

	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&id, &code, &themeID, &topic, &answer, &status, &hostUserID,
		&createdAt, &startedAt,
		pq.Array(&originalEmojis), pq.Array(&displayedEmojis),
		&dummyIndex, &dummyEmoji, pq.Array(&assignments),
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("room not found")
		}
		return nil, err
	}

	fmt.Printf("[RoomRepository.scanRoom] Scanned from DB:\n")
	fmt.Printf("  ID: %s\n", id)
	fmt.Printf("  Topic Valid: %v, Topic Value: '%s'\n", topic.Valid, topic.String)
	fmt.Printf("  Answer Valid: %v, Answer Value: '%s'\n", answer.Valid, answer.String)
	fmt.Printf("  Status: %s\n", status)

	roomID, _ := room.NewRoomIDFromString(id)
	roomCode, _ := room.NewRoomCodeFromString(code)
	roomThemeID, _ := room.NewThemeIDFromString(themeID)
	roomHostUserID, _ := room.NewHostUserIDFromString(hostUserID)

	rm := room.NewRoom(roomID, roomCode, roomThemeID, roomHostUserID)

	if topic.Valid {
		if t, err := room.NewTopic(topic.String); err == nil {
			rm.SetTopicUnchecked(&t)
		}
	}
	if answer.Valid {
		if a, err := room.NewAnswer(answer.String); err == nil {
			rm.SetAnswer(a)
		}
	}

	roomStatus, _ := room.NewRoomStatusFromString(status)
	rm.SetStatus(roomStatus) // Use SetStatus instead of ChangeStatus to bypass validation when reconstructing from DB

	if len(originalEmojis) > 0 && len(displayedEmojis) > 0 && dummyIndex.Valid {
		origEmojis := room.NewEmojiList(originalEmojis)
		dispEmojis := room.NewEmojiList(displayedEmojis)
		dummyIdx, _ := room.NewDummyIndex(int(dummyIndex.Int64))
		dummyEmo, _ := room.NewDummyEmoji(dummyEmoji.String)
		rm.SetGameData(origEmojis, dispEmojis, dummyIdx, dummyEmo)
	}

	if len(assignments) > 0 {
		rm.SetAssignments(room.NewAssignments(assignments))
	}

	return rm, nil
}

// Delete removes a room
func (r *RoomRepository) Delete(ctx context.Context, id room.RoomID) error {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}
