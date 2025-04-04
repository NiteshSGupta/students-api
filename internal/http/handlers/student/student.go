package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/NiteshSGupta/students-api/internal/storage"
	"github.com/NiteshSGupta/students-api/internal/types"
	"github.com/NiteshSGupta/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")

		//we can't get the request data directly in golang , we have to decode data in struct
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		// Checks for io.EOF (End Of File) error, which occurs when the request body is empty
		// Returns a 400 Bad Request status with an error message if the body is empty
		if errors.Is(err, io.EOF) {

			//for sedning response we sended the response in json common file , to send response in json
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))  //eof
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		//if there no empty body error , so we have to catch error again here
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation, for 0 trust policy on client side.
		//https://github.com/go-playground/validator
		validationError := validator.New().Struct(student)

		if validationError != nil {

			//typecasting the validation error in validator.ValidationErrors
			validateErrs := validationError.(validator.ValidationErrors)

			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("User created successfully", slog.String("user_id", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
		// response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//we can fetch id from request using musk
		//r.PathValue("id") same name as mension in route
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error("errr getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)

	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("getting all student")

		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusOK, students)

	}
}
