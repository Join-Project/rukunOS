package middleware

import (
	"net/http"
	"os"
	"strings"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "secret" // Default for development only
	}
	
	config := echojwt.Config{
		SigningKey: []byte(jwtSecret),
		ContextKey: "user",
		ErrorHandler: func(c echo.Context, err error) error {
			// Check if token is missing
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing or malformed jwt",
					"message": "Authorization header is missing",
				})
			}
			
			// Check if Bearer prefix is present
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing or malformed jwt",
					"message": "Authorization header must start with 'Bearer '",
				})
			}
			
			// Token exists but is invalid
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "missing or malformed jwt",
				"message": err.Error(),
			})
		},
	}
	return echojwt.WithConfig(config)
}

