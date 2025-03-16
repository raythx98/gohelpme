package errorhelper

import "github.com/go-playground/validator/v10"

type ErrorResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

func NewInternalServerError(err error) *ErrorResponse {
	return &ErrorResponse{
		Message: "Something went wrong, please try again later",
		Code:    500,
		Data:    err.Error(),
	}
}

func NewValidationError(fieldErrs []validator.FieldError, err error) *ErrorResponse {
	message := "Please check your inputs and try again"
	if fieldErrs != nil && len(fieldErrs) > 0 {
		message = validationMsg(fieldErrs[0])
	}
	return &ErrorResponse{
		Message: message,
		Code:    422,
		Data:    err.Error(),
	}
}

func validationMsg(fe validator.FieldError) string {
	genericMessage := "Please check your inputs and try again"
	if fe == nil {
		return genericMessage
	}

	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return "Invalid email"
	case "alphanum":
		return fe.Field() + " should only contain letters and numbers"
	case "min":
		return fe.Field() + " should at least have " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " should at most have " + fe.Param() + " characters"
	}
	return genericMessage
}
