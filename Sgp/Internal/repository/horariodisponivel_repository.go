package repository

import (
	"context"
	"fmt"
	"sgp/Internal/model"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type HorarioDisponivelRepositoryImpl struct {
	Client *firestore.Client
}

func NewHorarioDisponivelRepository(client *firestore.Client) HorarioDisponivelRepository {
	return &HorarioDisponivelRepositoryImpl{Client: client}
}

func (r *HorarioDisponivelRepositoryImpl) CriarHorario(ctx context.Context, horario model.HorarioDisponivel) (*model.HorarioDisponivel, error) {
	docRef, _, err := r.Client.Collection("horariosDisponiveis").Add(ctx, horario)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar horário: %w", err)
	}
	horario.ID = docRef.ID
	return &horario, nil
}

func (r *HorarioDisponivelRepositoryImpl) ListarHorariosPorPsicologo(ctx context.Context, psicologoID string, status string) ([]*model.HorarioDisponivel, error) {
	var horarios []*model.HorarioDisponivel
	query := r.Client.Collection("horariosDisponiveis").Where("psicologoId", "==", psicologoID)
	if status != "" {
		query = query.Where("status", "==", status)
	}

	iter := query.Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("erro ao listar horários: %w", err)
		}
		var h model.HorarioDisponivel
		if err := doc.DataTo(&h); err != nil {
			continue // Logar erro em um cenário real
		}
		h.ID = doc.Ref.ID
		horarios = append(horarios, &h)
	}
	return horarios, nil
}

func (r *HorarioDisponivelRepositoryImpl) BuscarHorarioPorID(ctx context.Context, id string) (*model.HorarioDisponivel, error) {
	doc, err := r.Client.Collection("horariosDisponiveis").Doc(id).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("horário com ID '%s' não encontrado", id)
	}
	var horario model.HorarioDisponivel
	if err := doc.DataTo(&horario); err != nil {
		return nil, err
	}
	horario.ID = doc.Ref.ID
	return &horario, nil
}

func (r *HorarioDisponivelRepositoryImpl) AtualizarStatusHorario(ctx context.Context, id string, novoStatus string) error {
	_, err := r.Client.Collection("horariosDisponiveis").Doc(id).Update(ctx, []firestore.Update{
		{Path: "status", Value: novoStatus},
	})
	return err
}

func (r *HorarioDisponivelRepositoryImpl) DeletarHorario(ctx context.Context, id string) error {
	_, err := r.Client.Collection("horariosDisponiveis").Doc(id).Delete(ctx)
	return err
}
