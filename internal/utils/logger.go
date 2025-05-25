package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLogger *zap.Logger

// InitLogger inicializa el logger con la configuración apropiada
func InitLogger() error {
	// Configurar el encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Configurar el core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	// Crear el logger
	zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}

// Logger devuelve la instancia del logger
func Logger() *zap.Logger {
	return zapLogger
}

// LogRequest registra una petición HTTP
func LogRequest(method, path string, status int, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", status),
		zap.Duration("duration", duration),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		zapLogger.Error("HTTP Request", fields...)
		return
	}

	zapLogger.Info("HTTP Request", fields...)
}

// LogDB registra operaciones de base de datos
func LogDB(operation string, err error) {
	fields := []zap.Field{
		zap.String("operation", operation),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		zapLogger.Error("Database Operation", fields...)
		return
	}

	zapLogger.Info("Database Operation", fields...)
}

// LogAsset registra operaciones con assets
func LogAsset(operation string, assetName string, err error) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("asset", assetName),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		zapLogger.Error("Asset Operation", fields...)
		return
	}

	zapLogger.Info("Asset Operation", fields...)
}

// LogServer registra eventos del servidor
func LogServer(event string, err error) {
	fields := []zap.Field{
		zap.String("event", event),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		zapLogger.Error("Server Event", fields...)
		return
	}

	zapLogger.Info("Server Event", fields...)
} 