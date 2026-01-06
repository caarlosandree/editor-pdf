package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupCORS configura o middleware CORS
func SetupCORS(allowedOrigins []string) echo.MiddlewareFunc {
	origins := allowedOrigins
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     origins,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
		MaxAge:           86400, // 24 horas
	})
}
