package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func Ok() Response {
	return Response{
		Status: StatusOk,
	}
}

func Error(message string) Response {
	return Response{
		Status: StatusError,
		Error:  message,
	}
}

func ValidateErrors(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		errField := err.Field()

		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", errField))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid URL", errField))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", errField))
		}
	}

	return Error(strings.Join(errMsgs, ", "))
}
