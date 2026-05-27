package main

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/seminhnva/e-com/internal/store"
)

func TestGetUser(t *testing.T) {
	app := newTestApplication(t)
	mux := app.mount()
	testtoken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Helper: trả về mock store để set behaviour trong từng test case
	userMock := app.store.Users.(*store.MockUserStore)

	t.Run("should not allow unauthenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkReponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated request and return 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testtoken)

		rr := executeRequest(req, mux)

		checkReponseCode(t, http.StatusOK, rr.Code)
	})

	t.Run("should return 404 when user not found", func(t *testing.T) {
		// JWT sub=42: middleware gọi getUser(42) phải thành công.
		// Handler gọi getUser(999) → ErrNotFound → 404.
		userMock.GetByIdFn = func(_ context.Context, id int64) (*store.User, error) {
			if id == 42 {
				return &store.User{}, nil
			}
			return nil, store.ErrNotFound
		}
		defer func() { userMock.GetByIdFn = nil }()

		req, err := http.NewRequest(http.MethodGet, "/v1/users/999", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testtoken)

		rr := executeRequest(req, mux)

		checkReponseCode(t, http.StatusNotFound, rr.Code)
	})

	t.Run("should return 400 for non-numeric user ID", func(t *testing.T) {
		// strconv.ParseInt fail trước khi gọi DB → 400.
		req, err := http.NewRequest(http.MethodGet, "/v1/users/abc", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testtoken)

		rr := executeRequest(req, mux)

		checkReponseCode(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("should return 500 on internal server error", func(t *testing.T) {
		// JWT sub=42: middleware gọi getUser(42) phải thành công.
		// Handler gọi getUser(1) → lỗi DB → 500.
		userMock.GetByIdFn = func(_ context.Context, id int64) (*store.User, error) {
			if id == 42 {
				return &store.User{}, nil
			}
			return nil, errors.New("unexpected db error")
		}
		defer func() { userMock.GetByIdFn = nil }()

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testtoken)

		rr := executeRequest(req, mux)

		checkReponseCode(t, http.StatusInternalServerError, rr.Code)
	})
}
