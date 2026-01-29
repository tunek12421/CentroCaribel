package dto

import (
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	if err := validator.RequiredString(r.Email, "email"); err != nil {
		return err
	}
	if err := validator.ValidEmail(r.Email); err != nil {
		return err
	}
	if err := validator.RequiredString(r.Password, "password"); err != nil {
		return err
	}
	return nil
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *RefreshRequest) Validate() error {
	if r.RefreshToken == "" {
		return apperrors.NewBadRequest("El campo 'refresh_token' es requerido")
	}
	return nil
}
