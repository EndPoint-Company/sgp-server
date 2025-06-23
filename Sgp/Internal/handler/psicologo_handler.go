package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sgp/Internal/model"
	"sgp/Internal/repository"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const TimeoutTime = 5 * time.Second

type PsicologoHandler struct {
	Repo repository.PsicologoRepository
}

func NewPsicologoHandler(repo repository.PsicologoRepository) *PsicologoHandler {
	return &PsicologoHandler{Repo: repo}
}

func httpError(w http.ResponseWriter, message string, code int) {
	log.Printf("Erro na requisição: %s", message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *PsicologoHandler) HandlerCriarPsicologo(
	w http.ResponseWriter, r *http.Request,
) {
	var psicologo model.Psicologo
	if err := json.NewDecoder(r.Body).Decode(&psicologo); err != nil {
		httpError(w, "Requisição inválida", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if psicologo.Nome == "" || psicologo.Email == "" || psicologo.CRP == "" {
		httpError(w, "Campos 'nome', 'email' e 'crp' são obrigatórios",
			http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	createdPsicologo, err := h.Repo.CriarPsicologo(ctx, psicologo)
	if err != nil {
		httpError(w, "Erro ao criar psicólogo no banco de dados",
			http.StatusInternalServerError)
			
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPsicologo)
}

func (h *PsicologoHandler) HandlerListarPsicologos(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	psicologos, err := h.Repo.ListarPsicologos(ctx)
	if err != nil {
		httpError(w, "Erro ao listar psicólogos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(psicologos)
}

func (h *PsicologoHandler) HandlerBuscarPsicologoPorID(
	w http.ResponseWriter, r *http.Request,
) {
	id := r.PathValue("id")
	if id == "" {
		httpError(w, "O ID do psicólogo é obrigatório", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	psicologo, err := h.Repo.BuscarPsicologoPorID(ctx, id)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			httpError(w, "Psicólogo não encontrado", http.StatusNotFound)
		} else {
			httpError(w, "Erro ao buscar psicólogo", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(psicologo)
}

func (h *PsicologoHandler) HandlerAtualizarPsicologo(
	w http.ResponseWriter, r *http.Request,
) {
	id := r.PathValue("id")
	if id == "" {
		httpError(w, "O ID do psicólogo é obrigatório", http.StatusBadRequest)
		return
	}

	var psicologo model.Psicologo
	if err := json.NewDecoder(r.Body).Decode(&psicologo); err != nil {
		httpError(w, "Requisição inválida", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	if err := h.Repo.AtualizarPsicologo(ctx, id, psicologo); err != nil {
        if status.Code(err) == codes.NotFound {
			httpError(w, "Não é possível atualizar um psicólogo que não existe.", http.StatusNotFound)
		} else {
            httpError(w, "Erro ao atualizar psicólogo.", http.StatusInternalServerError)
        }

		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]string{"message": "Psicólogo atualizado com sucesso"},
	)
}

func (h *PsicologoHandler) HandlerDeletarPsicologo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		httpError(w, "O ID do psicólogo é obrigatório", http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	if err := h.Repo.DeletarPsicologo(ctx, id); err != nil {
		httpError(w, "Erro ao deletar psicólogo", 
			http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PsicologoHandler) HandlerBuscarPsicologoPorNome(w http.ResponseWriter, r *http.Request) {
	nome := r.URL.Query().Get("nome")
	if nome == "" {
		httpError(w, "O 'nome' é obrigatório",
			http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutTime)
	defer cancel()

	id, err := h.Repo.GetPsicologoIDPorNome(ctx, nome)
	if err != nil {
		httpError(w, err.Error(), http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}
