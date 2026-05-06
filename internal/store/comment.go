package store

import (
	"context"
	"database/sql"
)

type CommentUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type Comment struct {
	ID        int64       `json:"id"`
	PostID    int64       `json:"post_id"`
	UserID    int64       `json:"-"`
	Content   string      `json:"content"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	User      CommentUser `json:"user"`
}
type CommentsStore struct {
	db *sql.DB
}

type CommentsRepository interface {
	Create(context.Context, *Comment) error
	GetByPostID(context.Context, int64) ([]*Comment, error)
}

func (s *CommentsStore) GetByPostID(ctx context.Context, postID int64) ([]*Comment, error) {
	query := `
	SELECT c.id, c.post_id, c.user_id, c.content, u.id, u.username, c.created_at, c.updated_at
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC
	`
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*Comment{}
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.User.ID, &c.User.Username, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `
	INSERT INTO comments(post_id, user_id, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
