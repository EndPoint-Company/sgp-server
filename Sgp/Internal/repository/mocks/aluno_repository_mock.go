package mocks

import (
	"context"
	"sgp/Internal/model"
	"sgp/Internal/repository"
)

var _ repository.AlunoRepository = &AlunoRepositoryMock{}
type AlunoRepositoryMock struct{
	CriarAlunoFunc func(ctx context.Context, aluno model.Aluno)(*model.Aluno, error)
	ListarAlunosFunc func(ctx context.Context)([]*model.Aluno, error)
	AtualizarAlunosFunc func(ctx context.Context, id string, aluno model.Aluno)(error)
	BuscarAlunoPorIDFunc func(ctx context.Context, id string)(*model.Aluno, error)
	DeletarAlunoFunc func(ctx context.Context, id string)(error)
	GetAlunoIDPorNomeFunc func(ctx context.Context, nome string)(string, error)
}


func (m *AlunoRepositoryMock) CriarAluno(ctx context.Context, aluno model.Aluno) (*model.Aluno, error) {
	return m.CriarAlunoFunc(ctx, aluno)
}

func (m *AlunoRepositoryMock) ListarAlunos(ctx context.Context) ([]*model.Aluno, error) {
	return m.ListarAlunosFunc(ctx)
}

func (m *AlunoRepositoryMock) BuscarAlunoPorID(ctx context.Context, id string) (*model.Aluno, error) {
	return m.BuscarAlunoPorIDFunc(ctx, id)
}

func (m *AlunoRepositoryMock) AtualizarAluno(ctx context.Context, id string, aluno model.Aluno) error {
	return m.AtualizarAlunosFunc(ctx, id, aluno)
}

func (m *AlunoRepositoryMock) DeletarAluno(ctx context.Context, id string) error {
	return m.DeletarAlunoFunc(ctx, id)
}

func (m *AlunoRepositoryMock) GetAlunoIDPorNome(ctx context.Context, nome string) (string, error) {
	return m.GetAlunoIDPorNomeFunc(ctx, nome)
}