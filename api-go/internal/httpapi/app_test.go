package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DanielRCor/talsory_interseguro/api-go/internal/client"
	"github.com/DanielRCor/talsory_interseguro/api-go/internal/config"
)

func TestAnalyzeHandler(t *testing.T) {
	cfg := config.Load()
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
	app := NewApp(AppDependencies{
		Config:           config.Load(),
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
	app := NewApp(AppDependencies{
		Config:           config.Load(),
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
