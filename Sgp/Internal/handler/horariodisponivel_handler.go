package handler

import (
	"encoding/json"
	"net/http"
	"sgp/Internal/model"
	"sgp/Internal/repository"
	"time"
	"context"
)

type HorarioDisponivelHandler struct {
	Repo repository.HorarioDisponivelRepository
}

func NewHorarioDisponivelHandler(repo repository.HorarioDisponivelRepository) *HorarioDisponivelHandler {
	return &HorarioDisponivelHandler{Repo: repo}
}

func (h *HorarioDisponivelHandler) HandlerCriarHorario(w http.ResponseWriter, r *http.Request) {
	var horario model.HorarioDisponivel
	
	if err := json.NewDecoder(r.Body).Decode(&horario); err != nil {
		httpError(w, "Requisição inválida", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if horario.PsicologoID == "" || horario.Inicio.IsZero() || horario.Fim.IsZero() {
		httpError(w, "Campos 'psicologoId', 'inicio' e 'fim' são obrigatórios", http.StatusBadRequest)
		return
	}

	if horario.Status != "bloqueado" {
		horario.Status = "disponivel"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	novoHorario, err := h.Repo.CriarHorario(ctx, horario)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(novoHorario)
}

func (h *HorarioDisponivelHandler) HandlerListarHorarios(w http.ResponseWriter, r *http.Request) {
	psicologoId := r.URL.Query().Get("psicologoId")
	status := r.URL.Query().Get("status")

	if psicologoId == "" {
		httpError(w, "O 'psicologoId' é obrigatório", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	horarios, err := h.Repo.ListarHorariosPorPsicologo(ctx, psicologoId, status)
	if err != nil {
		httpError(w, "Erro ao listar horários", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(horarios)
}

func (h *HorarioDisponivelHandler) HandlerDeletarHorario(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpError(w, "O ID do horário é obrigatório", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.Repo.DeletarHorario(ctx, id); err != nil {
		httpError(w, "Erro ao deletar horário/bloqueio", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}