package logger

import (
	"context"
	"go.uber.org/zap"
)

const (
	// LoggerIDKey é a chave usada para armazenar o loggerId no context
	LoggerIDKey = "loggerId"
)

// FromContext retorna um logger com loggerId extraído do context
// Se o loggerId não existir no context, retorna o logger global sem loggerId
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Logger
	}

	loggerID, ok := ctx.Value(LoggerIDKey).(string)
	if !ok || loggerID == "" {
		return Logger
	}

	return Logger.With(zap.String("loggerId", loggerID))
}

// WithLoggerID adiciona um loggerId ao logger
func WithLoggerID(logger *zap.Logger, loggerID string) *zap.Logger {
	if loggerID == "" {
		return logger
	}
	return logger.With(zap.String("loggerId", loggerID))
}

// WithContext adiciona o loggerId do context ao logger
func WithContext(logger *zap.Logger, ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}

	loggerID, ok := ctx.Value(LoggerIDKey).(string)
	if !ok || loggerID == "" {
		return logger
	}

	return logger.With(zap.String("loggerId", loggerID))
}

// WithFields adiciona campos customizados ao logger
func WithFields(logger *zap.Logger, fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

// WithUserID adiciona user_id ao logger
func WithUserID(logger *zap.Logger, userID string) *zap.Logger {
	if userID == "" {
		return logger
	}
	return logger.With(zap.String("user_id", userID))
}

// WithDocumentID adiciona document_id ao logger
func WithDocumentID(logger *zap.Logger, documentID string) *zap.Logger {
	if documentID == "" {
		return logger
	}
	return logger.With(zap.String("document_id", documentID))
}
