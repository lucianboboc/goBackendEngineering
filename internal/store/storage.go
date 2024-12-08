package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type PostsStorage interface {
	Create(context.Context, *Post) error
	GetAllPosts(ctx context.Context) ([]Post, error)
	GetPostByID(ctx context.Context, id int64) (*Post, error)
	UpdatePost(ctx context.Context, post *Post) error
	DeletePost(ctx context.Context, postID int64) error
}

type UsersStorage interface {
	Create(context.Context, *User) error
}

type CommentsStorage interface {
	GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	Create(context.Context, *Comment) error
}

type Storage struct {
	Posts    PostsStorage
	Users    UsersStorage
	Comments CommentsStorage
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}
