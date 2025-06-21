package main

import (
	"context"
	//"fmt"
	//"time"

	"log"
	//"net/http"
	//"sgp/Internal/model"
	//"sgp/Internal/model"
	//"sgp/Internal/repository"

	firebase "firebase.google.com/go/v4"

	//"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

func main() {
	//---------------Conexao com o firebase------------------//
	ctx := context.Background()
	opt := option.WithCredentialsFile("/home/marco/sgp-server/key.json")
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

	/*			funcionando
	alunodId, err := alunoRepo.GetAlunoIDPorNome(ctx, "Marcos")
	fmt.Print(alunodId)
	*/
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
}
