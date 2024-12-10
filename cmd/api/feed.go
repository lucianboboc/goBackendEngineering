package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// pagination, filter

	feed, err := app.store.Posts.GetUserFeed(r.Context(), int64(42))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	err = app.jsonResponse(w, http.StatusOK, feed)
	if err != nil {
		app.internalServerError(w, r, err)
	}
}
