package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sgp/Internal/model"
	"sgp/Internal/repository/mocks"
	"testing"
	"time"
)

func TestHandlerAgendarConsulta(t *testing.T) {
	consulta := model.Consulta{
		AlunoID:     "aluno-1",
		PsicologoID: "psico-1",
		Horario:     time.Now(),
	}
	consultaJSON, _ := json.Marshal(consulta)

	t.Run("sucesso ao agendar", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/consultas", bytes.NewBuffer(consultaJSON))
		rr := httptest.NewRecorder()

		mockRepo := &mocks.ConsultaRepositoryMock{
			AgendarConsultaFunc: func(ctx context.Context, c model.Consulta) (*model.Consulta, error) {
				c.ID = "consulta-123"
				c.Status = "aguardando aprovacao"
				return &c, nil
			},
		}

		h := NewConsultaHandler(mockRepo)
		h.HandlerAgendarConsulta(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusCreated)
		}

		var novaConsulta model.Consulta
		json.NewDecoder(rr.Body).Decode(&novaConsulta)
		if novaConsulta.ID != "consulta-123" {
			t.Errorf("ID da consulta incorreto, esperava 'consulta-123', obteve '%s'", novaConsulta.ID)
		}
	})
}

func TestHandlerListarConsultasPorPsicologo(t *testing.T) {
	t.Run("sucesso ao listar", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/consultas/psicologo?psicologoId=psico-1&status=aprovada", nil)
		rr := httptest.NewRecorder()

		mockRepo := &mocks.ConsultaRepositoryMock{
			ListarConsultasPorPsicologoFunc: func(ctx context.Context, psicologoID string, statusFiltro string) ([]*model.Consulta, error) {
				if psicologoID == "psico-1" && statusFiltro == "aprovada" {
					return []*model.Consulta{{ID: "1"}, {ID: "2"}}, nil
				}
				return nil, nil
			},
		}

		h := NewConsultaHandler(mockRepo)
		h.HandlerListarConsultasPorPsicologo(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusOK)
		}
		var consultas []*model.Consulta
		json.NewDecoder(rr.Body).Decode(&consultas)
		if len(consultas) != 2 {
			t.Errorf("número incorreto de consultas: obteve %d, esperava 2", len(consultas))
		}
	})
}

func TestHandlerAtualizarStatusConsulta(t *testing.T) {
	payload := map[string]string{"status": "confirmada"}
	jsonBody, _ := json.Marshal(payload)

	t.Run("sucesso ao atualizar status", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/consultas/consulta-1/status", bytes.NewBuffer(jsonBody))
		req.SetPathValue("id", "consulta-1")
		rr := httptest.NewRecorder()

		mockRepo := &mocks.ConsultaRepositoryMock{
			AtualizaStatusConsultaFunc: func(ctx context.Context, id string, novoStatus string) error {
				return nil // Sucesso
			},
		}

		h := NewConsultaHandler(mockRepo)
		h.HandlerAtualizarStatusConsulta(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusOK)
		}
	})
}

func TestHandlerDeletarConsulta(t *testing.T) {
	t.Run("sucesso ao deletar", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/consultas/consulta-1", nil)
		req.SetPathValue("id", "consulta-1")
		rr := httptest.NewRecorder()

		mockRepo := &mocks.ConsultaRepositoryMock{
			DeletarConsultaFunc: func(ctx context.Context, id string) error {
				return nil // Sucesso
			},
		}

		h := NewConsultaHandler(mockRepo)
		h.HandlerDeletarConsulta(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusNoContent)
		}
	})
}