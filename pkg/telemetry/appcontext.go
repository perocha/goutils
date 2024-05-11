package telemetry

import (
	"context"
	"log"
)

// OperationIDKey represents the key type for the operation ID in context
type OperationIDKey string
type TelemetryObj string

const (
	// OperationIDKeyContextKey is the key used to store the operation ID in context
	OperationIDKeyContextKey OperationIDKey = "operationID"

	// TelemetryContextKey represents the key type for the telemetry object in context
	TelemetryContextKey TelemetryObj = "telemetry"

	// Service name key
	ServiceNameKey = "ServiceName"
)

// Helper function to retrieve the telemetry client from the context
func GetXTelemetryClient(ctx context.Context) *XTelemetryObjectImpl {
	telemetryClient, ok := ctx.Value(TelemetryContextKey).(*XTelemetryObjectImpl)
	if !ok {
		log.Panic("Telemetry client not found in context")
	}
	return telemetryClient
}

// Helper function to retrieve the operation ID from the context
func GetOperationID(ctx context.Context) string {
	operationID, ok := ctx.Value(OperationIDKeyContextKey).(string)
	if !ok {
		log.Panic("Operation ID not found in context")
	}
	return operationID
}

// Helper function to retrieve the service name from the context
func GetServiceName(ctx context.Context) string {
	serviceName, ok := ctx.Value(ServiceNameKey).(string)
	if !ok {
		log.Panic("Service name not found in context")
	}
	return serviceName
}
