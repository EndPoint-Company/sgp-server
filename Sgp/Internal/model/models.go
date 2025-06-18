package model

import "time"

type Aluno struct {
	ID string `json:"id" firestore:"-"`
	Nome string `json:"nome" firestore:"nome"`
	Email string `json:"email" firestore:"email"`
}

type Psicologo struct{
	ID string `json:"id" firestore:"-"`
	Nome string `json:"nome" firestore:"nome"`
	Email string `json:"email" firestore:"email"`
	CRP string `json:"crp" firestore:"crp"`
}

type Consulta struct{
	ID string `json:"id" firestore:"-"`
	AlunoID string `json:"alunoId" firestore:"alunoId"`
	PsicologoID string `json:"psicologoId" firestore:"psicologoId"`
	Horario time.Time `json:"horario" firestore:"horario"`
	Status string `json:"status" firestore:"status"`
}