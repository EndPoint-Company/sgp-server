package repository

import (
	"context"
	"log"
	"sgp/Internal/model"
	"google.golang.org/api/iterator"
	"cloud.google.com/go/firestore"
)

type PsicologoRepository struct {
	Client *firestore.Client
}

func NewPsicologoRepository(client *firestore.Client) *PsicologoRepository {
	return &PsicologoRepository{ Client: client }
}

func (r *PsicologoRepository) CriarPsicologo(
	ctx context.Context, psicologo model.Psicologo)(*model.Psicologo, error) {

	docRef, _, err := r.Client.Collection("Psicologos").
		Add(ctx, map[string] interface{}{
			"nome": psicologo.Nome,
			"email": psicologo.Email,
			"crp": psicologo.CRP,
		});

	if err != nil{
		return nil, err
	}

	psicologo.ID = docRef.ID

	return &psicologo, nil
}

func(r *PsicologoRepository) BuscarPsicologoPorID(ctx context.Context, id string)(*model.Psicologo,error ){
	doc, err := r.Client.Collection("Psicologos").Doc(id).Get(ctx)

	if err != nil{
		return nil, err
	}
	var Psicologo model.Psicologo

	if err := doc.DataTo(&Psicologo); err != nil{
		return nil, err
	} 
	Psicologo.ID = doc.Ref.ID
	return &Psicologo, nil
}

func (r *PsicologoRepository) ListarPsicologos(ctx context.Context)([]*model.Psicologo, error){
	var Psicologos []*model.Psicologo

	iter := r.Client.Collection("Psicologos").Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("erro ao iterar sobre psicologos: %v", err)	

			return nil, err		
		}

		var Psicologo model.Psicologo

		if err := doc.DataTo(&Psicologo); err != nil {
			log.Printf("erro ao converter dados do doc '%s': %v", 
				doc.Ref.ID, err)	
			
			continue
		}

		Psicologo.ID = doc.Ref.ID
		Psicologos = append(Psicologos, &Psicologo)
	}

	return Psicologos, nil
}

func (r *PsicologoRepository) AtualizarPsicologo(
	ctx context.Context, id string, Psicologo model.Psicologo) error {

	_, err := r.Client.Collection("Psicologos").
		Doc(id).
		Set(ctx, map[string] interface{} {
			"nome": Psicologo.Nome,
			"email": Psicologo.Email,
			"crp": Psicologo.CRP,
		})

	if err != nil {
		log.Printf("erro ao atualizar o psicologo com ID '%s': %v", id, err)
		return err
	}
	return nil
}

func(r *PsicologoRepository) DeletarPsicologo(
	ctx context.Context, id string) error {

	_, err := r.Client.Collection("Psicologos").Doc(id).Delete(ctx)

	if err != nil {
		log.Printf("erro ao deletar psicologo com ID '%s': %v")
		return err
	}

	return nil
}