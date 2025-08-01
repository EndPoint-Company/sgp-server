package handler

import (
	"context"
	"encoding/json"
	"log" // MODIFICADO: Garantir que o pacote de log está sendo usado
	"net/http"
	"sgp/Internal/model"
	"sgp/Internal/repository"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const TimeoutAluno = 5 * time.Second

type AlunoHandler struct {
	Repo repository.AlunoRepository
}

func NewAlunoHandler(repo repository.AlunoRepository) *AlunoHandler {
	return &AlunoHandler{Repo: repo}
}

func (h *AlunoHandler) HandlerCriarAluno(
	w http.ResponseWriter, r *http.Request,
) {
	var aluno model.Aluno
	if err := json.NewDecoder(r.Body).Decode(&aluno); err != nil {
		httpError(w, "Requisição inválida", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if aluno.Nome == "" || aluno.Email == "" {
		httpError(w, "Campos 'nome' e 'email' são obrigatórios",
			http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutAluno)
	defer cancel()

	createdAluno, err := h.Repo.CriarAluno(ctx, aluno)
	if err != nil {
		// MODIFICADO: Adicionado log de erro
		log.Printf("ERRO ao criar aluno no banco de dados: %v", err)
		httpError(w, "Erro ao criar aluno no banco de dados",
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdAluno)
}

func (h *AlunoHandler) HandlerListarAlunos(
	w http.ResponseWriter, r *http.Request,
) {
	ctx, cancel := context.WithTimeout(r.Context(), TimeoutAluno)
	defer cancel()

	alunos, err := h.Repo.ListarAlunos(ctx)
	if err != nil {
		// MODIFICADO: Adicionado log de erro
		log.Printf("ERRO ao listar alunos: %v", err)
		httpError(w, "Erro ao listar alunos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(alunos)
}

func (h *AlunoHandler) HandlerBuscarAlunoPorID(
	w http.ResponseWriter, r *http.Request,
) {
	id := r.PathValue("id")
	if id == "" {
		httpError(w, "O ID do aluno é obrigatório", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutAluno)
	defer cancel()

	aluno, err := h.Repo.BuscarAlunoPorID(ctx, id)
	if err != nil {
		// MODIFICADO: Adicionado log de erro
		log.Printf("ERRO ao buscar aluno por ID (%s): %v", id, err)
		if status.Code(err) == codes.NotFound {
			httpError(w, "Aluno não encontrado", http.StatusNotFound)
		} else {
			httpError(w, "Erro ao buscar aluno", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(aluno)
}

func (h *AlunoHandler) HandlerAtualizarAluno(
	w http.ResponseWriter, r *http.Request,
) {
	id := r.PathValue("id")
	if id == "" {
		httpError(w, "O ID do aluno é obrigatório", http.StatusBadRequest)
		return
	}

	var aluno model.Aluno
	if err := json.NewDecoder(r.Body).Decode(&aluno); err != nil {
		httpError(w, "Requisição inválida", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutAluno)
	defer cancel()

	if err := h.Repo.AtualizarAluno(ctx, id, aluno); err != nil {
		// MODIFICADO: Adicionado log de erro
		log.Printf("ERRO ao atualizar aluno por ID (%s): %v", id, err)
		if status.Code(err) == codes.NotFound {
			httpError(w, "Não é possível atualizar um aluno que não existe.",
				http.StatusNotFound)
		} else {
			httpError(w, "Erro ao atualizar aluno.",
				http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]string{"message": "Aluno atualizado com sucesso"},
	)
}

func (h *AlunoHandler) HandlerDeletarAluno(
	w http.ResponseWriter, r *http.Request,
) {
	id := r.PathValue("id")
	if id == "" {
		httpError(w, "O ID do aluno é obrigatório", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutAluno)
	defer cancel()

	if err := h.Repo.DeletarAluno(ctx, id); err != nil {
		// MODIFICADO: Adicionado log de erro
		log.Printf("ERRO ao deletar aluno por ID (%s): %v", id, err)
		httpError(w, "Erro ao deletar aluno",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AlunoHandler) HandlerBuscarAlunoPorNome(
	w http.ResponseWriter, r *http.Request,
) {
	nome := r.URL.Query().Get("nome")
	if nome == "" {
		httpError(w, "O 'nome' é obrigatório",
			http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), TimeoutAluno)
	defer cancel()

	id, err := h.Repo.GetAlunoIDPorNome(ctx, nome)
	if err != nil {
		// MODIFICADO: Adicionado log de erro
		log.Printf("ERRO ao buscar aluno por nome (%s): %v", nome, err)
		httpError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}




// A função httpError não foi fornecida, mas estou assumindo que ela existe
// em algum lugar do seu pacote 'handler' para que o código compile.
// Exemplo de como ela poderia ser:
/*
func httpError(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}
*/