package middleware

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRedact(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		paths    []string
		expected map[string]interface{}
	}{
		{
			name: "Simple redaction",
			data: map[string]interface{}{
				"password": "secret123",
				"user":     "ray",
			},
			paths: []string{"password"},
			expected: map[string]interface{}{
				"password": "*REDACTED*",
				"user":     "ray",
			},
		},
		{
			name: "Nested redaction",
			data: map[string]interface{}{
				"variables": map[string]interface{}{
					"password": "secret123",
					"email":    "ray@gmail.com",
				},
			},
			paths: []string{"variables.password"},
			expected: map[string]interface{}{
				"variables": map[string]interface{}{
					"password": "*REDACTED*",
					"email":    "ray@gmail.com",
				},
			},
		},
		{
			name: "Redact JSON string",
			data: map[string]interface{}{
				"body": `{"query":"mutation","variables":{"password":"secret123","email":"ray@gmail.com"}}`,
			},
			paths: []string{"body.variables.password"},
			expected: map[string]interface{}{
				"body": `{"query":"mutation","variables":{"email":"ray@gmail.com","password":"*REDACTED*"}}`,
			},
		},
		{
			name: "Multiple paths and nested JSON strings",
			data: map[string]interface{}{
				"request": map[string]interface{}{
					"body": `{"variables":{"password":"p1","token":"t1"}}`,
				},
				"response": map[string]interface{}{
					"body": `{"data":{"login":{"token":"t2"}}}`,
				},
			},
			paths: []string{"request.body.variables.password", "response.body.data.login.token"},
			expected: map[string]interface{}{
				"request": map[string]interface{}{
					"body": `{"variables":{"password":"*REDACTED*","token":"t1"}}`,
				},
				"response": map[string]interface{}{
					"body": `{"data":{"login":{"token":"*REDACTED*"}}}`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			redact(tt.data, tt.paths)
			if tt.name == "Redact JSON string" || tt.name == "Multiple paths and nested JSON strings" {
				// For JSON strings, the order of keys might change after marshal/unmarshal.
				// So we compare by unmarshaling back.
				compareMapsWithJSONStrings(t, tt.expected, tt.data)
			} else {
				if !reflect.DeepEqual(tt.data, tt.expected) {
					t.Errorf("expected %v, got %v", tt.expected, tt.data)
				}
			}
		})
	}
}

func compareMapsWithJSONStrings(t *testing.T, expected, actual map[string]interface{}) {
	for k, vExp := range expected {
		vAct, ok := actual[k]
		if !ok {
			t.Errorf("key %s missing in actual", k)
			continue
		}

		sExp, okExp := vExp.(string)
		sAct, okAct := vAct.(string)

		if okExp && okAct {
			var mExp, mAct interface{}
			errExp := json.Unmarshal([]byte(sExp), &mExp)
			errAct := json.Unmarshal([]byte(sAct), &mAct)

			if errExp == nil && errAct == nil {
				if !reflect.DeepEqual(mExp, mAct) {
					t.Errorf("key %s: expected JSON %s, got %s", k, sExp, sAct)
				}
				continue
			}
		}

		mExp, okExpM := vExp.(map[string]interface{})
		mAct, okActM := vAct.(map[string]interface{})
		if okExpM && okActM {
			compareMapsWithJSONStrings(t, mExp, mAct)
			continue
		}

		if !reflect.DeepEqual(vExp, vAct) {
			t.Errorf("key %s: expected %v, got %v", k, vExp, vAct)
		}
	}
}
