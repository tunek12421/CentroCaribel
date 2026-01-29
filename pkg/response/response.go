package response

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/tunek/centro-caribel/pkg/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorBody struct {
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

type Meta struct {
	Page      int   `json:"page,omitempty"`
	PerPage   int   `json:"per_page,omitempty"`
	Total     int64 `json:"total,omitempty"`
	TotalPage int   `json:"total_pages,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&Response{Success: true, Data: data})
}

func JSONWithMeta(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&Response{Success: true, Data: data, Meta: meta})
}

func Error(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.Code)
		json.NewEncoder(w).Encode(&Response{
			Success: false,
			Error:   &ErrorBody{Message: appErr.Message, Detail: appErr.Detail},
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(&Response{
		Success: false,
		Error:   &ErrorBody{Message: "Error interno del servidor"},
	})
}
