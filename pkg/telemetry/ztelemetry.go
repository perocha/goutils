package telemetry

import (
	"context"

	"go.uber.org/zap"
)

type ZTelemetry struct {
	logger *zap.Logger
}

type ZField struct {
	Key   string
	Value string
}

// NewZTelemetry creates a new ZTelemetry instance
func NewZTelemetry(LogLevel string, callerSkip int64) (*ZTelemetry, error) {
	// Define configuration for the logger
	config := zap.NewProductionConfig()

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
	logger, err := zap.NewProduction(zap.AddCallerSkip(int(callerSkip)))
	if err != nil {
		return nil, err
	}
	return &ZTelemetry{logger: logger}, nil
}

func String(key string, val string) ZField {
	return ZField{Key: key, Value: val}
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
