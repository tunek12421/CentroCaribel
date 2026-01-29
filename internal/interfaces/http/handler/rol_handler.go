package handler

import (
	"net/http"

	"github.com/tunek/centro-caribel/internal/domain"
	"github.com/tunek/centro-caribel/pkg/response"
)

type RolHandler struct {
	repo domain.RolRepository
}

func NewRolHandler(repo domain.RolRepository) *RolHandler {
	return &RolHandler{repo: repo}
}

func (h *RolHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	roles, err := h.repo.GetAll(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.JSON(w, http.StatusOK, roles)
}
