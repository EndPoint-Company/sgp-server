package mocks

import (
	"context"
	"sgp/Internal/model"
	"sgp/Internal/repository"
)

// Garante que o Mock implementa a interface PsicologoRepository.
var _ repository.PsicologoRepository = &PsicologoRepositoryMock{}

type PsicologoRepositoryMock struct {
	CriarPsicologoFunc        func(ctx context.Context, psicologo model.Psicologo) (*model.Psicologo, error)
	ListarPsicologosFunc      func(ctx context.Context) ([]*model.Psicologo, error)
	BuscarPsicologoPorIDFunc  func(ctx context.Context, id string) (*model.Psicologo, error)
	AtualizarPsicologoFunc    func(ctx context.Context, id string, psicologo model.Psicologo) error
	DeletarPsicologoFunc      func(ctx context.Context, id string) error
	GetPsicologoIDPorNomeFunc func(ctx context.Context, nome string) (string, error)
}

func (m *PsicologoRepositoryMock) CriarPsicologo(ctx context.Context, p model.Psicologo) (*model.Psicologo, error) {
	return m.CriarPsicologoFunc(ctx, p)
}

func (m *PsicologoRepositoryMock) ListarPsicologos(ctx context.Context) ([]*model.Psicologo, error) {
	return m.ListarPsicologosFunc(ctx)
}

func (m *PsicologoRepositoryMock) BuscarPsicologoPorID(ctx context.Context, id string) (*model.Psicologo, error) {
	return m.BuscarPsicologoPorIDFunc(ctx, id)
}

func (m *PsicologoRepositoryMock) AtualizarPsicologo(ctx context.Context, id string, p model.Psicologo) error {
	return m.AtualizarPsicologoFunc(ctx, id, p)
}

func (m *PsicologoRepositoryMock) DeletarPsicologo(ctx context.Context, id string) error {
	return m.DeletarPsicologoFunc(ctx, id)
}

func (m *PsicologoRepositoryMock) GetPsicologoIDPorNome(ctx context.Context, nome string) (string, error) {
	return m.GetPsicologoIDPorNomeFunc(ctx, nome)
}