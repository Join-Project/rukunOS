package middleware

import (
	"os"
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
	}
	return echojwt.WithConfig(config)
}

