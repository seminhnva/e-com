package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type PostsRepository interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *Post) error
}

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Version   int        `json:"version"`
	Comments  []*Comment `json:"comments,omitempty"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts(content, title , user_id, tags)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags)).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, postID int64) (*Post, error) {
	query := `
	SELECT id, user_id, content, title, tags, version, created_at, updated_at
	FROM posts, pg_sleep(2)
	WHERE id = $1
	`
	post := &Post{}
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&post.ID, &post.UserID, &post.Content, &post.Title, pq.Array(&post.Tags), &post.Version, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return post, nil
}

func (s *PostsStore) Delete(ctx context.Context, postID int64) error {
	query := `
	DELETE FROM posts
	WHERE id = $1
	`
	result, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostsStore) Update(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts
	SET content = $1, title = $2, tags = $3, updated_at = NOW(), version = version + 1
	WHERE id = $4 AND version = $5
	`
	result, err := s.db.ExecContext(ctx, query, post.Content, post.Title, pq.Array(post.Tags), post.ID, post.Version)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrEditConflict
	}
	return nil
}
