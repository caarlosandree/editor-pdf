package middleware

import (
	"strings"

	"github.com/editor-pdf/backend/internal/usecase"
	"github.com/editor-pdf/backend/pkg/response"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware cria um middleware de autenticação JWT
func AuthMiddleware(authUseCase *usecase.AuthUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extrai o token do header Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.ErrorUnauthorized(c, nil, "token de autenticação não fornecido")
			}

			// Remove o prefixo "Bearer "
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return response.ErrorUnauthorized(c, nil, "formato de token inválido (esperado: Bearer <token>)")
			}

			// Valida o token
			user, err := authUseCase.ValidateToken(c.Request().Context(), tokenString)
			if err != nil {
				return response.ErrorUnauthorized(c, err, "token inválido ou expirado")
			}

			// Armazena o usuário no contexto
			c.Set("user", user)
			c.Set("user_id", user.ID)

			return next(c)
		}
	}
}

// GetUserID extrai o ID do usuário do contexto
func GetUserID(c echo.Context) string {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return ""
	}
	return userID
}

// GetUser extrai os dados do usuário do contexto
func GetUser(c echo.Context) interface{} {
	return c.Get("user")
}
