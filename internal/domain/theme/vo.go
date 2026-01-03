package theme

import (
	"errors"

	"github.com/shooooooma415/guess-title-game-api/utils"
)

// ThemeID represents a theme identifier
type ThemeID struct {
	value string
}

func NewThemeID() ThemeID {
	return ThemeID{value: utils.GenerateUUID()}
}

func NewThemeIDFromString(value string) (ThemeID, error) {
	if err := utils.ValidateUUID(value); err != nil {
		return ThemeID{}, err
	}
	return ThemeID{value: value}, nil
}

func (id ThemeID) String() string {
	return id.value
}

func (id ThemeID) Equals(other ThemeID) bool {
	return id.value == other.value
}

// ThemeTitle represents a theme title
type ThemeTitle struct {
	value string
}

func NewThemeTitle(value string) (ThemeTitle, error) {
	if value == "" {
		return ThemeTitle{}, errors.New("theme title cannot be empty")
	}
	if len(value) > 255 {
		return ThemeTitle{}, errors.New("theme title is too long")
	}
	return ThemeTitle{value: value}, nil
}

func (t ThemeTitle) String() string {
	return t.value
}

func (t ThemeTitle) Equals(other ThemeTitle) bool {
	return t.value == other.value
}

// Hint represents a theme hint
type Hint struct {
	value string
}

func NewHint(value string) Hint {
	return Hint{value: value}
}

func (h Hint) String() string {
	return h.value
}

func (h Hint) IsEmpty() bool {
	return h.value == ""
}
