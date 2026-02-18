package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status int
	Error  string
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GenError(err error) Response {
	return Response{
		Status: http.StatusInternalServerError,
		Error:  err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("feild %s is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("feild %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: http.StatusBadRequest,
		Error:  strings.Join(errMsgs, ","),
	}
}
