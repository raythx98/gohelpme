package middleware

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRedact(t *testing.T) {
	tests := []struct {
		name           string
		logs           map[string]interface{}
		redactedFields []string
		expected       map[string]interface{}
	}{
		{
			name: "Redact request body password",
			logs: map[string]interface{}{
				"request": map[string]interface{}{
					"body": `{"email":"test@gmail.com","password":"SENSITIVE"}`,
				},
			},
			redactedFields: []string{"request.body.password"},
			expected: map[string]interface{}{
				"request": map[string]interface{}{
					"body": `{"email":"test@gmail.com","password":"*REDACTED*"}`,
				},
			},
		},
		{
			name: "Redact headers case-insensitive",
			logs: map[string]interface{}{
				"request": map[string]interface{}{
					"headers": map[string][]string{
						"Authorization": {"Bearer secret-token"},
						"X-API-Key":     {"key123"},
					},
				},
			},
			redactedFields: []string{"request.headers.authorization"},
			expected: map[string]interface{}{
				"request": map[string]interface{}{
					"headers": map[string][]string{
						"Authorization": {"*REDACTED*"},
						"X-API-Key":     {"key123"},
					},
				},
			},
		},
		{
			name: "Redact GraphQL nested variables",
			logs: map[string]interface{}{
				"request": map[string]interface{}{
					"body": `{"query":"mutation...","variables":{"email":"raytoh@gmail.com","password":"SENSITIVE"}}`,
				},
			},
			redactedFields: []string{"request.body.variables.password"},
			expected: map[string]interface{}{
				"request": map[string]interface{}{
					"body": `{"query":"mutation...","variables":{"email":"raytoh@gmail.com","password":"*REDACTED*"}}`,
				},
			},
		},
		{
			name: "Redact nested response body",
			logs: map[string]interface{}{
				"response": map[string]interface{}{
					"body": `{"data":{"login":{"id":6,"access_token":"secret"}}}`,
				},
			},
			redactedFields: []string{"response.body.data.login.access_token"},
			expected: map[string]interface{}{
				"response": map[string]interface{}{
					"body": `{"data":{"login":{"id":6,"access_token":"*REDACTED*"}}}`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := redact(tt.logs, tt.redactedFields)
			if !equalValues(got, tt.expected) {
				gotJSON, _ := json.Marshal(got)
				expectedJSON, _ := json.Marshal(tt.expected)
				t.Errorf("redact() = %s, want %s", string(gotJSON), string(expectedJSON))
			}
		})
	}
}

func equalValues(a, b interface{}) bool {
	if reflect.DeepEqual(a, b) {
		return true
	}

	// If both are strings, they might be JSON
	if sa, ok := a.(string); ok {
		if sb, ok := b.(string); ok {
			var ma, mb interface{}
			errA := json.Unmarshal([]byte(sa), &ma)
			errB := json.Unmarshal([]byte(sb), &mb)
			if errA == nil && errB == nil {
				return equalValues(ma, mb)
			}
		}
	}

	// Recurse into maps
	if ma, ok := a.(map[string]interface{}); ok {
		if mb, ok := b.(map[string]interface{}); ok {
			if len(ma) != len(mb) {
				return false
			}
			for k, va := range ma {
				vb, ok := mb[k]
				if !ok || !equalValues(va, vb) {
					return false
				}
			}
			return true
		}
	}
	
	// Handle cases where types might differ after unmarshal (e.g. map[string][]string vs map[string]interface{})
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	var am, bm interface{}
	json.Unmarshal(aj, &am)
	json.Unmarshal(bj, &bm)
	
	// Final check on normalized interface{}
	return reflect.DeepEqual(am, bm)
}
