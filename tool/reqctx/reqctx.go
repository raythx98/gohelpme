package reqctx

import (
	"context"
	"encoding/json"
)

// Key is used when writing to context using context.WithValue(parent, Key, Value).
var Key = "ReqCtx"

// Value stores commonly used metadata to log consistently across microservices.
//
// Uninitialized fields will not be logged due to json tag `json:"omitempty"`.
type Value struct {
	RequestId      string
	UserId         *int64  `json:"userId,omitempty"`
	IdempotencyKey *string `json:"idempotencyKey,omitempty"`
	Error          error   `json:"error,omitempty"`
	ErrorStack     string  `json:"errorStack,omitempty"`
}

// MarshalJSON customizes the JSON marshaling for the Value struct.
func (v *Value) MarshalJSON() ([]byte, error) {
	type Alias Value // Create an alias to avoid recursion
	return json.Marshal(&struct {
		*Alias
		Error string `json:"error,omitempty"`
	}{
		Alias: (*Alias)(v),
		Error: func() string {
			if v.Error != nil {
				return v.Error.Error()
			}
			return ""
		}(),
	})
}

// New initializes *Value with a required requestId.
//
// Value is mutable, and can be modified in-place using setter methods.
// Example:
//
//	valuePtr := reqctx.New(uuid.NewString())
//	valuePtr.SetUserId(1)
//	fmt.Println(valuePtr.UserId) // 1
func New(requestId string) *Value {
	return &Value{RequestId: requestId}
}

// SetUserId modifies UserId of Value in place.
func (v *Value) SetUserId(userId int64) *Value {
	if v == nil {
		return v
	}
	v.UserId = &userId
	return v
}

// SetIdempotencyKey modifies IdempotencyKey of Value in place.
func (v *Value) SetIdempotencyKey(idemKey string) *Value {
	if v == nil {
		return v
	}
	v.IdempotencyKey = &idemKey
	return v
}

// SetError modifies Error of Value in place.
func (v *Value) SetError(error error) *Value {
	if v == nil {
		return v
	}
	v.Error = error
	return v
}

// SetErrorStack modifies ErrorStack of Value in place.
func (v *Value) SetErrorStack(errorStack []byte) *Value {
	if v == nil {
		return v
	}
	v.ErrorStack = string(errorStack)
	return v
}

// GetValue retrieves a pointer to Value.
//
// Value is mutable, a common usage is to retrieve from context to set values.
// Example:
//
//	ctx := context.WithValue(context.Background(), reqctx.Key, reqctx.New(uuid.NewString()))
//	reqctx.GetValue(ctx).SetUserId(1)
func GetValue(ctx context.Context) *Value {
	if reqCtx, ok := ctx.Value(Key).(*Value); ok {
		return reqCtx
	}
	return nil
}
