package main

import (
	"net/http"
	"testing"
	"time"
)

func TestGetUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()

	testToken, _ := app.authenticator.GenerateToken(0, "", "", time.Hour)

	t.Run("it should not allowed unauthenticated request", func(t *testing.T) {
		// check for the 401 code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
