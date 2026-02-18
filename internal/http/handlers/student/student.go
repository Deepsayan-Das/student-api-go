package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/Deepsayan-Das/student-api-go/internal/types"
	"github.com/Deepsayan-Das/student-api-go/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	slog.Info("Creating a new Student")
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) { //empty body provided
			err := response.WriteJson(w, http.StatusBadRequest, response.GenError(err))
			if err != nil {
				slog.Error("failed to write response", slog.String("error", err.Error()))
			}
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GenError(err))
			return
		}
		//validating request  -> go-playground/validator

		if err := validator.New().Struct(student); err != nil {
			validationErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validationErr))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"Success": "OK"})
		w.Write([]byte("Welcome to students api"))
	}
}
