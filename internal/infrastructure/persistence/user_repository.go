package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/shooooooma415/guess-title-game-api/internal/domain/user"
)

// UserRepository implements the user.Repository interface
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Save persists a user
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (id, name, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		u.ID().String(),
		u.Name().String(),
		u.CreatedAt(),
	)

	return err
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	query := `
		SELECT id, name, created_at
		FROM users
		WHERE id = $1
	`

	var (
		userID    string
		name      string
		createdAt interface{}
	)

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(&userID, &name, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	uid, err := user.NewUserIDFromString(userID)
	if err != nil {
		return nil, err
	}

	userName, err := user.NewUserName(name)
	if err != nil {
		return nil, err
	}

	return user.NewUser(uid, userName), nil
}

// Delete removes a user
func (r *UserRepository) Delete(ctx context.Context, id user.UserID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}
