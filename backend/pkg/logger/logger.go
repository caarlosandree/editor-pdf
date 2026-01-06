package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger é o logger global da aplicação
	Logger *zap.Logger
)

// colorLevelEncoder adiciona cores aos níveis de log para desenvolvimento
func colorLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var color string
	switch level {
	case zapcore.DebugLevel:
		color = "\033[36m" // Cyan
	case zapcore.InfoLevel:
		color = "\033[32m" // Green
	case zapcore.WarnLevel:
		color = "\033[33m" // Yellow
	case zapcore.ErrorLevel:
		color = "\033[31m" // Red
	case zapcore.FatalLevel:
		color = "\033[35m" // Magenta
	default:
		color = "\033[0m" // Reset
	}
	enc.AppendString(fmt.Sprintf("%s%s\033[0m", color, level.CapitalString()))
}

// InitLogger inicializa o logger com configuração baseada no ambiente
func InitLogger(env string) error {
	var config zap.Config

	if env == "production" {
		// Produção: JSON estruturado
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.MessageKey = "message"
		config.EncoderConfig.LevelKey = "level"
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.StacktraceKey = "stacktrace"
	} else {
		// Desenvolvimento: Console colorido
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.EncoderConfig.EncodeLevel = colorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		config.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
		config.EncoderConfig.ConsoleSeparator = " | "
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	var err error
	Logger, err = config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	if err != nil {
		return err
	}

	// Substitui o logger global do zap
	zap.ReplaceGlobals(Logger)

	return nil
}

// Sync sincroniza o logger (deve ser chamado antes de encerrar a aplicação)
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}
