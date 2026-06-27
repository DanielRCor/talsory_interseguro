package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type StatisticsRequest struct {
	Matrices struct {
		Q [][]float64 `json:"q"`
		R [][]float64 `json:"r"`
	} `json:"matrices"`
}

type StatisticsResponse struct {
	Max               float64  `json:"max"`
	Min               float64  `json:"min"`
	Average           float64  `json:"average"`
	Sum               float64  `json:"sum"`
	HasDiagonalMatrix bool     `json:"hasDiagonalMatrix"`
	DiagonalMatrices  []string `json:"diagonalMatrices"`
}

type StatisticsCalculator interface {
	Calculate(ctx context.Context, q, r [][]float64) (StatisticsResponse, error)
}

type StatisticsClient struct {
	baseURL    string
	httpClient *http.Client
	jwtSecret  string
}

func NewStatisticsClient(baseURL string, httpClient *http.Client, jwtSecret string) *StatisticsClient {
	return &StatisticsClient{
		baseURL:    baseURL,
		httpClient: httpClient,
		jwtSecret:  jwtSecret,
	}
}

func (c *StatisticsClient) Calculate(ctx context.Context, q, r [][]float64) (StatisticsResponse, error) {
	var payload StatisticsRequest
	payload.Matrices.Q = q
	payload.Matrices.R = r

	body, err := json.Marshal(payload)
	if err != nil {
		return StatisticsResponse{}, fmt.Errorf("marshal statistics payload: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/statistics", bytes.NewReader(body))
	if err != nil {
		return StatisticsResponse{}, fmt.Errorf("build statistics request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")

	if c.jwtSecret != "" {
		token, err := c.issueServiceToken()
		if err != nil {
			return StatisticsResponse{}, fmt.Errorf("issue service jwt: %w", err)
		}
		request.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return StatisticsResponse{}, fmt.Errorf("call node statistics api: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return StatisticsResponse{}, fmt.Errorf("node statistics api returned status %d", response.StatusCode)
	}

	var statistics StatisticsResponse
	if err := json.NewDecoder(response.Body).Decode(&statistics); err != nil {
		return StatisticsResponse{}, fmt.Errorf("decode statistics response: %w", err)
	}

	return statistics, nil
}

func (c *StatisticsClient) issueServiceToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "go-api-service",
		"iss": "go-qr-api",
		"aud": "node-statistics-api",
		"exp": time.Now().Add(5 * time.Minute).Unix(),
	})

	return token.SignedString([]byte(c.jwtSecret))
}
