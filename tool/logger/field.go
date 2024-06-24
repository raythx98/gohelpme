package logger

// Field is a function type that returns a map of string keys to interface{} values.
// It is used to provide additional metadata for logging.
// The function is expected to return a map where each key-value pair represents a field to be logged.
type Field func() map[string]interface{}

// EmptyField returns an empty map.
func EmptyField() map[string]interface{} {
	return make(map[string]interface{})
}

// GetMapFromFields returns a map from the variadic Field.
//
// It merges all the maps returned by the Field functions into a single map.
func GetMapFromFields(fields ...Field) map[string]interface{} {
	m := make(map[string]interface{})
	for _, o := range fields {
		for k, v := range o() {
			m[k] = v
		}
	}
	return m
}

// WithError creates an option that contains the error.
func WithError(err error) Field {
	return func() map[string]interface{} {
		return map[string]interface{}{
			"error": err,
		}
	}
}

// WithField creates an option that contains the key-value pair.
func WithField(key string, value ...interface{}) Field {
	return func() map[string]interface{} {
		return map[string]interface{}{
			key: value,
		}
	}
}

// WithFields creates an option that contains the metadata.
func WithFields(metadata map[string]interface{}) Field {
	return func() map[string]interface{} {
		return metadata
	}
}
