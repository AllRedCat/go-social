package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Contract
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	UpdateAvatar(ctx context.Context, id uint, avatarURL string) error
	UpdateUser(ctx context.Context, user *User) error
	SoftDelete(ctx context.Context, id uint) error
}

// Structure and Constructor
type sqliteRepository struct {
	db *sql.DB
}

// New instance of SQLite repository
func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{
		db: db,
	}
}

// Queries
func (r *sqliteRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (name, email, password, avatar_url, created_at)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.AvatarURL,
		user.CreatedAt,
	).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("Error on create new user: %w", err)
	}

	return nil
}

func (r *sqliteRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, email, password, avatar_url, created_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("User not found: %w", err)
		}
		// Faltava esse retorno para tratar outros erros de banco!
		return nil, fmt.Errorf("Error on fetch user by email: %w", err)
	}

	return user, nil
}

func (r *sqliteRepository) UpdateAvatar(ctx context.Context, id uint, avatarURL string) error {
	query := `UPDATE users SET avatar_url = ? WHERE id = ? AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, avatarURL, id)
	if err != nil {
		return fmt.Errorf("Error on update avatar: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("User not found or deleted")
	}

	return nil
}

func (r *sqliteRepository) UpdateUser(ctx context.Context, user *User) error {
	query := `UPDATE users SET name = ?, email = ? WHERE id = ? AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("Error on update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("User not found or deleted")
	}

	return nil
}

func (r *sqliteRepository) SoftDelete(ctx context.Context, id uint) error {
	now := time.Now()
	query := `UPDATE users SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, now, id)
	if err != nil {
		return fmt.Errorf("Error on soft delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("User not found or already deleted")
	}

	return nil
}
