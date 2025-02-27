package validator

import (
	"github.com/go-playground/validator/v10"
)

type GoValidator struct {
	validate *validator.Validate
}

func New() *GoValidator {
	return &GoValidator{
		validate: validator.New(),
	}
}

// GetInstance returns the underlying instance, cast to *validator.Validate
func (v *GoValidator) GetInstance() any {
	return v.validate
}

func (v *GoValidator) StructCtx(ctx, s any) error {
	return v.StructCtx(ctx, s)
}
