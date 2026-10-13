package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func New(
	w http.ResponseWriter,
	status string,
	code int,
	message string,
	data any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if data == nil {
		data = map[string]any{}
	}

	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(
		&Response{
			Status:  status,
			Code:    code,
			Message: message,
			Data:    data,
		},
	)
}

func Success(
	w http.ResponseWriter,
	data any,
	message string,
	code int,
) {
	New(
		w,
		"success",
		code,
		message,
		data,
	)
}

func Error(
	w http.ResponseWriter,
	message string,
	data any,
	code int,
) {
	New(
		w,
		"error",
		code,
		message,
		data,
	)
}

func SuccessOK(
	w http.ResponseWriter,
	data any,
	message string,
) {
	Success(
		w,
		data,
		message,
		http.StatusOK,
	)
}

func SuccessCreated(
	w http.ResponseWriter,
	data any,
	message string,
) {
	Success(
		w,
		data,
		message,
		http.StatusCreated,
	)
}

func ErrorNoContent(
	w http.ResponseWriter,
	data any,
) {
	status := http.StatusNoContent

	Error(
		w,
		http.StatusText(status),
		data,
		status,
	)
}

func ErrorBadRequest(
	w http.ResponseWriter,
	data any,
) {
	status := http.StatusBadRequest

	Error(
		w,
		http.StatusText(status),
		data,
		status,
	)
}

func ErrorUnauthorized(
	w http.ResponseWriter,
	data any,
) {
	status := http.StatusUnauthorized

	Error(
		w,
		http.StatusText(status),
		data,
		status,
	)
}

func ErrorForbidden(
	w http.ResponseWriter,
	data any,
) {
	status := http.StatusForbidden

	Error(
		w,
		http.StatusText(status),
		data,
		status,
	)
}

func ErrorNotFound(
	w http.ResponseWriter,
	data any,
) {
	status := http.StatusNotFound

	Error(
		w,
		http.StatusText(status),
		data,
		status,
	)
}

func ErrorConflict(
	w http.ResponseWriter,
	data any,
) {
	status := http.StatusConflict

	Error(
		w,
		http.StatusText(status),
		data,
		status,
	)
}

func ErrorInternalServer(
	w http.ResponseWriter,
	data any,
) {
	status := http.StatusInternalServerError

	Error(
		w,
		http.StatusText(status),
		data,
		status,
	)
}
