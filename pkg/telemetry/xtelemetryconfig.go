package telemetry

type XTelemetryConfig interface {
	GetInstrumentationKey() string
	GetServiceName() string
	GetLogLevel() string
	GetCallerSkip() int
	SetInstrumentationKey(string)
	SetServiceName(string)
	SetLogLevel(string)
	SetCallerSkip(int)
}

type xTelemetryConfig struct {
	instrumentationKey string
	serviceName        string
	logLevel           string
	callerSkip         int
}

func NewXTelemetryConfig(instrumentationKey string, serviceName string, logLevel string, callerSkip int) XTelemetryConfig {
	return &xTelemetryConfig{
		instrumentationKey: instrumentationKey,
		serviceName:        serviceName,
		logLevel:           logLevel,
		callerSkip:         callerSkip,
	}
}

func (c *xTelemetryConfig) GetInstrumentationKey() string {
	return c.instrumentationKey
}

func (c *xTelemetryConfig) GetServiceName() string {
	return c.serviceName
}

func (c *xTelemetryConfig) GetLogLevel() string {
	return c.logLevel
}

func (c *xTelemetryConfig) GetCallerSkip() int {
	return c.callerSkip
}

func (c *xTelemetryConfig) SetInstrumentationKey(instrumentationKey string) {
	c.instrumentationKey = instrumentationKey
}

func (c *xTelemetryConfig) SetServiceName(serviceName string) {
	c.serviceName = serviceName
}

func (c *xTelemetryConfig) SetLogLevel(logLevel string) {
	c.logLevel = logLevel
}

func (c *xTelemetryConfig) SetCallerSkip(callerSkip int) {
	c.callerSkip = callerSkip
}
