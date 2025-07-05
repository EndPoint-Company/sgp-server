package handler

import(
	"context"
	"encoding/json"
	"net/http"
	"sgp/Internal/model"
	"sgp/Internal/repository"
	"time"
)

const Timeout = 5 * time.Second

type ConsultaHandler struct{
	Repo repository.ConsultaRepository
}

func NewConsultaHandler(repo repository.ConsultaRepository) *ConsultaHandler{
	return &ConsultaHandler{Repo:repo}
}



func (h *ConsultaHandler) HandlerAgendarConsulta(w http.ResponseWriter, r *http.Request){
	var consulta model.Consulta

	if err := json.NewDecoder(r.Body).Decode(&consulta); err != nil{
		http.Error(w, "Requisicao invalida", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if consulta.AlunoID == "" || consulta.PsicologoID == ""{
		http.Error(w, "campos alunoid e psicologoid sao obrigatorios", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), Timeout)
	defer cancel()

	novaConsulta, err := h.Repo.AgendarConsulta(ctx, consulta)

	if err != nil{
		http.Error(w, "erro ao agendar consulta no banco de dados", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(novaConsulta)
}

func (h *ConsultaHandler) HandlerListarConsultasPorPsicologo(w http.ResponseWriter, r * http.Request){
	psicologoId := r.URL.Query().Get("psicologoId")
	
	if psicologoId == ""{
		http.Error(w, "o psicologoId é obrigatorio", http.StatusBadRequest)
		return
	}

	statusFiltro := r.URL.Query().Get("status")

	ctx, cancel := context.WithTimeout(r.Context(), Timeout)
	defer cancel()

	consultas, err := h.Repo.ListarConsultasPorPsicologo(ctx, psicologoId, statusFiltro)

	if err != nil{
		http.Error(w, "erro ao lsitar consultas por psicologo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(consultas)
}

func (h *ConsultaHandler) HandlerAtualizarStatusConsulta(w http.ResponseWriter, r * http.Request){
	id := r.PathValue("id")

	if id == ""{
		http.Error(w, "o id da consulta e obrigatorio", http.StatusBadRequest)
		return
	}
	var payload struct{
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil{
		http.Error(w, "requisicao invalida", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()


	if payload.Status == ""{
		http.Error(w, "o campo status é obrigatorio", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), Timeout)
	defer cancel()

	if err := h.Repo.AtualizaStatusConsulta(ctx, id, payload.Status); err != nil {
		http.Error(w, "erro ao atualizar o status da consulta", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "status da consulta atualizado com sucesso"})
}

func (h *ConsultaHandler) HandlerDeletarConsulta(w http.ResponseWriter, r * http.Request){
	id := r.PathValue("id")
	if id == ""{
		http.Error(w, "o id da consulta e obrigatorio", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), Timeout)
	defer cancel()

	if err := h.Repo.DeletarConsulta(ctx, id); err != nil{
		http.Error(w, "erro ao deletar consulta", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}