package telemetry

type XField struct {
	Key   string
	Value interface{}
}

func String(key string, val string) XField {
	return XField{Key: key, Value: val}
}

func Int(key string, val int) XField {
	return XField{Key: key, Value: val}
}

func Bool(key string, val bool) XField {
	return XField{Key: key, Value: val}
}
