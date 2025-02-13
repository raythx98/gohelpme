package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/reqctx"

	"github.com/go-playground/validator/v10"
)

// ErrorHandler handles errors and returns the appropriate response.
func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if err := reqctx.GetValue(r.Context()).Error; err != nil {
			var appError *errorhelper.AppError
			if errors.As(err, &appError) {
				HandleAppError(w, appError)
				return
			}

			var authError *errorhelper.AuthError
			if errors.As(err, &authError) {
				HandleAuthError(w, authError)
				return
			}

			var invalidValidationErr *validator.InvalidValidationError
			if errors.As(err, &invalidValidationErr) {
				HandleInvalidValidationError(w, invalidValidationErr)
				return
			}

			var validationErr validator.ValidationErrors
			if errors.As(err, &validationErr) {
				HandleValidationError(w, validationErr)
				return
			}

			HandleInternalServerError(w, err)
		}
	}
}

func HandleAppError(w http.ResponseWriter, appError *errorhelper.AppError) {
	marshal, err := json.Marshal(&errorhelper.ErrorResponse{
		Message: appError.Message,
		Code:    appError.Code,
		Data:    appError.Error(),
	})
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(marshal)
}

func HandleAuthError(w http.ResponseWriter, appError *errorhelper.AuthError) {
	marshal, err := json.Marshal(&errorhelper.ErrorResponse{
		Message: "Unauthorized",
		Code:    404,
		Data:    appError.Error(),
	})
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write(marshal)
}

func HandleInvalidValidationError(w http.ResponseWriter, validationErr *validator.InvalidValidationError) {
	marshal, err := json.Marshal(errorhelper.NewValidationError(validationErr))
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	_, _ = w.Write(marshal)
}

func HandleValidationError(w http.ResponseWriter, validationErr validator.ValidationErrors) {
	marshal, err := json.Marshal(errorhelper.NewValidationError(validationErr))
	if err != nil {
		HandleInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	_, _ = w.Write(marshal)
}

func HandleInternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	marshal, err := json.Marshal(errorhelper.NewInternalServerError(err))
	if err != nil {
		_, _ = w.Write([]byte("Internal Server Error"))
	}

	_, _ = w.Write(marshal)
}
