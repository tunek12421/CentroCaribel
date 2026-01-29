package handler

import (
	"net/http"

	"github.com/tunek/centro-caribel/internal/application/auth"
	"github.com/tunek/centro-caribel/internal/interfaces/http/dto"
	"github.com/tunek/centro-caribel/pkg/response"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	token, refreshToken, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    28800, // 8 horas en segundos
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	token, refreshToken, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    28800,
	})
}
