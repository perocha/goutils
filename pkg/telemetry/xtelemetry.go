package telemetry

import (
	"context"
	"errors"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type XTelemetry interface {
	Debug(ctx context.Context, message string, fields ...XField)
	Info(ctx context.Context, message string, fields ...XField)
	Warn(ctx context.Context, message string, fields ...XField)
	Error(ctx context.Context, message string, fields ...XField)
}

type XField struct {
	Key   string
	Value interface{}
}

type XTelemetryConfig struct {
	instrumentationKey string
	serviceName        string
	logLevel           string
	callerSkip         int
}

type XTelemetryImpl struct {
	logger      *zap.Logger
	appinsights appinsights.TelemetryClient
}

func NewXTelemetry(customConfig XTelemetryConfig) (*XTelemetryImpl, error) {
	if customConfig.instrumentationKey == "" {
		return nil, errors.New("app insights instrumentation key not initialized")
	}

	// Initialize telemetry client
	appInsightsClient := appinsights.NewTelemetryClient(customConfig.instrumentationKey)

	// Set the role name
	appInsightsClient.Context().Tags.Cloud().SetRole(customConfig.serviceName)

	// Define configuration for the logger
	zapconfig := zap.Config{
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
	switch customConfig.logLevel {
	case "debug":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapconfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapconfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Create a new logger, with the provided caller skip (1 would skip the current frame, since we are in the telemetry package)
	logger, err := zapconfig.Build(zap.AddCallerSkip(customConfig.callerSkip))
	if err != nil {
		return nil, err
	}

	return &XTelemetryImpl{
		logger:      logger,
		appinsights: appInsightsClient,
	}, nil
}

func (t *XTelemetryImpl) Debug(ctx context.Context, message string, fields ...XField) {
	t.logger.Debug(message, convertFields(fields)...)
}

func (t *XTelemetryImpl) Info(ctx context.Context, message string, fields ...XField) {
	// Create the new trace
	t.logger.Info(message, convertFields(fields)...)

	// Create the new trace
	trace := appinsights.NewTraceTelemetry(message, contracts.Information)

	// Get the operation ID from the context
	operationID, ok := ctx.Value("OperationID").(string)
	if !ok {
		operationID = ""
	}

	// Add properties to the trace
	for _, field := range fields {
		trace.Properties[field.Key] = field.Value.(string)
	}

	// Set parent id, using the operationID from the context
	if operationID != "" {
		trace.Tags.Operation().SetParentId(operationID)
	}

	// Send the trace to App Insights
	t.appinsights.Track(trace)
}

func (t *XTelemetryImpl) Warn(ctx context.Context, message string, fields ...XField) {
	t.logger.Warn(message, convertFields(fields)...)
}

func (t *XTelemetryImpl) Error(ctx context.Context, message string, fields ...XField) {
	t.logger.Error(message, convertFields(fields)...)
}

func convertFields(fields []XField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}
