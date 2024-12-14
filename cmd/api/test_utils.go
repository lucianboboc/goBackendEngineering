package main

import (
	"github.com/lucianboboc/goBackendEngineering/internal/auth"
	"github.com/lucianboboc/goBackendEngineering/internal/store"
	"github.com/lucianboboc/goBackendEngineering/internal/store/cache"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()
	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, got %d", expected, actual)
	}
}
