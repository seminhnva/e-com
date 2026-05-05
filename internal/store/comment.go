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
