package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/editor-pdf/backend/pkg/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggingMiddleware cria um middleware de logging customizado
// que captura RequestID, adiciona ao context e loga requisições HTTP
func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Obtém RequestID do Echo (gerado pelo middleware.RequestID())
			// O middleware.RequestID() armazena o ID no header da resposta
			requestID := c.Response().Header().Get(echo.HeaderXRequestID)
			if requestID == "" {
				// Tenta obter do header da requisição (caso venha do cliente)
				requestID = c.Request().Header.Get(echo.HeaderXRequestID)
				if requestID == "" {
					// Se não houver RequestID, gera um novo
					requestID = generateRequestID()
				}
				// Garante que o RequestID está no header da resposta
				c.Response().Header().Set(echo.HeaderXRequestID, requestID)
			}

			// Adiciona loggerId ao context
			ctx := context.WithValue(c.Request().Context(), logger.LoggerIDKey, requestID)
			c.SetRequest(c.Request().WithContext(ctx))

			// Cria logger com loggerId
			log := logger.WithLoggerID(logger.Logger, requestID)

			// Obtém informações da requisição
			method := c.Request().Method
			path := c.Request().URL.Path
			query := c.Request().URL.RawQuery
			if query != "" {
				path += "?" + query
			}
			ip := c.RealIP()
			userAgent := c.Request().UserAgent()

			// Loga início da requisição
			log.Info("Iniciando requisição",
				zap.String("method", method),
				zap.String("path", path),
				zap.String("ip", ip),
				zap.String("user_agent", userAgent),
			)

			// Executa o próximo handler
			err := next(c)

			// Calcula duração
			duration := time.Since(start)

			// Obtém status code
			status := c.Response().Status

			// Obtém user_id se autenticado
			userID := GetUserID(c)

			// Prepara campos do log
			fields := []zap.Field{
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.Duration("duration", duration),
				zap.Int64("duration_ms", duration.Milliseconds()),
			}

			if userID != "" {
				fields = append(fields, zap.String("user_id", userID))
			}

			// Loga fim da requisição
			if err != nil {
				// Se houver erro, loga como erro
				log.Error("Requisição concluída com erro",
					append(fields, zap.Error(err))...,
				)
			} else {
				// Loga baseado no status code
				if status >= 500 {
					log.Error("Requisição concluída com erro do servidor", fields...)
				} else if status >= 400 {
					log.Warn("Requisição concluída com erro do cliente", fields...)
				} else {
					log.Info("Requisição concluída com sucesso", fields...)
				}
			}

			return err
		}
	}
}

// generateRequestID gera um ID único para a requisição
// (fallback caso o middleware.RequestID() não tenha gerado)
func generateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// Fallback para timestamp se rand falhar
		return time.Now().Format("20060102150405")
	}
	return hex.EncodeToString(b)
}
