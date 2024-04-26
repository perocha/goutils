package telemetry

type XTelemetry interface {
	String(key string, val string) XField
	Int(key string, val int) XField
	Error(key string, val error) XField
}

type XField struct {
	Key   string
	Value interface{}
}

func XTelemetryInit() XTelemetry {
	return &XField{}
}

func (f *XField) String(key string, val string) XField {
	return XField{Key: key, Value: val}
}

func (f XField) Int(key string, val int) XField {
	return XField{Key: key, Value: val}
}

func (f XField) Error(key string, val error) XField {
	return XField{Key: key, Value: val}
}

func XString(key string, val string) XField {
	return XField{Key: key, Value: val}
}
