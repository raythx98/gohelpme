package validator

import "context"

type IValidator interface {
	// GetInstance returns the underlying validator instance, useful for custom validation rules.
	GetInstance() any
	StructCtx(ctx context.Context, s any) error
}
