package main

import (
	"net/http"
	"testing"

	"github.com/raffyshira/project-rest-api/internal/store"
	"github.com/raffyshira/project-rest-api/internal/store/cache"
	"github.com/stretchr/testify/mock"
)

func TestLogoutAndTokenBlacklist(t *testing.T) {
	withRedis := config{
		redisCfg: redisConfig{
			enabled: true,
		},
	}

	app := newTestApplication(t, withRedis)
	mux := app.mount()

	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated logout", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/v1/authentication/logout", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should successfully logout and blacklist token", func(t *testing.T) {
		mockUserStore := app.cacheStorage.Users.(*cache.MockUserStore)
		mockUserStore.On("Get", int64(1)).Return(&store.User{ID: 1}, nil).Maybe()

		mockTokenStore := app.cacheStorage.Tokens.(*cache.MockTokenStore)
		mockTokenStore.On("Blacklist", "mock-test-jti", mock.AnythingOfType("time.Duration")).Return(nil).Once()

		req, err := http.NewRequest(http.MethodPost, "/v1/authentication/logout", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusNoContent, rr.Code)

		mockTokenStore.AssertCalled(t, "Blacklist", "mock-test-jti", mock.AnythingOfType("time.Duration"))
	})

	t.Run("should reject revoked token on protected routes", func(t *testing.T) {
		mockTokenStore := app.cacheStorage.Tokens.(*cache.MockTokenStore)
		// Simulasikan token sudah diblacklist
		mockTokenStore.ExpectedCalls = nil // reset
		mockTokenStore.On("IsBlacklisted", "mock-test-jti").Return(true, nil).Once()

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)

		// Kembalikan default mock
		mockTokenStore.ExpectedCalls = nil
		mockTokenStore.On("IsBlacklisted", mock.Anything).Return(false, nil).Maybe()
	})
}
