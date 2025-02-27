package validator

type IValidator interface {
	// GetInstance returns the underlying validator instance, useful for custom validation rules.
	GetInstance() any
	StructCtx(ctx, s any) error
}
