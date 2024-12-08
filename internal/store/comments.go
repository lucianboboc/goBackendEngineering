package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64      `json:"id"`
	PostID    int64      `json:"post_id"`
	UserID    int64      `json:"user_id"`
	Content   string     `json:"content"`
	CreatedAt *time.Time `json:"created_at"`
	User      User       `json:"user"`
}

type CommentsStore struct {
	db *sql.DB
}

func (c *CommentsStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id, users.email FROM comments AS c
	JOIN users ON users.id = c.user_id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC`

	rows, err := c.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		c := Comment{}
		c.User = User{}
		err = rows.Scan(
			&c.ID,
			&c.PostID,
			&c.UserID,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
			&c.User.ID,
			&c.User.Email,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) 
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := c.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
