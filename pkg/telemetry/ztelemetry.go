package telemetry

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZTelemetry struct {
	logger *zap.Logger
}

type ZField struct {
	Key   string
	Value string
}

// NewZTelemetry creates a new ZTelemetry instance
func NewZTelemetry(LogLevel string, callerSkip int) (*ZTelemetry, error) {
	// Define configuration for the logger
	config := zap.Config{
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	// Set the log level
	switch LogLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}

	// Create a new logger, with the provided caller skip (1 will skip the current frame, since we are in the telemetry package)
	logger, err := config.Build(zap.AddCallerSkip(callerSkip))
	if err != nil {
		return nil, err
	}

	return &ZTelemetry{logger: logger}, nil
}

// Debug logs a debug message with the given tags
func (z *ZTelemetry) Debug(ctx context.Context, msg string, tags ...ZField) {
	fields := make([]zap.Field, len(tags))
	for i, tag := range tags {
		fields[i] = zap.String(tag.Key, tag.Value)
	}
	z.logger.Debug(msg, fields...)
}

// Info logs an info message with the given tags
func (z *ZTelemetry) Info(ctx context.Context, msg string, tags ...ZField) {
	fields := make([]zap.Field, len(tags))
	for i, tag := range tags {
		fields[i] = zap.String(tag.Key, tag.Value)
	}

	z.logger.Info(msg, fields...)
}

// Warn logs a warning message with the given tags
func (z *ZTelemetry) Warn(ctx context.Context, msg string, tags ...ZField) {
	fields := make([]zap.Field, len(tags))
	for i, tag := range tags {
		fields[i] = zap.String(tag.Key, tag.Value)
	}
	z.logger.Warn(msg, fields...)
}

// Error logs an error message with the given tags
func (z *ZTelemetry) Error(ctx context.Context, msg string, tags ...ZField) {
	fields := make([]zap.Field, len(tags))
	for i, tag := range tags {
		fields[i] = zap.String(tag.Key, tag.Value)
	}
	z.logger.Error(msg, fields...)
}

func (z *ZTelemetry) Sync() error {
	return z.logger.Sync()
}
