package room

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/participant"
	"github.com/shooooooma415/guess-title-game-api/internal/domain/room"
)

// SetTopicInput represents the input for setting a topic
type SetTopicInput struct {
	RoomID          string
	UserID          string
	Topic           string
	Emojis          []string
	DisplayedEmojis []string
	OriginalEmojis  []string
	DummyIndex      int
	DummyEmoji      string
}

// SetTopicUseCase handles the logic for setting a topic
type SetTopicUseCase struct {
	roomRepo        room.Repository
	participantRepo participant.Repository
}

// NewSetTopicUseCase creates a new SetTopicUseCase
func NewSetTopicUseCase(
	roomRepo room.Repository,
	participantRepo participant.Repository,
) *SetTopicUseCase {
	return &SetTopicUseCase{
		roomRepo:        roomRepo,
		participantRepo: participantRepo,
	}
}

// Execute sets a topic for the room
func (uc *SetTopicUseCase) Execute(ctx context.Context, input SetTopicInput) error {
	// Find room
	roomID, err := room.NewRoomIDFromString(input.RoomID)
	if err != nil {
		return err
	}

	foundRoom, err := uc.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return errors.New("room not found")
	}

	// Verify user is host
	participantRoomID, _ := participant.NewRoomIDFromString(input.RoomID)
	participantUserID, _ := participant.NewUserIDFromString(input.UserID)

	foundParticipant, err := uc.participantRepo.FindByRoomAndUser(ctx, participantRoomID, participantUserID)
	if err != nil {
		return errors.New("participant not found")
	}

	if foundParticipant.Role() != participant.RoleHost {
		return errors.New("only host can set topic")
	}

	// Set topic
	fmt.Printf("[SetTopic] Received topic from input: '%s'\n", input.Topic)
	topic, err := room.NewTopic(input.Topic)
	if err != nil {
		fmt.Printf("[SetTopic] Failed to create topic: %v\n", err)
		return err
	}
	if err := foundRoom.SetTopic(topic); err != nil {
		fmt.Printf("[SetTopic] Failed to set topic: %v\n", err)
		return err
	}
	fmt.Printf("[SetTopic] Successfully set topic: '%s'\n", topic.String())

	// Set game data if provided (dummy emoji information)
	if len(input.DisplayedEmojis) > 0 && len(input.OriginalEmojis) > 0 {
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

		foundRoom.SetGameData(originalEmojis, displayedEmojis, dummyIndex, dummyEmoji)

		// Generate emoji assignments for players
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

		fmt.Printf("[SetTopic] Generated %d assignments: %v\n", len(assignmentsJSON), assignmentsJSON)

		// Set assignments
		foundRoom.SetAssignments(room.NewAssignments(assignmentsJSON))

		// Change status to discussing
		fmt.Printf("[SetTopic] Current status: %s, attempting to change to discussing\n", foundRoom.Status().String())
		if err := foundRoom.ChangeStatus(room.StatusDiscussing); err != nil {
			fmt.Printf("[SetTopic] Failed to change status: %v\n", err)
			return err
		}
		fmt.Printf("[SetTopic] Successfully changed status to discussing\n")
	}

	// Save room
	topicBeforeSave := ""
	if foundRoom.Topic() != nil {
		topicBeforeSave = foundRoom.Topic().String()
	}
	fmt.Printf("[SetTopic] About to save room with status: %s, topic: '%s'\n", foundRoom.Status().String(), topicBeforeSave)
	if err := uc.roomRepo.Save(ctx, foundRoom); err != nil {
		fmt.Printf("[SetTopic] Save failed: %v\n", err)
		return err
	}
	fmt.Printf("[SetTopic] Room saved successfully\n")

	return nil
}
