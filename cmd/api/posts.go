package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/lucianboboc/goBackendEngineering/internal/store"
	"net/http"
	"strconv"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Content string   `json:"content" validate:"required,max=100"`
	Title   string   `json:"title" validate:"required,max=1000"`
	Tags    []string `json:"tags" validate:"required"`
}

type UpdatePostPayload struct {
	Content *string `json:"content" validate:"omitempty,max=100"`
	Title   *string `json:"title" validate:"omitempty,max=1000"`
}

// createPostsHandler godoc
//
//	@Summary		Create a post
//	@Description	Create a new post
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	store.Post
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/ [post]
func (app *application) createPostsHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := getUserFromCtx(r)
	post := &store.Post{
		Content: payload.Content,
		Title:   payload.Title,
		Tags:    payload.Tags,
		UserID:  user.ID,
	}
	err := app.store.Posts.Create(r.Context(), post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusCreated, post)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

// getPostsHandler godoc
//
//	@Summary		Fetches all posts
//	@Description	Fetches all posts
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	store.Post
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts [get]
func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.store.Posts.GetAllPosts(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	_ = app.jsonResponse(w, http.StatusOK, posts)
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err = app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
	}
}

// updatePostHandler godoc
//
//	@Summary		Updates a post
//	@Description	Updates a post by id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Post ID"
//	@Param			payload	body		UpdatePostPayload	true	"Post payload"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(&payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	err := app.store.Posts.UpdatePost(r.Context(), post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusOK, post)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	err := app.store.Posts.DeletePost(r.Context(), post.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	err = app.jsonResponse(w, http.StatusOK, post)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDStr := chi.URLParam(r, "post_id")
		postID, err := strconv.ParseInt(postIDStr, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		post, err := app.store.Posts.GetPostByID(r.Context(), postID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx := context.WithValue(r.Context(), postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	return r.Context().Value(postCtx).(*store.Post)
}
