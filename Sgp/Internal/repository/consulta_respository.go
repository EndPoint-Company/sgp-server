package repository

import (
	"context"
	"fmt"
	"sgp/Internal/model"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type ConsultaRepositoryImpl struct {
	Client *firestore.Client
}

func NewConsultaRepository(client *firestore.Client) *ConsultaRepositoryImpl {
	return &ConsultaRepositoryImpl{Client: client}
}

func (r *ConsultaRepositoryImpl) AgendarConsulta(ctx context.Context, consulta model.Consulta) (*model.Consulta, error) {
	docRef, _, err := r.Client.Collection("Consultas").Add(ctx, map[string]interface{}{
		"alunoId":     consulta.AlunoID,
		"psicologoId": consulta.PsicologoID,
		"horario":     consulta.Horario,
		"status":      "aguardando aprovacao",
	})

	if err != nil {
		return nil, err
	}

	consulta.ID = docRef.ID
	consulta.Status = "aguardando aprovacao"
	return &consulta, nil
}

func (r *ConsultaRepositoryImpl) AtualizaStatusConsulta(ctx context.Context, id string, novoStatus string) error {
	_, err := r.Client.Collection("Consultas").Doc(id).Update(ctx, []firestore.Update{
		{Path: "status", Value: novoStatus},
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar status da consulta com ID '%s': %v", id, err)
	}
	return nil
}

func (r *ConsultaRepositoryImpl) ListarConsultasPorPsicologo(ctx context.Context, psicologoID string, statusFiltro string) ([]*model.Consulta, error) {
	var consultas []*model.Consulta

	query := r.Client.Collection("Consultas").Where("psicologoId", "==", psicologoID)
	if statusFiltro != "" {
		query = query.Where("status", "==", statusFiltro)
	}

	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("erro ao iterar sobre consultas do psicologo '%s': %v", psicologoID, err)
		}

		var consulta model.Consulta
		if err := doc.DataTo(&consulta); err != nil {
			fmt.Printf("erro ao converter dados do doc '%s': %v", doc.Ref.ID, err)
			continue
		}

		consulta.ID = doc.Ref.ID
		consultas = append(consultas, &consulta)
	}

	return consultas, nil
}

func (r *ConsultaRepositoryImpl) ListarConsultasPorAluno(ctx context.Context, alunoID string) ([]*model.Consulta, error) {
	var consultas []*model.Consulta

	iter := r.Client.Collection("Consultas").Where("alunoId", "==", alunoID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("erro ao iterar sobre consultas do aluno '%s': %v", alunoID, err)
		}

		var consulta model.Consulta
		if err := doc.DataTo(&consulta); err != nil {
			fmt.Printf("erro ao converter dados do doc '%s': %v", doc.Ref.ID, err)
			continue
		}

		consulta.ID = doc.Ref.ID
		consultas = append(consultas, &consulta)
	}

	return consultas, nil
}

func (r *ConsultaRepositoryImpl) DeletarConsulta(ctx context.Context, id string) error {
	_, err := r.Client.Collection("Consultas").Doc(id).Delete(ctx)
	if err != nil {
		return fmt.Errorf("erro ao deletar consulta com ID '%s': %v", id, err)
	}
	return nil
}
