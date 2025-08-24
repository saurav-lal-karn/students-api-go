package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/saurav-lal-karn/students-api-go/internal/storage"
	"github.com/saurav-lal-karn/students-api-go/internal/types"
	"github.com/saurav-lal-karn/students-api-go/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// We need to serialize the data so we can use it in golang
		// For that we need to serialize it in a struct

		var student types.Student

		// Decode the request body and store it to student
		err := json.NewDecoder(r.Body).Decode(&student)
		// Check if the error is the correct one
		if errors.Is(err, io.EOF) {
			// Case where the request body is empty
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		// General error handling
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		// Validate the user request
		// Use no trust policy
		// Request validation
		// Use validator
		if err := validator.New().Struct(student); err != nil {
			// Need to typecast as the Validation function expects validation errors not general error
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		id, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("User Created successfully", slog.String("userId", fmt.Sprint(id)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": id})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the id from the request
		id := r.PathValue("id")
		slog.Info("Getting the user by id", slog.String("Student Id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid id passed by user", slog.String("student_id", id))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("Error getting user: ", slog.String("student_id", fmt.Sprint(intId)))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting the user list")

		student, err := storage.GetStudents()
		if err != nil {
			slog.Error("Error getting user list: ")
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}
