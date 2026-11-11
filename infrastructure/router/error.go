package router

import (
	"encoding/json/v2"
	"errors"
	"net/http"

	"github.com/vitor-franz/book-log/entities"
)

// HandleError writes a JSON error response with the status matching the application
// error code. Unknown errors return a generic internal server error response to
// hide internal details. A nil error leaves the response unchanged; response
// encoding failures write a bare 500, and response write errors are ignored.
func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	response := entities.NewInternalServerError()
	status := http.StatusInternalServerError

	applicationError, ok := errors.AsType[entities.ApplicationError](err)
	if ok {
		switch applicationError.Code {
		case "BAD_REQUEST":
			status = http.StatusBadRequest
		case "UNAUTHORIZED":
			status = http.StatusUnauthorized
		case "NOT_FOUND":
			status = http.StatusNotFound
		}

		if status != http.StatusInternalServerError {
			response = applicationError
		}
	}

	body, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(body)
}
