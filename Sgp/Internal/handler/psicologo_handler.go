package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sgp/Internal/model"
	"sgp/Internal/repository"
	"time"
)

const TimeoutTime = 5 * time.Second

type PsicologoHandler struct {
	Repo repository.PsicologoRepository
}

func NewPsicologoRepository(repo repository.PsicologoRepository) *PsicologoHandler {
	return &PsicologoHandler { Repo: repo }
}

func (h *PsicologoHandler) HandlerCriarPsicologo(
	w http.ResponseWriter,
	r *http.Request,
) {
	var psicologo model.Psicologo
	err := json.NewDecoder(r.Body).Decode(&psicologo)

	if err != nil {
		// tem que tratar isso aí		
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	createdPsicologo, err := h.Repo.CriarPsicologo(ctx, psicologo)
	if err != nil {
		// tem que tratar isso aí
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPsicologo)
}