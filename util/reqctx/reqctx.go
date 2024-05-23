package reqctx

import (
	"context"
)

var Key = "ReqCtx"

type Value struct {
	RequestId string
	UserId    *int64 `json:"Userid,omitempty"`
}

func New(requestId string, userId *int64) *Value {
	return &Value{
		RequestId: requestId,
		UserId:    userId,
	}
}

func GetValue(ctx context.Context) *Value {
	if reqCtx, ok := ctx.Value(Key).(*Value); ok {
		return reqCtx
	}
	return nil
}
