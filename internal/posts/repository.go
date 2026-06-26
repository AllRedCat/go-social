package posts

import (
	"context"
	"database/sql"
	"fmt"
)

// Repository - Contract
type Repository interface {
	Create(ctx context.Context, post *Post) error
}

// Structure and Constructor
type sqliteRepository struct {
	db *sql.DB
}

// NewRepository - Create a new instance of SQLite repository
func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{
		db: db,
	}
}

func (r *sqliteRepository) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (user_id, title, content, image_url, created_at)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		post.UserId,
		post.Title,
		post.Content,
		post.ImageURL,
		post.CreatedAt,
	).Scan(&post.Id)
	if err != nil {
		return fmt.Errorf("error on create new post: %w", err)
	}

	return nil
}
