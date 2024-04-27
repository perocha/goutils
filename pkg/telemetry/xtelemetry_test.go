package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MockXTelemetryConfig struct {
	logLevel string
}

func (m *MockXTelemetryConfig) GetLogLevel() string {
	return m.logLevel
}

func (m *MockXTelemetryConfig) GetInstrumentationKey() string {
	return ""
}

func (m *MockXTelemetryConfig) GetServiceName() string {
	return ""
}

func (m *MockXTelemetryConfig) GetCallerSkip() int {
	return 0
}

func (m *MockXTelemetryConfig) SetInstrumentationKey(string) {
}

func (m *MockXTelemetryConfig) SetServiceName(string) {
}

func (m *MockXTelemetryConfig) SetLogLevel(string) {
}

func (m *MockXTelemetryConfig) SetCallerSkip(int) {
}

// captureStdout captures the output of a function that writes to stdout
func captureStdout(f func()) string {
	// Redirect stdout to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	// Restore stdout
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNewXTelemetry(t *testing.T) {
	testCases := []struct {
		name     string
		logLevel string
		want     zapcore.Level
	}{
		{"Debug level", "debug", zap.DebugLevel},
		{"Info level", "info", zap.InfoLevel},
		{"Warn level", "warn", zap.WarnLevel},
		{"Error level", "error", zap.ErrorLevel},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &MockXTelemetryConfig{logLevel: tc.logLevel}
			telemetry, err := NewXTelemetry(config)
			if err != nil {
				t.Fatalf("NewXTelemetry() error = %v", err)
			}

			if got := telemetry.logger.Core().Enabled(tc.want); !got {
				t.Errorf("NewXTelemetry() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestTelemetryOutput(t *testing.T) {
	testCases := []struct {
		level       string
		message     string
		expectedMsg string
	}{
		{"debug", "test debug message", "test debug message"},
		{"info", "test info message", "test info message"},
		{"warn", "test warn message", "test warn message"},
		{"error", "test error message", "test error message"},
	}

	for _, tc := range testCases {
		// Capture stdout output
		output := captureStdout(func() {
			// Create telemetry object
			config := &MockXTelemetryConfig{logLevel: tc.level}
			telemetry, err := NewXTelemetry(config)
			if err != nil {
				t.Fatalf("NewXTelemetry() error = %v", err)
			}

			// Call the appropriate method based on test case
			switch tc.level {
			case "debug":
				telemetry.Debug(context.Background(), tc.message)
			case "info":
				telemetry.Info(context.Background(), tc.message)
			case "warn":
				telemetry.Warn(context.Background(), tc.message)
			case "error":
				telemetry.Error(context.Background(), tc.message)
			}
		})

		// Parse the JSON output
		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(output), &logEntry)
		if err != nil {
			t.Fatalf("Failed to parse JSON: %v", err)
		}

		// Assert against specific fields or values
		expectedLevel := tc.level
		if level, ok := logEntry["level"].(string); !ok || level != expectedLevel {
			t.Errorf("Expected level %q, got %q", expectedLevel, level)
		}

		expectedMsg := tc.expectedMsg
		if msg, ok := logEntry["msg"].(string); !ok || msg != expectedMsg {
			t.Errorf("Expected message %q, got %q", expectedMsg, msg)
		}
	}
}
