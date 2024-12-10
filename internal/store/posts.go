package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Version   int        `json:"version"`
	Comments  []Comment  `json:"comments"`
	User      User       `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags) 
	VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostsStore) GetAllPosts(ctx context.Context) ([]Post, error) {
	query := `SELECT id, title, user_id, content, tags, created_at, updated_at, version FROM posts`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	posts := make([]Post, 0)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.UserID,
			&post.Content,
			pq.Array(&post.Tags),
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Version,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostsStore) GetPostByID(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, title, user_id, content, tags, created_at, updated_at, version FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.UserID,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostsStore) UpdatePost(ctx context.Context, post *Post) error {
	query := `UPDATE posts 
	SET content = $1, title = $2, version = $3
	WHERE id = $4
	AND version = $5
	RETURNING version`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.Version+1,
		post.ID,
		post.Version,
	).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *PostsStore) DeletePost(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostsStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	query := `SELECT p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
    u.username, COUNT(c.id) as comments_count
	FROM posts AS p
	LEFT JOIN comments AS c ON c.post_id = p.id
	LEFT JOIN users AS u ON p.user_id = u.id
	JOIN followers AS f ON f.follower_id = p.user_id OR p.user_id = $1
	WHERE f.user_id = $1 OR p.user_id = $1
	GROUP BY p.id, u.username
	ORDER BY p.created_at desc
	LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]PostWithMetadata, 0)
	for rows.Next() {
		var postWithMetadata PostWithMetadata
		err := rows.Scan(
			&postWithMetadata.ID,
			&postWithMetadata.UserID,
			&postWithMetadata.Title,
			&postWithMetadata.Content,
			&postWithMetadata.CreatedAt,
			&postWithMetadata.Version,
			pq.Array(&postWithMetadata.Tags),
			&postWithMetadata.User.Username,
			&postWithMetadata.CommentsCount,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, postWithMetadata)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
