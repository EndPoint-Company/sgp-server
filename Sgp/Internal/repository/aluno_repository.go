package repository

import (
	"context"
	"fmt"
	"sgp/Internal/model"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)


type AlunoRepositoryImpl struct {
	Client *firestore.Client
}

func NewAlunoRepository(client *firestore.Client) *AlunoRepositoryImpl {
	return &AlunoRepositoryImpl{Client: client}
}

func (r *AlunoRepositoryImpl) CriarAluno(ctx context.Context, aluno model.Aluno) (*model.Aluno, error) {
	docRef, _, err := r.Client.Collection("Alunos").Add(ctx, map[string]interface{}{
		"nome":  aluno.Nome,
		"email": aluno.Email,
	})

	if err != nil {
		return nil, err
	}

	aluno.ID = docRef.ID
	return &aluno, nil
}

func (r *AlunoRepositoryImpl) BuscarAlunoPorID(ctx context.Context, id string) (*model.Aluno, error) {
	doc, err := r.Client.Collection("Alunos").Doc(id).Get(ctx)

	if err != nil {
		return nil, err
	}
	var aluno model.Aluno

	if err := doc.DataTo(&aluno); err != nil {
		return nil, err
	}
	aluno.ID = doc.Ref.ID
	return &aluno, nil
}

func (r *AlunoRepositoryImpl) ListarAlunos(ctx context.Context) ([]*model.Aluno, error) {
	var alunos []*model.Aluno

	iter := r.Client.Collection("Alunos").Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("erro ao iterar sobre alunos: %v", err)
		}

		var aluno model.Aluno

		if err := doc.DataTo(&aluno); err != nil {
			fmt.Printf("erro ao converter dados do doc '%s': %v", doc.Ref.ID, err)
			continue
		}

		aluno.ID = doc.Ref.ID
		alunos = append(alunos, &aluno)
	}
	return alunos, nil
}

func (r *AlunoRepositoryImpl) AtualizarAluno(ctx context.Context, id string, aluno model.Aluno) error {
	_, err := r.Client.Collection("Alunos").Doc(id).Set(ctx, map[string]interface{}{
		"nome":  aluno.Nome,
		"email": aluno.Email,
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar o aluno com ID '%s': %v", id, err)
	}
	return nil
}

func (r *AlunoRepositoryImpl) DeletarAluno(ctx context.Context, id string) error {
	_, err := r.Client.Collection("Alunos").Doc(id).Delete(ctx)
	if err != nil {
		//return fmt.Errorf("erro ao deletar aluno com ID '%s': %v")
	}
	return nil
}

func (r *AlunoRepositoryImpl) GetAlunoIDPorNome(ctx context.Context, nome string) (string, error) {
	query := r.Client.Collection("Alunos").Where("nome", "==", nome).Limit(1)

	iter := query.Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return "", fmt.Errorf("aluno com o nome '%s' não encontrado", nome)
	}
	if err != nil {
		return "", fmt.Errorf("erro ao buscar aluno por nome: %w", err)
	}

	return doc.Ref.ID, nil
}
