// sgp/Internal/handler/consulta_handler.go

package handler

import (
	"encoding/json"
	"log" // MODIFICADO: Garantir que o pacote de log está sendo usado
	"net/http"
	"sgp/Internal/model"
	"sgp/Internal/repository"
	"time"
    "context"
)

const Timeout = 5 * time.Second

type ConsultaHandler struct {
	Repo repository.ConsultaRepository
}

func NewConsultaHandler(repo repository.ConsultaRepository) *ConsultaHandler {
	return &ConsultaHandler{Repo: repo}
}

func (h *ConsultaHandler) HandlerAgendarConsulta(w http.ResponseWriter, r *http.Request) {
	var consulta model.Consulta

	if err := json.NewDecoder(r.Body).Decode(&consulta); err != nil {
		http.Error(w, "Requisicao invalida", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if consulta.AlunoID == "" || consulta.PsicologoID == "" {
		http.Error(w, "campos alunoid e psicologoid sao obrigatorios", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), Timeout)
	defer cancel()

	novaConsulta, err := h.Repo.AgendarConsulta(ctx, consulta)

	if err != nil {
        // MODIFICADO: Adicionado log para registrar o erro no terminal
		log.Printf("ERRO ao agendar consulta: %v", err)
		http.Error(w, "erro ao agendar consulta no banco de dados", http.StatusInternalServerError)
        return // MODIFICADO: Adicionado 'return' para parar a execução
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(novaConsulta)
}

func (h *ConsultaHandler) HandlerListarConsultasPorPsicologo(w http.ResponseWriter, r *http.Request) {
    // MODIFICADO: Adicionado log de rastreamento no início
	log.Println("--- INÍCIO: HandlerListarConsultasPorPsicologo foi chamado ---")
    
	psicologoId := r.URL.Query().Get("psicologoId")

	if psicologoId == "" {
		http.Error(w, "o psicologoId é obrigatorio", http.StatusBadRequest)
		return
	}
    // MODIFICADO: Adicionado log para ver o ID recebido
    log.Printf("Buscando consultas para o psicologoId: %s", psicologoId)

	statusFiltro := r.URL.Query().Get("status")

	ctx, cancel := context.WithTimeout(r.Context(), Timeout)
	defer cancel()

	consultas, err := h.Repo.ListarConsultasPorPsicologo(ctx, psicologoId, statusFiltro)

	if err != nil {
        // MODIFICADO: Adicionado log para registrar o erro no terminal
		log.Printf("ERRO ao listar consultas por psicologo: %v", err)
		http.Error(w, "erro ao listar consultas por psicologo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(consultas)
    // MODIFICADO: Adicionado log de rastreamento no fim
    log.Println("--- FIM: HandlerListarConsultasPorPsicologo finalizado com sucesso ---")
}


func (h *ConsultaHandler) HandlerListarConsultasPorAluno(w http.ResponseWriter, r *http.Request) {
	// MODIFICADO: Log de início do handler
	log.Println("--- INÍCIO: HandlerListarConsultasPorAluno foi chamado ---")

	alunoId := r.URL.Query().Get("alunoId")
	if alunoId == "" {
		http.Error(w, "o alunoId é obrigatorio", http.StatusBadRequest)
		return
	}

	// MODIFICADO: Log do alunoId recebido
	log.Printf("Buscando consultas para o alunoId: %s", alunoId)

	ctx, cancel := context.WithTimeout(r.Context(), Timeout)
	defer cancel()

	consultas, err := h.Repo.ListarConsultasPorAluno(ctx, alunoId)
	if err != nil {
		log.Printf("ERRO ao listar consultas por aluno: %v", err)
		http.Error(w, "erro ao listar consultas por aluno", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(consultas)

	// MODIFICADO: Log de finalização
	log.Println("--- FIM: HandlerListarConsultasPorAluno finalizado com sucesso ---")
}


func (h *ConsultaHandler) HandlerAtualizarStatusConsulta(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		http.Error(w, "o id da consulta e obrigatorio", http.StatusBadRequest)
		return
	}
	var payload struct {
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "requisicao invalida", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if payload.Status == "" {
		http.Error(w, "o campo status é obrigatorio", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), Timeout)
	defer cancel()

	if err := h.Repo.AtualizaStatusConsulta(ctx, id, payload.Status); err != nil {
        // MODIFICADO: Adicionado log para registrar o erro no terminal
		log.Printf("ERRO ao atualizar status da consulta: %v", err)
		http.Error(w, "erro ao atualizar o status da consulta", http.StatusInternalServerError)
        return // MODIFICADO: Adicionado 'return' para parar a execução
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "status da consulta atualizado com sucesso"})
}

func (h *ConsultaHandler) HandlerDeletarConsulta(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "o id da consulta e obrigatorio", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), Timeout)
	defer cancel()

	if err := h.Repo.DeletarConsulta(ctx, id); err != nil {
        // MODIFICADO: Adicionado log para registrar o erro no terminal
		log.Printf("ERRO ao deletar consulta: %v", err)
		http.Error(w, "erro ao deletar consulta", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}