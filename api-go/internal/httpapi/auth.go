package httpapi

import (
	"strings"
	"time"

	"github.com/DanielRCor/talsory_interseguro/api-go/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const demoTokenTTL = 2 * time.Minute

func authMiddleware(cfg config.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		header := ctx.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return writeError(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid bearer token", nil)
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(cfg.JWTSecret), nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil {
			return writeError(ctx, fiber.StatusUnauthorized, "UNAUTHORIZED", "Invalid bearer token", []string{err.Error()})
		}

		return ctx.Next()
	}
}

func issueDemoToken(cfg config.Config) (string, time.Time, error) {
	expiresAt := time.Now().Add(demoTokenTTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "demo-user",
		"iss": "go-qr-api",
		"aud": "interseguro-challenge",
		"exp": expiresAt.Unix(),
	})

	signed, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}
