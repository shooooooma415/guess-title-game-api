package theme

import "context"

// Repository defines the interface for theme persistence
type Repository interface {
	// Save persists a theme
	Save(ctx context.Context, theme *Theme) error

	// FindByID retrieves a theme by ID
	FindByID(ctx context.Context, id ThemeID) (*Theme, error)

	// FindAll retrieves all themes
	FindAll(ctx context.Context) ([]*Theme, error)

	// Delete removes a theme
	Delete(ctx context.Context, id ThemeID) error
}
