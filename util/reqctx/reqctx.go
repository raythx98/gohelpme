package reqctx

import (
	"context"
)

var Key = "ReqCtx"

type Value struct {
	RequestId string
	UserId    *int64 `json:"Userid,omitempty"`
}

func (v *Value) SetUserId(userId int64) {
	if v == nil {
		return
	}
	v.UserId = &userId
}

func New(requestId string) *Value {
	return &Value{RequestId: requestId}
}

func GetValue(ctx context.Context) *Value {
	if reqCtx, ok := ctx.Value(Key).(*Value); ok {
		return reqCtx
	}
	return nil
}
