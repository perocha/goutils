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

type xConfig struct {
	instrumentationKey string
	serviceName        string
	logLevel           string
	callerSkip         int
	otlpEndPoint       string
}

func NewXTelemetryConfig(instrumentationKey string, serviceName string, logLevel string, callerSkip int, otlpEndPoint string) XTelemetryConfig {
	return &xConfig{
		instrumentationKey: instrumentationKey,
		serviceName:        serviceName,
		logLevel:           logLevel,
		callerSkip:         callerSkip,
		otlpEndPoint:       otlpEndPoint,
	}
}

func (c *xConfig) GetInstrumentationKey() string {
	return c.instrumentationKey
}

func (c *xConfig) GetServiceName() string {
	return c.serviceName
}

func (c *xConfig) GetLogLevel() string {
	return c.logLevel
}

func (c *xConfig) GetCallerSkip() int {
	return c.callerSkip
}

func (c *xConfig) SetInstrumentationKey(instrumentationKey string) {
	c.instrumentationKey = instrumentationKey
}

func (c *xConfig) SetServiceName(serviceName string) {
	c.serviceName = serviceName
}

func (c *xConfig) SetLogLevel(logLevel string) {
	c.logLevel = logLevel
}

func (c *xConfig) SetCallerSkip(callerSkip int) {
	c.callerSkip = callerSkip
}
