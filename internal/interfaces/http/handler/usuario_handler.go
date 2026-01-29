package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/tunek/centro-caribel/internal/application/usuario"
	"github.com/tunek/centro-caribel/internal/interfaces/http/dto"
	apperrors "github.com/tunek/centro-caribel/pkg/errors"
	"github.com/tunek/centro-caribel/pkg/response"
	"github.com/tunek/centro-caribel/pkg/validator"
)

type UsuarioHandler struct {
	service *usuario.Service
}

func NewUsuarioHandler(service *usuario.Service) *UsuarioHandler {
	return &UsuarioHandler{service: service}
}

func (h *UsuarioHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUsuarioRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	user, err := h.service.Create(r.Context(), req.NombreCompleto, req.Email, req.Password, req.RolID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

func (h *UsuarioHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	usuarios, total, err := h.service.GetAll(r.Context(), page, perPage)
	if err != nil {
		response.Error(w, err)
		return
	}

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	response.JSONWithMeta(w, http.StatusOK, usuarios, &response.Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     total,
		TotalPage: totalPages,
	})
}

func (h *UsuarioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID inválido"))
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UsuarioHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID inválido"))
		return
	}

	var req dto.UpdateUsuarioRequest
	if err := validator.DecodeAndValidate(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	user, err := h.service.Update(r.Context(), id, req.NombreCompleto, req.Email, req.RolID, req.Activo)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *UsuarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, apperrors.NewBadRequest("ID inválido"))
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Usuario eliminado"})
}
