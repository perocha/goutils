package telemetry

import (
	"context"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type XTelemetryObject interface {
	Debug(ctx context.Context, message string, fields ...XField)
	Info(ctx context.Context, message string, fields ...XField)
	Warn(ctx context.Context, message string, fields ...XField)
	Error(ctx context.Context, message string, fields ...XField)
	Dependency(ctx context.Context, dependencyType string, target string, success bool, startTime time.Time, endTime time.Time, message string, fields ...XField)
	Request(ctx context.Context, method string, url string, duration time.Duration, responseCode string, success bool, source string, message string, fields ...XField)
	GetContextInfo(ctx context.Context, key string) string
}

// XTelemetryObjectImpl will store the logger, the app insights client and the service name
type XTelemetryObjectImpl struct {
	logger      *zap.Logger
	appinsights appinsights.TelemetryClient
	xConfig     XTelemetryConfig
}

// Initialize the telemetry object
func NewXTelemetry(cc XTelemetryConfig) (*XTelemetryObjectImpl, error) {
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
	switch cc.GetLogLevel() {
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
	logger, err := zapconfig.Build(zap.AddCallerSkip(cc.GetCallerSkip()))
	if err != nil {
		return nil, err
	}

	// Initialize App Insights client, skip initialization if instrumentation key is not provided
	var appInsightsClient appinsights.TelemetryClient
	if cc.GetInstrumentationKey() != "" {
		// Initialize telemetry client
		appInsightsClient = appinsights.NewTelemetryClient(cc.GetInstrumentationKey())

		// Set the role name
		appInsightsClient.Context().Tags.Cloud().SetRole(cc.GetServiceName())
	}

	return &XTelemetryObjectImpl{
		logger:      logger,
		appinsights: appInsightsClient,
		xConfig:     cc,
	}, nil
}

// Debug will log the message using xTelemetry (no trace to App Insights)
func (t *XTelemetryObjectImpl) Debug(ctx context.Context, message string, fields ...XField) {
	t.logger.Debug(message, convertFields(fields)...)
}

// Info will log the message using xTelemetry and also send a trace to App Insights
func (t *XTelemetryObjectImpl) Info(ctx context.Context, message string, fields ...XField) {
	// Get the operation ID from the context
	operationID, ok := ctx.Value(OperationIDKeyContextKey).(string)
	if !ok {
		operationID = ""
	}

	// Create the new log trace
	telemFields := convertFields(fields)
	telemFields = append(telemFields, zap.String(string(ServiceNameKey), t.xConfig.GetServiceName()))
	if operationID != "" {
		telemFields = append(telemFields, zap.String(string(OperationIDKeyContextKey), operationID))
	}
	t.logger.Info(message, telemFields...)

	// Create the new trace in App Insights, if the client is initialized
	if t.appinsights != nil {
		// Create the new trace
		trace := appinsights.NewTraceTelemetry(constructTraceMessage(t.xConfig.GetServiceName(), operationID, message), contracts.Information)
		// Add properties to App Insights trace
		for _, field := range fields {
			trace.Properties[field.Key] = field.Value.(string)
		}
		// Set parent id, using the operationID from the context
		if operationID != "" {
			trace.Tags.Operation().SetParentId(operationID)
			trace.Properties[string(OperationIDKeyContextKey)] = operationID
		}
		// Add service name as a property
		trace.Properties[string(ServiceNameKey)] = t.xConfig.GetServiceName()

		// Send the trace to App Insights
		t.appinsights.Track(trace)
	}
}

// Warn will log the message using xTelemetry and also send a trace to App Insights
// TODO
func (t *XTelemetryObjectImpl) Warn(ctx context.Context, message string, fields ...XField) {
	t.logger.Warn(message, convertFields(fields)...)
}

// Error will log the message using xTelemetry and also send an exception to App Insights
func (t *XTelemetryObjectImpl) Error(ctx context.Context, message string, fields ...XField) {
	// Get the operation ID from the context
	operationID, ok := ctx.Value(OperationIDKeyContextKey).(string)
	if !ok {
		operationID = ""
	}

	// Create the new error trace
	telemFields := convertFields(fields)
	telemFields = append(telemFields, zap.String(string(ServiceNameKey), t.xConfig.GetServiceName()))
	if operationID != "" {
		telemFields = append(telemFields, zap.String(string(OperationIDKeyContextKey), operationID))
	}
	t.logger.Error(message, telemFields...)

	// Create the new exception in App Insights, if the client is initialized
	if t.appinsights != nil {
		// Create the new exception
		exception := appinsights.NewExceptionTelemetry(constructTraceMessage(t.xConfig.GetServiceName(), operationID, message))
		exception.SeverityLevel = contracts.Error
		// Add properties to the exception
		for _, field := range fields {
			exception.Properties[field.Key] = field.Value.(string)
		}

		// Set parent id, using the operationID from the context
		if operationID != "" {
			exception.Tags.Operation().SetParentId(operationID)
			exception.Properties[string(OperationIDKeyContextKey)] = operationID
		}
		// Add service name as a property
		exception.Properties[string(ServiceNameKey)] = t.xConfig.GetServiceName()

		t.appinsights.Track(exception)
	}
}

// Dependency will log the dependency using xTelemetry and also send a dependency to App Insights
func (t *XTelemetryObjectImpl) Dependency(ctx context.Context, dependencyType string, target string, success bool, startTime time.Time, endTime time.Time, message string, fields ...XField) {
	// Get the operation ID from the context
	operationID, ok := ctx.Value(OperationIDKeyContextKey).(string)
	if !ok {
		operationID = ""
	}

	// Create an info trace
	telemFields := convertFields(fields)
	telemFields = append(telemFields, zap.String(string(ServiceNameKey), t.xConfig.GetServiceName()))
	if operationID != "" {
		telemFields = append(telemFields, zap.String(string(OperationIDKeyContextKey), operationID))
	}
	t.logger.Info(message, telemFields...)

	// Create the new dependency in App Insights, if the client is initialized
	if t.appinsights != nil {
		// Create the new dependency
		dependency := appinsights.NewRemoteDependencyTelemetry(constructTraceMessage(t.xConfig.GetServiceName(), operationID, message), dependencyType, target, success)
		dependency.MarkTime(startTime, endTime)
		// Add properties to the dependency
		for _, field := range fields {
			dependency.Properties[field.Key] = field.Value.(string)
		}

		// Set parent id, using the operationID from the context
		if operationID != "" {
			dependency.Tags.Operation().SetParentId(operationID)
			dependency.Properties[string(OperationIDKeyContextKey)] = operationID
		}
		// Add service name as a property
		dependency.Properties[string(ServiceNameKey)] = t.xConfig.GetServiceName()

		t.appinsights.Track(dependency)
	}
}

// Request will log the request using xTelemetry and also send a request to App Insights
func (t *XTelemetryObjectImpl) Request(ctx context.Context, method string, url string, duration time.Duration, responseCode string, success bool, source string, message string, fields ...XField) {
	// Get the operation ID from the context
	operationID, ok := ctx.Value(OperationIDKeyContextKey).(string)
	if !ok {
		operationID = ""
	}

	// Create an info trace
	telemFields := convertFields(fields)
	telemFields = append(telemFields, zap.String(string(ServiceNameKey), t.xConfig.GetServiceName()))
	if operationID != "" {
		telemFields = append(telemFields, zap.String(string(OperationIDKeyContextKey), operationID))
	}
	t.logger.Info(message, telemFields...)

	// Create the new request in App Insights, if the client is initialized
	if t.appinsights != nil {
		// Create the new request
		request := appinsights.NewRequestTelemetry(method, url, duration, responseCode)
		request.Source = source
		request.Success = success
		// Add properties to the request
		for _, field := range fields {
			request.Properties[field.Key] = field.Value.(string)
		}
		// Append the operation ID to the request properties
		if operationID != "" {
			request.Tags.Operation().SetId(operationID)
			request.Properties[string(OperationIDKeyContextKey)] = operationID
		}
		// Add service name as a property
		request.Properties[string(ServiceNameKey)] = t.xConfig.GetServiceName()

		t.appinsights.Track(request)
	}
}

// Convert the telemetry property fields to zap fields
func convertFields(fields []XField) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

// Helper function to construct the trace message
func constructTraceMessage(serviceName string, operationID string, message string) string {
	if operationID == "" {
		return serviceName + "::" + message
	} else {
		return serviceName + "::" + operationID + "::" + message
	}
}

// Get context info using a key
func GetContextInfo(ctx context.Context, key string) string {
	info, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return info
}
