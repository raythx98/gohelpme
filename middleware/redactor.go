package middleware

import (
	"encoding/json"
	"strings"
)

const redactedValue = "*REDACTED*"

func redact(logs map[string]interface{}, redactedFields []string) map[string]interface{} {
	if len(redactedFields) == 0 {
		return logs
	}

	logsCopy := deepCopy(logs)

	for _, path := range redactedFields {
		parts := strings.Split(path, ".")
		redactByPath(logsCopy, parts)
	}

	return logsCopy
}

func deepCopy(m map[string]interface{}) map[string]interface{} {
	b, _ := json.Marshal(m)
	var copy map[string]interface{}
	json.Unmarshal(b, &copy)
	return copy
}

func redactByPath(current interface{}, path []string) {
	if len(path) == 0 {
		return
	}

	key := path[0]

	switch v := current.(type) {
	case map[string]interface{}:
		if strings.EqualFold(key, "headers") && len(path) > 1 {
			if headersVal, ok := v["headers"]; ok {
				targetHeader := path[1]

				// Case 1: map[string][]string
				if headers, ok := headersVal.(map[string][]string); ok {
					for hKey := range headers {
						if strings.EqualFold(hKey, targetHeader) {
							for i := range headers[hKey] {
								headers[hKey][i] = redactedValue
							}
						}
					}
					return
				}

				// Case 2: map[string]interface{} (common after json unmarshal)
				if headers, ok := headersVal.(map[string]interface{}); ok {
					for hKey, hVal := range headers {
						if strings.EqualFold(hKey, targetHeader) {
							if hValSlice, ok := hVal.([]interface{}); ok {
								for i := range hValSlice {
									hValSlice[i] = redactedValue
								}
							} else {
								headers[hKey] = redactedValue
							}
						}
					}
					return
				}
			}
		}

		val, ok := v[key]
		if !ok {
			return
		}

		if len(path) == 1 {
			v[key] = redactedValue
			return
		}

		if strVal, ok := val.(string); ok {
			var nested interface{}
			if err := json.Unmarshal([]byte(strVal), &nested); err == nil {
				redactByPath(nested, path[1:])
				if redactedJSON, err := json.Marshal(nested); err == nil {
					v[key] = string(redactedJSON)
				}
			}
		} else {
			redactByPath(val, path[1:])
		}

	case []interface{}:
		for _, item := range v {
			redactByPath(item, path)
		}
	}
}
