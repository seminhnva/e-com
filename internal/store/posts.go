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
	GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]*PostWithMetadata, int, error)
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
	User      *User      `json:"user,omitempty"`
}

type PostWithMetadata struct {
	Post
	CommentCount int `json:"comment_count"`
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
	FROM posts
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (s *PostsStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]*PostWithMetadata, int, error) {
	query := `
	SELECT
		p.id, p.user_id, p.title, p.content, p.tags,
		p.version, p.created_at, p.updated_at,
		u.username,
		COUNT(c.id) AS comment_count,
		COUNT(*) OVER() AS total_count
	FROM posts p
	LEFT JOIN comments c ON c.post_id = p.id
	JOIN users u ON u.id = p.user_id
	JOIN followers f ON f.user_id = p.user_id AND f.follower_id = $1
	WHERE (
    $4 = ''
    OR p.title ILIKE '%' || $4 || '%'
    OR p.content ILIKE '%' || $4 || '%'
    OR EXISTS (
        SELECT 1 FROM unnest(p.tags) tag
        WHERE tag ILIKE '%' || $4 || '%'
    )
	)
	GROUP BY p.id, u.id, u.username
	ORDER BY p.created_at DESC
	LIMIT $2 OFFSET $3
	`
	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset(), fq.Search)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var posts []*PostWithMetadata
	var total int
	for rows.Next() {
		var p PostWithMetadata
		p.Post.User = &User{}
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			pq.Array(&p.Tags),
			&p.Version,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Post.User.Username,
			&p.CommentCount,
			&total,
		)
		if err != nil {
			return nil, 0, err
		}
		posts = append(posts, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}
