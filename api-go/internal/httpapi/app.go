package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/DanielRCor/talsory_interseguro/api-go/internal/client"
	"github.com/DanielRCor/talsory_interseguro/api-go/internal/config"
	"github.com/DanielRCor/talsory_interseguro/api-go/internal/qr"
	"github.com/DanielRCor/talsory_interseguro/api-go/internal/rotation"
	"github.com/DanielRCor/talsory_interseguro/api-go/internal/validation"
	"github.com/gofiber/fiber/v2"
	fibercors "github.com/gofiber/fiber/v2/middleware/cors"
)

type analyzeRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

type rotateRequest struct {
	Matrix    [][]float64 `json:"matrix"`
	Direction string      `json:"direction"`
}

type errorResponse struct {
	Error struct {
		Code    string   `json:"code"`
		Message string   `json:"message"`
		Details []string `json:"details,omitempty"`
	} `json:"error"`
}

type AppDependencies struct {
	Config           config.Config
	StatisticsClient client.StatisticsCalculator
}

func NewApp(dependencies AppDependencies) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return writeError(ctx, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Unexpected server error", nil)
		},
	})

	app.Use(fibercors.New())

	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":  "ok",
			"service": "go-qr-api",
		})
	})

	app.Post("/auth/demo-token", func(ctx *fiber.Ctx) error {
		token, expiresAt, err := issueDemoToken(dependencies.Config)
		if err != nil {
			return writeError(ctx, fiber.StatusInternalServerError, "TOKEN_ISSUE_ERROR", "Failed to issue demo token", []string{err.Error()})
		}

		return ctx.JSON(fiber.Map{
			"token":            token,
			"tokenPreview":     previewToken(token),
			"expiresAt":        expiresAt.Format(time.RFC3339),
			"expiresInSeconds": int(demoTokenTTL.Seconds()),
			"enabled":          dependencies.Config.EnableAuth,
		})
	})

	api := app.Group("/api/v1")
	if dependencies.Config.EnableAuth {
		api.Use(authMiddleware(dependencies.Config))
	}

	api.Post("/qr/analyze", func(ctx *fiber.Ctx) error {
		var request analyzeRequest
		if err := ctx.BodyParser(&request); err != nil {
			return writeError(ctx, fiber.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload", []string{err.Error()})
		}

		if err := validation.ValidateQRMatrix(request.Matrix); err != nil {
			validationError, ok := err.(validation.ValidationError)
			if ok {
				return writeError(ctx, fiber.StatusUnprocessableEntity, "VALIDATION_ERROR", "Matrix validation failed", validationError.Details)
			}
			return writeError(ctx, fiber.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error(), nil)
		}

		qMatrix, rMatrix, err := qr.Factorize(request.Matrix)
		if err != nil {
			return writeError(ctx, fiber.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error(), nil)
		}

		statistics, err := dependencies.StatisticsClient.Calculate(context.Background(), qMatrix, rMatrix)
		if err != nil {
			return writeError(ctx, http.StatusBadGateway, "BAD_GATEWAY", "Failed to communicate with statistics service", []string{err.Error()})
		}

		return ctx.JSON(fiber.Map{
			"input": fiber.Map{
				"rows":    len(request.Matrix),
				"columns": len(request.Matrix[0]),
			},
			"qr": fiber.Map{
				"q": qMatrix,
				"r": rMatrix,
			},
			"statistics": statistics,
			"metadata": fiber.Map{
				"algorithm": dependencies.Config.Algorithm,
				"tolerance": dependencies.Config.Tolerance,
			},
		})
	})

	api.Post("/matrix/rotate", func(ctx *fiber.Ctx) error {
		var request rotateRequest
		if err := ctx.BodyParser(&request); err != nil {
			return writeError(ctx, fiber.StatusBadRequest, "BAD_REQUEST", "Invalid JSON payload", []string{err.Error()})
		}

		if err := validation.ValidateRectangularMatrix(request.Matrix); err != nil {
			validationError, ok := err.(validation.ValidationError)
			if ok {
				return writeError(ctx, fiber.StatusUnprocessableEntity, "VALIDATION_ERROR", "Matrix validation failed", validationError.Details)
			}
			return writeError(ctx, fiber.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error(), nil)
		}

		rotated, direction, err := rotation.Rotate(request.Matrix, request.Direction)
		if err != nil {
			return writeError(ctx, fiber.StatusUnprocessableEntity, "VALIDATION_ERROR", err.Error(), nil)
		}

		return ctx.JSON(fiber.Map{
			"input": fiber.Map{
				"rows":    len(request.Matrix),
				"columns": len(request.Matrix[0]),
			},
			"operation": direction,
			"result":    rotated,
		})
	})

	return app
}

func writeError(ctx *fiber.Ctx, statusCode int, code, message string, details []string) error {
	var payload errorResponse
	payload.Error.Code = code
	payload.Error.Message = message
	payload.Error.Details = details
	ctx.Status(statusCode)
	return ctx.JSON(payload)
}

type stubStatisticsClient struct {
	response client.StatisticsResponse
	err      error
}

func (s stubStatisticsClient) Calculate(_ context.Context, _ [][]float64, _ [][]float64) (client.StatisticsResponse, error) {
	if s.err != nil {
		return client.StatisticsResponse{}, s.err
	}
	return s.response, nil
}

func NewStubStatisticsClient(response client.StatisticsResponse, err error) client.StatisticsCalculator {
	return stubStatisticsClient{response: response, err: err}
}

func previewToken(token string) string {
	if len(token) <= 24 {
		return token
	}

	return token[:12] + "..." + token[len(token)-12:]
}
