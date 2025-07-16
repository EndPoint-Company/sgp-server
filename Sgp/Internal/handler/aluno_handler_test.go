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

func TestHandlerCriarAluno(t *testing.T) {
	aluno := model.Aluno{Nome: "John Doe", Email: "john.doe@example.com"}
	alunoJSON, _ := json.Marshal(aluno)

	t.Run("sucesso ao criar aluno", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/alunos", bytes.NewBuffer(alunoJSON))
		rr := httptest.NewRecorder()
		mockRepo := &mocks.AlunoRepositoryMock{
			CriarAlunoFunc: func(ctx context.Context, a model.Aluno) (*model.Aluno, error) {
				return &model.Aluno{ID: "123", Nome: a.Nome, Email: a.Email}, nil
			},
		}

		h := NewAlunoHandler(mockRepo)
		h.HandlerCriarAluno(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusCreated)
		}

		var createdAluno model.Aluno
		json.NewDecoder(rr.Body).Decode(&createdAluno)
		if createdAluno.Nome != aluno.Nome {
			t.Errorf("nome do aluno incorreto: obteve %s, esperava %s", createdAluno.Nome, aluno.Nome)
		}
	})

	t.Run("erro no corpo da requisição", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/alunos", bytes.NewBufferString("corpo invalido"))
		rr := httptest.NewRecorder()
		h := NewAlunoHandler(&mocks.AlunoRepositoryMock{})
		h.HandlerCriarAluno(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusBadRequest)
		}
	})
}

func TestHandlerListarAlunos(t *testing.T) {
	t.Run("sucesso ao listar alunos", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/alunos", nil)
		rr := httptest.NewRecorder()
		mockRepo := &mocks.AlunoRepositoryMock{
			ListarAlunosFunc: func(ctx context.Context) ([]*model.Aluno, error) {
				return []*model.Aluno{
					{ID: "1", Nome: "Aluno 1", Email: "aluno1@test.com"},
					{ID: "2", Nome: "Aluno 2", Email: "aluno2@test.com"},
				}, nil
			},
		}

		h := NewAlunoHandler(mockRepo)
		h.HandlerListarAlunos(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusOK)
		}

		var alunos []*model.Aluno
		json.NewDecoder(rr.Body).Decode(&alunos)
		if len(alunos) != 2 {
			t.Errorf("número incorreto de alunos: obteve %d, esperava 2", len(alunos))
		}
	})
}

func TestHandlerBuscarAlunoPorID(t *testing.T) {
	t.Run("sucesso ao buscar aluno por id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/alunos/123", nil)
		req.SetPathValue("id", "123")
		rr := httptest.NewRecorder()

		mockRepo := &mocks.AlunoRepositoryMock{
			BuscarAlunoPorIDFunc: func(ctx context.Context, id string) (*model.Aluno, error) {
				if id == "123" {
					return &model.Aluno{ID: "123", Nome: "Aluno Encontrado"}, nil
				}
				return nil, status.Error(codes.NotFound, "não encontrado")
			},
		}

		h := NewAlunoHandler(mockRepo)
		h.HandlerBuscarAlunoPorID(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusOK)
		}
	})

	t.Run("aluno não encontrado", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/alunos/404", nil)
		req.SetPathValue("id", "404")
		rr := httptest.NewRecorder()

		mockRepo := &mocks.AlunoRepositoryMock{
			BuscarAlunoPorIDFunc: func(ctx context.Context, id string) (*model.Aluno, error) {
				return nil, status.Error(codes.NotFound, "não encontrado")
			},
		}

		h := NewAlunoHandler(mockRepo)
		h.HandlerBuscarAlunoPorID(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusNotFound)
		}
	})
}

func TestHandlerAtualizarAluno(t *testing.T) {
	alunoAtualizado := model.Aluno{Nome: "Nome Atualizado"}
	alunoJSON, _ := json.Marshal(alunoAtualizado)

	t.Run("sucesso ao atualizar aluno", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/alunos/123", bytes.NewBuffer(alunoJSON))
		req.SetPathValue("id", "123")
		rr := httptest.NewRecorder()

		mockRepo := &mocks.AlunoRepositoryMock{
			AtualizarAlunosFunc: func(ctx context.Context, id string, a model.Aluno) error {
				return nil // Simula sucesso
			},
		}

		h := NewAlunoHandler(mockRepo)
		h.HandlerAtualizarAluno(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusOK)
		}
	})

	t.Run("aluno a ser atualizado não encontrado", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/alunos/404", bytes.NewBuffer(alunoJSON))
		req.SetPathValue("id", "404")
		rr := httptest.NewRecorder()

		mockRepo := &mocks.AlunoRepositoryMock{
			AtualizarAlunosFunc: func(ctx context.Context, id string, a model.Aluno) error {
				return status.Error(codes.NotFound, "não encontrado")
			},
		}

		h := NewAlunoHandler(mockRepo)
		h.HandlerAtualizarAluno(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusNotFound)
		}
	})
}

func TestHandlerDeletarAluno(t *testing.T) {
	t.Run("sucesso ao deletar aluno", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/alunos/123", nil)
		req.SetPathValue("id", "123")
		rr := httptest.NewRecorder()

		mockRepo := &mocks.AlunoRepositoryMock{
			DeletarAlunoFunc: func(ctx context.Context, id string) error {
				return nil // Simula sucesso
			},
		}
		h := NewAlunoHandler(mockRepo)
		h.HandlerDeletarAluno(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusNoContent)
		}
	})
}

func TestHandlerBuscarAlunoPorNome(t *testing.T) {
	t.Run("sucesso ao buscar aluno por nome", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/alunos/nome?nome=Busca", nil)
		rr := httptest.NewRecorder()

		mockRepo := &mocks.AlunoRepositoryMock{
			GetAlunoIDPorNomeFunc: func(ctx context.Context, nome string) (string, error) {
				return "aluno-id-123", nil
			},
		}
		h := NewAlunoHandler(mockRepo)
		h.HandlerBuscarAlunoPorNome(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusOK)
		}
		var resp map[string]string
		json.NewDecoder(rr.Body).Decode(&resp)
		if resp["id"] != "aluno-id-123" {
			t.Errorf("ID do aluno incorreto: obteve %s, esperava aluno-id-123", resp["id"])
		}
	})

	t.Run("nome não encontrado", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/alunos/nome?nome=Inexistente", nil)
		rr := httptest.NewRecorder()

		mockRepo := &mocks.AlunoRepositoryMock{
			GetAlunoIDPorNomeFunc: func(ctx context.Context, nome string) (string, error) {
				return "", errors.New("não encontrado")
			},
		}
		h := NewAlunoHandler(mockRepo)
		h.HandlerBuscarAlunoPorNome(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("status code incorreto: obteve %v, esperava %v", status, http.StatusNotFound)
		}
	})
}