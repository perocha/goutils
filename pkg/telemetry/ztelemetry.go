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

func NewZTelemetry(callerSkip int64) (*ZTelemetry, error) {
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
	operationID, _ := ctx.Value(OperationIDKeyContextKey).(string)
	fields = append(fields, zap.String(string(OperationIDKeyContextKey), operationID))
	z.logger.Debug(msg, fields...)
}

// Info logs an info message with the given tags
func (z *ZTelemetry) Info(ctx context.Context, msg string, tags ...ZField) {
	fields := make([]zap.Field, len(tags))
	for i, tag := range tags {
		fields[i] = zap.String(tag.Key, tag.Value)
	}
	operationID, _ := ctx.Value(OperationIDKeyContextKey).(string)
	fields = append(fields, zap.String(string(OperationIDKeyContextKey), operationID))
	z.logger.Info(msg, fields...)
}

// Warn logs a warning message with the given tags
func (z *ZTelemetry) Warn(ctx context.Context, msg string, tags ...ZField) {
	fields := make([]zap.Field, len(tags))
	for i, tag := range tags {
		fields[i] = zap.String(tag.Key, tag.Value)
	}
	operationID, _ := ctx.Value(OperationIDKeyContextKey).(string)
	fields = append(fields, zap.String(string(OperationIDKeyContextKey), operationID))
	z.logger.Warn(msg, fields...)
}

// Error logs an error message with the given tags
func (z *ZTelemetry) Error(ctx context.Context, msg string, tags ...ZField) {
	fields := make([]zap.Field, len(tags))
	for i, tag := range tags {
		fields[i] = zap.String(tag.Key, tag.Value)
	}
	operationID, _ := ctx.Value(OperationIDKeyContextKey).(string)
	fields = append(fields, zap.String(string(OperationIDKeyContextKey), operationID))
	z.logger.Error(msg, fields...)
}

func (z *ZTelemetry) Sync() error {
	return z.logger.Sync()
}
