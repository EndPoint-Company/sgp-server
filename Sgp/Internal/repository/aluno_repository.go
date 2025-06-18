package repository

import (
	"context"
	"log"
	"sgp/Internal/model"
	"google.golang.org/api/iterator"
	"cloud.google.com/go/firestore"
)

type AlunoRepository struct{
	Client *firestore.Client
}

func NewAlunoRepository(client *firestore.Client) *AlunoRepository{
	return &AlunoRepository{Client: client}
}

func (r *AlunoRepository) CriarAluno(ctx context.Context, aluno model.Aluno)(*model.Aluno, error){
	docRef, _, err := r.Client.Collection("alunos").Add(ctx, map[string]interface{}{
		"nome": aluno.Nome,
		"email": aluno.Email,
	})

	if err != nil{
		return nil, err
	}

	aluno.ID = docRef.ID
	return &aluno, nil
}

func(r *AlunoRepository) BuscarAlunoPorID(ctx context.Context, id string)(*model.Aluno,error ){
	doc, err := r.Client.Collection("alunos").Doc(id).Get(ctx)

	if err != nil{
		return nil, err
	}
	var aluno model.Aluno

	if err := doc.DataTo(&aluno); err != nil{
		return nil, err
	} 
	aluno.ID = doc.Ref.ID
	return &aluno, nil
}

func (r *AlunoRepository) ListarAlunos(ctx context.Context)([]*model.Aluno, error){
	var alunos []*model.Aluno

	iter := r.Client.Collection("alunos").Documents(ctx)

	for{
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil{
			log.Printf("erro ao iterar sobre alunos: %v", err)	
			return nil, err		
		}

		var aluno model.Aluno

		if err := doc.DataTo(&aluno); err != nil{
			log.Printf("erro ao converter dados do doc '%s': %v", doc.Ref.ID, err)	
			continue
		}

		aluno.ID = doc.Ref.ID
		alunos = append(alunos, &aluno)
	}
	return alunos, nil
}

func (r *AlunoRepository) AtualizarAluno(ctx context.Context, id string, aluno model.Aluno)error{
	_, err := r.Client.Collection("alunos").Doc(id).Set(ctx, map[string]interface{}{
		"nome": aluno.Nome,
		"email": aluno.Email,
	})
	if err != nil {
		log.Printf("erro ao atualizar o alno com ID '%s': %v", id, err)
		return err
	}
	return nil
}

func(r *AlunoRepository) DeletarAluno(ctx context.Context, id string) error {
	_, err := r.Client.Collection("alunos").Doc(id).Delete(ctx)
	if err != err {
		log.Printf("erro ao deletar aluno com ID '%s': %v")
		return err
	}
	return nil
}