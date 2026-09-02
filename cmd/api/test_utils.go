package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raffyshira/project-rest-api/internal/auth"
	"github.com/raffyshira/project-rest-api/internal/ratelimiter"
	service "github.com/raffyshira/project-rest-api/internal/services"
	"github.com/raffyshira/project-rest-api/internal/store"
	"github.com/raffyshira/project-rest-api/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	// uncomment to enable logs
	// logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	testAuth := &auth.TestAuthenticator{}

	// rate limiter
	rateLimiter := ratelimiter.NewFixedWindowRateLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	services := service.NewServices(service.ServiceDependencies{
		Store: mockStore,
	})

	return &application{
		logger:        logger,
		store:         mockStore,
		services:      services,
		cacheStorage:  mockCacheStore,
		authenticator: testAuth,
		config:        cfg,
		rateLimiter:   rateLimiter,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
