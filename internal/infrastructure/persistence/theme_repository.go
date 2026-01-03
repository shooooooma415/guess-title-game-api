package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/theme"
)

// ThemeRepository implements the theme.Repository interface
type ThemeRepository struct {
	db *sql.DB
}

// NewThemeRepository creates a new ThemeRepository
func NewThemeRepository(db *sql.DB) *ThemeRepository {
	return &ThemeRepository{db: db}
}

// Save persists a theme
func (r *ThemeRepository) Save(ctx context.Context, t *theme.Theme) error {
	query := `
		INSERT INTO themes (id, title, hint)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET title = EXCLUDED.title,
			hint = EXCLUDED.hint
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		t.ID().String(),
		t.Title().String(),
		t.Hint().String(),
	)

	return err
}

// FindByID retrieves a theme by ID
func (r *ThemeRepository) FindByID(ctx context.Context, id theme.ThemeID) (*theme.Theme, error) {
	query := `
		SELECT id, title, hint
		FROM themes
		WHERE id = $1
	`

	var (
		themeID string
		title   string
		hint    sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&themeID, &title, &hint)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("theme not found")
		}
		return nil, err
	}

	tid, _ := theme.NewThemeIDFromString(themeID)
	themeTitle, _ := theme.NewThemeTitle(title)

	hintStr := ""
	if hint.Valid {
		hintStr = hint.String
	}
	hintVO := theme.NewHint(hintStr)

	return theme.NewTheme(tid, themeTitle, hintVO), nil
}

// FindAll retrieves all themes
func (r *ThemeRepository) FindAll(ctx context.Context) ([]*theme.Theme, error) {
	query := `
		SELECT id, title, hint
		FROM themes
		ORDER BY title ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var themes []*theme.Theme
	for rows.Next() {
		var (
			themeID string
			title   string
			hint    sql.NullString
		)

		if err := rows.Scan(&themeID, &title, &hint); err != nil {
			return nil, err
		}

		tid, _ := theme.NewThemeIDFromString(themeID)
		themeTitle, _ := theme.NewThemeTitle(title)

		hintStr := ""
		if hint.Valid {
			hintStr = hint.String
		}
		hintVO := theme.NewHint(hintStr)

		themes = append(themes, theme.NewTheme(tid, themeTitle, hintVO))
	}

	return themes, rows.Err()
}

// Delete removes a theme
func (r *ThemeRepository) Delete(ctx context.Context, id theme.ThemeID) error {
	query := `DELETE FROM themes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}
