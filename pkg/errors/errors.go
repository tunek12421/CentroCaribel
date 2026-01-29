package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Detail)
}

func NewBadRequest(detail string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: "Solicitud inválida", Detail: detail}
}

func NewNotFound(entity string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: "No encontrado", Detail: fmt.Sprintf("%s no encontrado", entity)}
}

func NewUnauthorized(detail string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: "No autorizado", Detail: detail}
}

func NewForbidden(detail string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: "Acceso denegado", Detail: detail}
}

func NewConflict(detail string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: "Conflicto", Detail: detail}
}

func NewInternal(detail string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: "Error interno", Detail: detail}
}
