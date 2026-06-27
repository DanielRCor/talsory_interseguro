package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DanielRCor/talsory_interseguro/api-go/internal/client"
	"github.com/DanielRCor/talsory_interseguro/api-go/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func TestAnalyzeHandler(t *testing.T) {
	cfg := config.Load()
	cfg.EnableAuth = false
	app := NewApp(AppDependencies{
		Config: cfg,
		StatisticsClient: NewStubStatisticsClient(client.StatisticsResponse{
			Max:               4,
			Min:               0,
			Average:           1.375,
			Sum:               11,
			HasDiagonalMatrix: true,
			DiagonalMatrices:  []string{"q"},
		}, nil),
	})

	request := httptest.NewRequest("POST", "/api/v1/qr/analyze", strings.NewReader(`{"matrix":[[1,2],[3,4],[5,6]]}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["statistics"] == nil {
		t.Fatalf("expected statistics in response")
	}
}

func TestAnalyzeHandlerValidationError(t *testing.T) {
	cfg := config.Load()
	cfg.EnableAuth = false
	app := NewApp(AppDependencies{
		Config:           cfg,
		StatisticsClient: NewStubStatisticsClient(client.StatisticsResponse{}, nil),
	})

	request := httptest.NewRequest("POST", "/api/v1/qr/analyze", strings.NewReader(`{"matrix":[[1,2,3],[4,5,6]]}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.StatusCode != 422 {
		t.Fatalf("expected status 422, got %d", response.StatusCode)
	}
}

func TestAnalyzeHandlerBadGateway(t *testing.T) {
	cfg := config.Load()
	cfg.EnableAuth = false
	app := NewApp(AppDependencies{
		Config:           cfg,
		StatisticsClient: NewStubStatisticsClient(client.StatisticsResponse{}, fmt.Errorf("node unavailable")),
	})

	request := httptest.NewRequest("POST", "/api/v1/qr/analyze", strings.NewReader(`{"matrix":[[1,2],[3,4],[5,6]]}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.StatusCode != 502 {
		t.Fatalf("expected status 502, got %d", response.StatusCode)
	}
}

func TestRotateHandler(t *testing.T) {
	cfg := config.Load()
	cfg.EnableAuth = false
	app := NewApp(AppDependencies{
		Config:           cfg,
		StatisticsClient: NewStubStatisticsClient(client.StatisticsResponse{}, nil),
	})

	request := httptest.NewRequest("POST", "/api/v1/matrix/rotate", strings.NewReader(`{"matrix":[[1,2,3],[4,5,6]],"direction":"counterclockwise"}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload["operation"] != "counterclockwise" {
		t.Fatalf("expected counterclockwise operation, got %v", payload["operation"])
	}
}

func TestAnalyzeHandlerRequiresJWTWhenEnabled(t *testing.T) {
	cfg := config.Load()
	cfg.EnableAuth = true
	cfg.JWTSecret = "test-secret"

	app := NewApp(AppDependencies{
		Config: cfg,
		StatisticsClient: NewStubStatisticsClient(client.StatisticsResponse{
			Max: 1,
		}, nil),
	})

	request := httptest.NewRequest("POST", "/api/v1/qr/analyze", strings.NewReader(`{"matrix":[[1,2],[3,4],[5,6]]}`))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.StatusCode != 401 {
		t.Fatalf("expected status 401, got %d", response.StatusCode)
	}
}

func TestAnalyzeHandlerAcceptsJWTWhenEnabled(t *testing.T) {
	cfg := config.Load()
	cfg.EnableAuth = true
	cfg.JWTSecret = "test-secret"

	app := NewApp(AppDependencies{
		Config: cfg,
		StatisticsClient: NewStubStatisticsClient(client.StatisticsResponse{
			Max:               4,
			Min:               0,
			Average:           1.375,
			Sum:               11,
			HasDiagonalMatrix: true,
			DiagonalMatrices:  []string{"q"},
		}, nil),
	})

	token := createTestJWT(t, cfg.JWTSecret)
	request := httptest.NewRequest("POST", "/api/v1/qr/analyze", strings.NewReader(`{"matrix":[[1,2],[3,4],[5,6]]}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}
}

func TestDemoTokenEndpoint(t *testing.T) {
	cfg := config.Load()
	cfg.EnableAuth = true
	cfg.JWTSecret = "test-secret"

	app := NewApp(AppDependencies{
		Config:           cfg,
		StatisticsClient: NewStubStatisticsClient(client.StatisticsResponse{}, nil),
	})

	request := httptest.NewRequest("POST", "/auth/demo-token", nil)
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}

	var payload struct {
		Token            string `json:"token"`
		TokenPreview     string `json:"tokenPreview"`
		ExpiresInSeconds int    `json:"expiresInSeconds"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Token == "" {
		t.Fatalf("expected token in response")
	}

	if payload.TokenPreview == "" {
		t.Fatalf("expected token preview in response")
	}

	if payload.ExpiresInSeconds != 120 {
		t.Fatalf("expected 120 seconds expiry, got %d", payload.ExpiresInSeconds)
	}

	parsedToken, err := jwt.Parse(payload.Token, func(token *jwt.Token) (any, error) {
		return []byte(cfg.JWTSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected map claims")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("expected exp claim")
	}

	expiresAt := time.Unix(int64(expFloat), 0)
	if expiresAt.Before(time.Now().Add(90*time.Second)) || expiresAt.After(time.Now().Add(3*time.Minute)) {
		t.Fatalf("expected token to expire near 2 minutes, got %v", expiresAt)
	}
}

func createTestJWT(t *testing.T, secret string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "integration-test",
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}
