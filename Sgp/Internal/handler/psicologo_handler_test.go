package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sgp/Internal/model"
	"sgp/Internal/repository/mocks"
	"testing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlerCriarPsicologo(t *testing.T) {
	psicologo := model.Psicologo{Nome: "Dr. Freud", Email: "freud@sigmund.com", CRP: "01/12345"}
	psicologoJSON, _ := json.Marshal(psicologo)

	t.Run("sucesso ao criar psicologo", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/psicologos", bytes.NewBuffer(psicologoJSON))
		rr := httptest.NewRecorder()

		mockRepo := &mocks.PsicologoRepositoryMock{
			CriarPsicologoFunc: func(ctx context.Context, p model.Psicologo) (*model.Psicologo, error) {
				return &model.Psicologo{ID: "psico-123", Nome: p.Nome, Email: p.Email, CRP: p.CRP}, nil
			},
		}

		h := NewPsicologoHandler(mockRepo)
		h.HandlerCriarPsicologo(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusCreated)
		}
	})

	t.Run("campos faltando", func(t *testing.T) {
		psicologoIncompleto := model.Psicologo{Nome: "Incompleto"}
		jsonIncompleto, _ := json.Marshal(psicologoIncompleto)
		req, _ := http.NewRequest("POST", "/psicologos", bytes.NewBuffer(jsonIncompleto))
		rr := httptest.NewRecorder()

		h := NewPsicologoHandler(&mocks.PsicologoRepositoryMock{})
		h.HandlerCriarPsicologo(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusBadRequest)
		}
	})
}

func TestHandlerListarPsicologos(t *testing.T) {
	t.Run("sucesso ao listar psicologos", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/psicologos", nil)
		rr := httptest.NewRecorder()

		mockRepo := &mocks.PsicologoRepositoryMock{
			ListarPsicologosFunc: func(ctx context.Context) ([]*model.Psicologo, error) {
				return []*model.Psicologo{
					{ID: "1", Nome: "Psico 1"},
					{ID: "2", Nome: "Psico 2"},
				}, nil
			},
		}
		h := NewPsicologoHandler(mockRepo)
		h.HandlerListarPsicologos(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusOK)
		}
		var psicologos []*model.Psicologo
		json.NewDecoder(rr.Body).Decode(&psicologos)
		if len(psicologos) != 2 {
			t.Errorf("número incorreto de psicólogos: obteve %d, esperava 2", len(psicologos))
		}
	})
}

func TestHandlerBuscarPsicologoPorID(t *testing.T) {
	t.Run("sucesso ao buscar por id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/psicologos/123", nil)
		req.SetPathValue("id", "123")
		rr := httptest.NewRecorder()

		mockRepo := &mocks.PsicologoRepositoryMock{
			BuscarPsicologoPorIDFunc: func(ctx context.Context, id string) (*model.Psicologo, error) {
				return &model.Psicologo{ID: id, Nome: "Encontrado"}, nil
			},
		}
		h := NewPsicologoHandler(mockRepo)
		h.HandlerBuscarPsicologoPorID(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusOK)
		}
	})

	t.Run("psicologo nao encontrado", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/psicologos/404", nil)
		req.SetPathValue("id", "404")
		rr := httptest.NewRecorder()
		mockRepo := &mocks.PsicologoRepositoryMock{
			BuscarPsicologoPorIDFunc: func(ctx context.Context, id string) (*model.Psicologo, error) {
				return nil, status.Error(codes.NotFound, "não encontrado")
			},
		}
		h := NewPsicologoHandler(mockRepo)
		h.HandlerBuscarPsicologoPorID(rr, req)
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusNotFound)
		}
	})
}

func TestHandlerAtualizarPsicologo(t *testing.T) {
	psicologoAtualizado := model.Psicologo{Nome: "Dr. Jung"}
	jsonBody, _ := json.Marshal(psicologoAtualizado)

	t.Run("sucesso ao atualizar", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/psicologos/123", bytes.NewBuffer(jsonBody))
		req.SetPathValue("id", "123")
		rr := httptest.NewRecorder()

		mockRepo := &mocks.PsicologoRepositoryMock{
			AtualizarPsicologoFunc: func(ctx context.Context, id string, p model.Psicologo) error {
				return nil // Sucesso
			},
		}
		h := NewPsicologoHandler(mockRepo)
		h.HandlerAtualizarPsicologo(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusOK)
		}
	})

	t.Run("psicologo a ser atualizado nao encontrado", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/psicologos/404", bytes.NewBuffer(jsonBody))
		req.SetPathValue("id", "404")
		rr := httptest.NewRecorder()
		mockRepo := &mocks.PsicologoRepositoryMock{
			AtualizarPsicologoFunc: func(ctx context.Context, id string, p model.Psicologo) error {
				return status.Error(codes.NotFound, "não encontrado")
			},
		}
		h := NewPsicologoHandler(mockRepo)
		h.HandlerAtualizarPsicologo(rr, req)
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusNotFound)
		}
	})
}

func TestHandlerDeletarPsicologo(t *testing.T) {
	t.Run("sucesso ao deletar", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/psicologos/123", nil)
		req.SetPathValue("id", "123")
		rr := httptest.NewRecorder()
		mockRepo := &mocks.PsicologoRepositoryMock{
			DeletarPsicologoFunc: func(ctx context.Context, id string) error {
				return nil // Sucesso
			},
		}
		h := NewPsicologoHandler(mockRepo)
		h.HandlerDeletarPsicologo(rr, req)
		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusNoContent)
		}
	})
}

func TestHandlerBuscarPsicologoPorNome(t *testing.T) {
	t.Run("sucesso ao buscar por nome", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/psicologos/nome?nome=Dr.+Adler", nil)
		rr := httptest.NewRecorder()
		mockRepo := &mocks.PsicologoRepositoryMock{
			GetPsicologoIDPorNomeFunc: func(ctx context.Context, nome string) (string, error) {
				return "psico-id-456", nil
			},
		}
		h := NewPsicologoHandler(mockRepo)
		h.HandlerBuscarPsicologoPorNome(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusOK)
		}
		var resp map[string]string
		json.NewDecoder(rr.Body).Decode(&resp)
		if resp["id"] != "psico-id-456" {
			t.Errorf("ID do psicólogo incorreto: obteve %s", resp["id"])
		}
	})

	t.Run("nome nao encontrado", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/psicologos/nome?nome=Inexistente", nil)
		rr := httptest.NewRecorder()
		mockRepo := &mocks.PsicologoRepositoryMock{
			GetPsicologoIDPorNomeFunc: func(ctx context.Context, nome string) (string, error) {
				return "", errors.New("não encontrado")
			},
		}
		h := NewPsicologoHandler(mockRepo)
		h.HandlerBuscarPsicologoPorNome(rr, req)
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusNotFound)
		}
	})
}