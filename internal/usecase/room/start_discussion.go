package room

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// StartDiscussionUseCase starts the discussion phase by setting game data and changing room status
type StartDiscussionUseCase struct {
	roomRepo        room.Repository
	participantRepo participant.Repository
}

// NewStartDiscussionUseCase creates a new StartDiscussionUseCase
func NewStartDiscussionUseCase(roomRepo room.Repository, participantRepo participant.Repository) *StartDiscussionUseCase {
	return &StartDiscussionUseCase{
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
	}
}

// StartDiscussionInput represents input for starting discussion
type StartDiscussionInput struct {
	RoomID          string
	OriginalEmojis  []string
	DisplayedEmojis []string
	DummyIndex      int
	DummyEmoji      string
}

// Execute starts the discussion phase
func (uc *StartDiscussionUseCase) Execute(ctx context.Context, input StartDiscussionInput) error {
	// Validate input
	if input.RoomID == "" {
		return errors.New("room ID is required")
	}

	// Convert room ID
	roomID, err := room.NewRoomIDFromString(input.RoomID)
	if err != nil {
		return err
	}

	// Find room
	foundRoom, err := uc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return errors.New("room not found")
	}

	// Validate and create game data
	originalEmojis := room.NewEmojiList(input.OriginalEmojis)
	displayedEmojis := room.NewEmojiList(input.DisplayedEmojis)

	dummyIndex, err := room.NewDummyIndex(input.DummyIndex)
	if err != nil {
		return err
	}

	dummyEmoji, err := room.NewDummyEmoji(input.DummyEmoji)
	if err != nil {
		return err
	}

	// Set game data
	fmt.Printf("[StartDiscussion] Setting game data\n")
	if err := foundRoom.SetGameData(
		originalEmojis,
		displayedEmojis,
		dummyIndex,
		dummyEmoji,
	); err != nil {
		fmt.Printf("[StartDiscussion] Failed to set game data: %v\n", err)
		return err
	}
	fmt.Printf("[StartDiscussion] Game data set successfully\n")

	// Generate emoji assignments for players (excluding host)
	participantRoomID, _ := participant.NewRoomIDFromString(input.RoomID)
	participants, err := uc.participantRepo.FindByRoomID(ctx, participantRoomID)
	if err != nil {
		return fmt.Errorf("failed to fetch participants: %w", err)
	}

	// Filter out host and build assignments
	type Assignment struct {
		UserID string `json:"user_id"`
		Emoji  string `json:"emoji"`
	}

	assignments := []Assignment{}
	emojiIndex := 0
	for _, p := range participants {
		if p.Role() != participant.RoleHost {
			if emojiIndex < len(input.DisplayedEmojis) {
				assignments = append(assignments, Assignment{
					UserID: p.UserID().String(),
					Emoji:  input.DisplayedEmojis[emojiIndex],
				})
				emojiIndex++
			}
		}
	}

	// Convert assignments to JSON string array
	assignmentsJSON := []string{}
	for _, assignment := range assignments {
		jsonBytes, _ := json.Marshal(assignment)
		assignmentsJSON = append(assignmentsJSON, string(jsonBytes))
	}

	fmt.Printf("[StartDiscussion] Generated %d assignments: %v\n", len(assignmentsJSON), assignmentsJSON)

	// Set assignments
	fmt.Printf("[StartDiscussion] Setting assignments\n")
	if err := foundRoom.SetAssignments(room.NewAssignments(assignmentsJSON)); err != nil {
		fmt.Printf("[StartDiscussion] Failed to set assignments: %v\n", err)
		return err
	}
	fmt.Printf("[StartDiscussion] Assignments set successfully\n")

	// Note: Status change to 'discussing' is already handled by SetTopicUseCase (HTTP endpoint)
	// This WebSocket handler only sets game data and assignments
	//
	// IMPORTANT: Fetch the latest room data again to avoid overwriting topic set by SetTopicUseCase
	// This is necessary because SetTopicUseCase (HTTP) and StartDiscussionUseCase (WebSocket)
	// are called almost simultaneously, causing a race condition.
	fmt.Printf("[StartDiscussion] Re-fetching room to get latest data (including topic)\n")
	latestRoom, err := uc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		fmt.Printf("[StartDiscussion] Failed to re-fetch room: %v\n", err)
		return errors.New("failed to re-fetch room")
	}

	topicStr := ""
	if latestRoom.Topic() != nil {
		topicStr = latestRoom.Topic().String()
	}
	fmt.Printf("[StartDiscussion] Re-fetched room topic: '%s', status: %s\n", topicStr, latestRoom.Status().String())

	// Re-apply game data and assignments to the latest room
	latestRoom.SetGameData(originalEmojis, displayedEmojis, dummyIndex, dummyEmoji)
	latestRoom.SetAssignments(room.NewAssignments(assignmentsJSON))

	fmt.Printf("[StartDiscussion] Saving room with game data and assignments. Current status: %s\n", latestRoom.Status().String())
	if err := uc.roomRepo.Save(ctx, latestRoom); err != nil {
		fmt.Printf("[StartDiscussion] Save failed: %v\n", err)
		return err
	}
	fmt.Printf("[StartDiscussion] Room saved successfully with game data\n")

	return nil
}
