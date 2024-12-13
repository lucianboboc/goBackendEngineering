package main

import (
	"log/slog"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(
		"internal server error",
		slog.Any("method", r.Method),
		slog.Any("path", r.URL.Path),
		slog.Any("error", err.Error()),
	)
	_ = writeJSONError(w, http.StatusInternalServerError, "The server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"bad request error",
		slog.Any("method", r.Method),
		slog.Any("path", r.URL.Path),
		slog.Any("error", err.Error()),
	)
	_ = writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn(
		"not found error",
		slog.Any("method", r.Method),
		slog.Any("path", r.URL.Path),
		slog.Any("error", err.Error()),
	)
	_ = writeJSONError(w, http.StatusNotFound, err.Error())
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(
		"conflict error",
		slog.Any("method", r.Method),
		slog.Any("path", r.URL.Path),
		slog.Any("error", err.Error()),
	)
	_ = writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(
		"unauthorized error",
		slog.Any("method", r.Method),
		slog.Any("path", r.URL.Path),
		slog.Any("error", err),
	)
	_ = writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(
		"unauthorized basic error",
		slog.Any("method", r.Method),
		slog.Any("path", r.URL.Path),
		slog.Any("error", err),
	)

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)

	_ = writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
