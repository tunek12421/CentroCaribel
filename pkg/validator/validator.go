package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	apperrors "github.com/tunek/centro-caribel/pkg/errors"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func DecodeAndValidate(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return apperrors.NewBadRequest("El cuerpo de la solicitud está vacío")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperrors.NewBadRequest("JSON inválido: " + err.Error())
	}

	if v, ok := dst.(Validatable); ok {
		return v.Validate()
	}
	return nil
}

type Validatable interface {
	Validate() error
}

func RequiredString(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return apperrors.NewBadRequest(fmt.Sprintf("El campo '%s' es requerido", field))
	}
	return nil
}

func ValidEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return apperrors.NewBadRequest("El email no es válido")
	}
	return nil
}

func MinLength(value, field string, min int) error {
	if len(value) < min {
		return apperrors.NewBadRequest(fmt.Sprintf("El campo '%s' debe tener al menos %d caracteres", field, min))
	}
	return nil
}
