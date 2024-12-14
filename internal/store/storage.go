package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource already exists")
)

type PostsStorage interface {
	Create(context.Context, *Post) error
	GetAllPosts(ctx context.Context) ([]Post, error)
	GetPostByID(ctx context.Context, id int64) (*Post, error)
	UpdatePost(ctx context.Context, post *Post) error
	DeletePost(ctx context.Context, postID int64) error
	GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error)
}

type UsersStorage interface {
	Create(context.Context, *sql.Tx, *User) error
	GetUserByID(context.Context, int64) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, userID int64) error
	CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
	Activate(ctx context.Context, token string) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type CommentsStorage interface {
	GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	Create(context.Context, *Comment) error
}

type FollowersStorage interface {
	Follow(ctx context.Context, followerID, userID int64) error
	Unfollow(ctx context.Context, followerID, userID int64) error
}

type RolesStorage interface {
	GetByName(ctx context.Context, role string) (*Role, error)
}

type Storage struct {
	Posts     PostsStorage
	Users     UsersStorage
	Comments  CommentsStorage
	Followers FollowersStorage
	Roles     RolesStorage
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowersStore{db},
		Roles:     &RolesStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, f func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
