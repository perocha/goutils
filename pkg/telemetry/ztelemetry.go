package telemetry

import (
	"context"

	"go.uber.org/zap"
)

type ZTelemetry struct {
	logger *zap.Logger
}

type ZField = zap.Field

func NewZTelemetry() (*ZTelemetry, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZTelemetry{logger: logger}, nil
}

func (z *ZTelemetry) Info(ctx context.Context, msg string, tags ...ZField) {
	operationID, _ := ctx.Value(OperationIDKeyContextKey).(string)
	tags = append(tags, zap.String(string(OperationIDKeyContextKey), operationID))
	z.logger.Info(msg, tags...)
}

func (z *ZTelemetry) Error(ctx context.Context, msg string, tags ...ZField) {
	operationID, _ := ctx.Value(OperationIDKeyContextKey).(string)
	tags = append(tags, zap.String(string(OperationIDKeyContextKey), operationID))
	z.logger.Error(msg, tags...)
}

func (z *ZTelemetry) Sync() error {
	return z.logger.Sync()
}
