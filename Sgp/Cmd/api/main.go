package main

import (
	"context"
	"fmt"
	"time"

	"github.com/joho/godotenv"

	//"time"

	"log"
	"net/http"
	"os"

	//"sgp/Internal/model"
	"sgp/Internal/handler"
	"sgp/Internal/repository"

	firebase "firebase.google.com/go/v4"

	//"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	creds := os.Getenv("CREDS")

	//---------------Conexao com o firebase------------------//

	ctx := context.Background()
	opt := option.WithCredentialsFile(creds)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("erro ao inicializar firebase: %v", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("erro ao conectar ao firestore: %v", err)
	}
	defer client.Close()
	//-------------------------------------------------------//

	//--------------TestAluno----------------------//
	//alunoRepo := repository.NewAlunoRepository(client)
	//marquin := &model.Aluno{Nome: "Marcolas2", Email: "marcoladopcc@gmail.com"}
	//miguel := &model.Aluno{Nome: "Manguel", Email: "mangas@gmail.com"}
	//larinha := &model.Aluno{Nome: "Lorax", Email: "Lorax@gmail.com"}

	/*			funcionando
	alunoRepo.CriarAluno(ctx, *marquin)
	alunoRepo.CriarAluno(ctx, *miguel)
	alunoRepo.CriarAluno(ctx, *larinha)
	*/

	/*			funcionando
	alunoRepo.DeletarAluno(ctx,"RvVGDg6y8CifRlDsHib7")
	*/

	/*			funcionando
	jefso,err := alunoRepo.BuscarAlunoPorID(ctx,"2QZfL6ICZwO6UUMBN14I")
	fmt.Printf(jefso.Nome +"\n")
	fmt.Printf(jefso.Email +"\n")
	fmt.Printf(jefso.ID +"\n")
	*/

	/*			funcionando
	alunadas, err := alunoRepo.ListarAlunos(ctx)

	for i := 0; i < len(alunadas); i++ {
		fmt.Printf(alunadas[i].Nome + "\n")
		fmt.Printf(alunadas[i].Email + "\n")
		fmt.Printf(alunadas[i].ID + "\n")
		fmt.Printf("\n")
	}
	*/

	/*			funcionando
	marcos := &model.Aluno{Nome: "Marcos", Email: "marcos@gmail.com"}
	alunoRepo.AtualizarAluno(ctx,"2QZfL6ICZwO6UUMBN14I",*marcos)
	*/

	//alunodId, err := alunoRepo.GetAlunoIDPorNome(ctx, "Marcos")
	//fmt.Print(alunodId)

	//-------------------------------------------------------//

	//--------------TestPsico----------------------//
	//psicologoRepo := repository.NewPsicologoRepository(client)
	//freudiano := &model.Psicologo{Nome: "Valdemir", Email: "valdinhoo@gmail.com", CRP: "2111312"}

	//psicologoRepo.CriarPsicologo(ctx, *freudiano)

	//psicologoRepo.AtualizarPsicologo(ctx,"CRsiWje2vKiLsr5fCpxW",*freud)

	//abacate,err:= psicologoRepo.BuscarPsicologoPorID(ctx,"CRsiWje2vKiLsr5fCpxW")

	//psicos, err := psicologoRepo.ListarPsicologos(ctx)

	//psicologoRepo.DeletarPsicologo(ctx,"9bCrnEXS1DSM1YQk8DoB")

	//psicoId, err := psicologoRepo.GetPsicologoIDPorNome(ctx, "Freudinho")
	//----------------tudo funcionando---------------------//

	//--------------TestConsultas----------------------//
	//ConsultaRepo := repository.NewConsultaRepository(client)
	//Consulta01 := model.Consulta{AlunoID: "2QZfL6ICZwO6UUMBN14I", PsicologoID: "CRsiWje2vKiLsr5fCpxW", Horario: time.Now(), Status: "" }
	//Consulta02 := model.Consulta{AlunoID: "2QZfL6ICZwO6UUMBN14I", PsicologoID: "CRsiWje2vKiLsr5fCpxW", Horario: time.Now(), Status: "" }
	//Consulta03 := model.Consulta{AlunoID: "mu3Mo6I0eSSD3aWZdYte", PsicologoID: "CRsiWje2vKiLsr5fCpxW", Horario: time.Now(), Status: "" }
	//ConsultaRepo.AgendarConsulta(ctx, Consulta02)
	//ConsultaRepo.AgendarConsulta(ctx, Consulta03)
	//ConsultaRepo.AtualizaStatusConsulta(ctx,"xxOX58BhUZRSoNkAJ4dV", "aguardando aprovation" )
	//consultasP, err := ConsultaRepo.ListarConsultasPorPsicologo(ctx, "CRsiWje2vKiLsr5fCpxW", "aguardando aprovacao")
	//consultasA, err := ConsultaRepo.ListarConsultasPorAluno(ctx, "2QZfL6ICZwO6UUMBN14I")
	//----------------tudo funcionando---------------------//

	alunoRepo := repository.NewAlunoRepository(client)
	psicologoRepo := repository.NewPsicologoRepository(client)
	consultaRepo := repository.NewConsultaRepository(client)

	alunoHandler := handler.NewAlunoHandler(*alunoRepo)
	psicologoHandler := handler.NewPsicologoHandler(*psicologoRepo)
	consultaHandler := handler.NewConsultaHandler(*consultaRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /alunos", alunoHandler.HandlerCriarAluno)
	mux.HandleFunc("GET /alunos", alunoHandler.HandlerListarAlunos)
	mux.HandleFunc("GET /alunos/{id}", alunoHandler.HandlerBuscarAlunoPorID)
    mux.HandleFunc("GET /alunos/nome", alunoHandler.HandlerBuscarAlunoPorNome)
	mux.HandleFunc("PUT /alunos/{id}", alunoHandler.HandlerAtualizarAluno)
	mux.HandleFunc("DELETE /alunos/{id}", alunoHandler.HandlerDeletarAluno)

	mux.HandleFunc("POST /psicologos", psicologoHandler.HandlerCriarPsicologo)
	mux.HandleFunc("GET /psicologos", psicologoHandler.HandlerListarPsicologos)
    mux.HandleFunc("GET /psicologos/{id}", psicologoHandler.HandlerBuscarPsicologoPorID)
    mux.HandleFunc("GET /psicologos/nome", psicologoHandler.HandlerBuscarPsicologoPorNome)
	mux.HandleFunc("PUT /psicologos/{id}", psicologoHandler.HandlerAtualizarPsicologo)
	mux.HandleFunc("DELETE /psicologos/{id}", psicologoHandler.HandlerDeletarPsicologo)

	mux.HandleFunc("POST /consultas", consultaHandler.HandlerAgendarConsulta)
	mux.HandleFunc("GET /consultas/psicologo", consultaHandler.HandlerListarConsultasPorPsicologo)
	mux.HandleFunc("PATCH /consultas/{id}/status", consultaHandler.HandlerAtualizarStatusConsulta)
	mux.HandleFunc("DELETE /consultas/{id}", consultaHandler.HandlerDeletarConsulta)

	port := ":8080"
	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("🐄 bovino na porta %s\n", port)
	log.Fatal(server.ListenAndServe())
}
