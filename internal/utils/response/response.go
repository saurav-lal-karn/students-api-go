package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "ok"
	StatusError = "error"
)

// Data can be any type
// interface{} is an alternative to any
func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Encode data to send to response
	return json.NewEncoder(w).Encode(data)
}

// General Error handling
func GeneralError(err error) ErrorResponse {
	return ErrorResponse{
		Status: StatusError,
		Error:  err.Error(),
	}
}

// Validation error
func ValidationError(errs validator.ValidationErrors) ErrorResponse {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s is required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s is invalid", err.Field()))
		}
	}

	return ErrorResponse{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ","),
	}
}
