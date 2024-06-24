package reqctx

import (
	"context"
)

// Key is used when writing to context using context.WithValue(parent, Key, Value).
var Key = "ReqCtx"

// Value stores commonly used metadata to log consistently across microservices.
//
// Uninitialized fields will not be logged due to json tag `json:"omitempty"`.
type Value struct {
	RequestId      string
	UserId         *int64  `json:"Userid,omitempty"`
	IdempotencyKey *string `json:"idempotencyKey,omitempty"`
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
func (v *Value) SetUserId(userId int64) {
	if v == nil {
		return
	}
	v.UserId = &userId
}

// SetIdempotencyKey modifies IdempotencyKey of Value in place.
func (v *Value) SetIdempotencyKey(idemKey string) {
	if v == nil {
		return
	}
	v.IdempotencyKey = &idemKey
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
