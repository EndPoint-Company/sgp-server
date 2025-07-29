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

type Consulta struct {
	ID              string    `json:"id" firestore:"-"`
	AlunoID         string    `json:"alunoId" firestore:"alunoId"`
	PsicologoID     string    `json:"psicologoId" firestore:"psicologoId"`
	HorarioID       string    `json:"horarioId" firestore:"horarioId"` 
	Inicio          time.Time `json:"inicio" firestore:"inicio"`
	Fim             time.Time `json:"fim" firestore:"fim"`
	Status          string    `json:"status" firestore:"status"`
	DataAgendamento time.Time `json:"dataAgendamento" firestore:"dataAgendamento"`
}

type HorarioDisponivel struct {
	ID          string    `json:"id" firestore:"-"`
	PsicologoID string    `json:"psicologoId" firestore:"psicologoId"`
	Inicio      time.Time `json:"inicio" firestore:"inicio"`
	Fim         time.Time `json:"fim" firestore:"fim"`
	Status      string    `json:"status" firestore:"status"`
}