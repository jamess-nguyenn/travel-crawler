package helper

func ResolveTraceID(traceID string) string {
	if traceID != "" {
		return traceID
	}

	return NewUUID()
}
