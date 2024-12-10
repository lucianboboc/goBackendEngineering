package main

import (
	"github.com/lucianboboc/goBackendEngineering/internal/store"
	"net/http"
)

type CommentPayload struct {
	UserID  int64  `json:"user_id" validate:"required"`
	Content string `json:"content" validate:"required,max=100"`
}

func (app *application) getCommentsByPost(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusOK, comments)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) createPostComment(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload CommentPayload
	err := readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err = Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  payload.UserID,
		Content: payload.Content,
	}
	err = app.store.Comments.Create(r.Context(), comment)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusCreated, comment)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
