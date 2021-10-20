package main

import (
  "database/sql" 
  "next/http"
  "text/template"
  _ "github.com/go-sql-driver/mysql"
)
//struct utilizada para exibir 
type Names struct{
  Id int
  Name string
  Email string
} 

func dbConn() (db *sql.DB) { //abre conexão com o BD
	dbDriver := "mysql"
	dbUser := ""
	dbPass := ""
	dbName := ""

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

//renderiza o arquivo Index
func Index(w http.ResponseWriter, r *http.Request){

  //abre conexão com o banco de dados utilizando a função dbConn()
  db:=dbConn()

  //realiza a consulta no bd e trata erros
  selDB, err := db.Query("SELECT * FROM names ORDER BY id DESC")
  if err != nil{
    panic(err.Error())
  }

  //monta a scrut do template
  n:= Names{}

  //monta um array para guardar os valores da struct
  res:=[]Names{}

  //pega todos os valores do banco
  for selDB.Next(){
    var id int
    var name, email string

    //faz o scan do SELECT
    err = selDB.Scan(&id, &name, &email)
    if err != nil{
      panic(err.Error())
    }

    //envia os resultados para a struct
    n.Id = id
    n.Name = name
    n.Email = email

    //junta a scruct com o array
    res = append(res, n) 
  }

  //abre a página Index e exibe os registros na template
  tmpl.ExecuteTemplate(w, "Index", res)

  //fecha conexão
  defer db.Close()
}

//apenas exibe o resultado
func Show(w http.ResponseWriter, r*http.Request){
  db := dbConn()

  //pega o ID do paramentro da URLs
  nId := r.URL.Query().Get("id")

  //usa o ID para fazer a consulta e tratar erros
  selDB, err := db.Query("SELECT * FROM names WHERE id=?", nId)
  
  if err != nil{
      panic(err.Error())
  }

  //Monta a struct para ser utilizada no template
  n := Names{}

  //pega os valores do BD
  for selDB.next(){
    var id int
    var name, email string

    //faz o scan do SELECT
    err = selDB.Scan(&id,&name,&email)
    if err != nil{
      panic(err.Error())
    }

    //envia os resultados para a struct

    n.Id = id
    n.Name = name
    n.Email = email 
  } 
  //mostra o template
  tmpl.ExecuteTemplate(w, "Show", n)

  //fecha a conexão
  defer db.Close()

}

// exibe o formulário para inserir novos dados
func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// Função Edit, edita os dados
func Edit(w http.ResponseWriter, r *http.Request) {
	// Abre a conexão com banco de dados
	db := dbConn()

	// Pega o ID do parametro da URL
	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM names WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	// Monta a struct para ser utilizada no template
	n := Names{}

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for selDB.Next() {
		//Armazena os valores em variaveis
		var id int
		var name, email string

		// Faz o Scan do SELECT
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		n.Id = id
		n.Name = name
		n.Email = email
	}

	// Mostra o template com formulário preenchido para edição
	tmpl.ExecuteTemplate(w, "Edit", n)

	// Fecha a conexão com o banco de dados
	defer db.Close()
}

//insere os valores do bd

// Função Insert, insere valores no banco de dados
func Insert(w http.ResponseWriter, r *http.Request) {

	//Abre a conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	// Verifica o METHOD do fomrulário passado
	if r.Method == "POST" {

		// Pega os campos do formulário
		name := r.FormValue("name")
		email := r.FormValue("email")

		// Prepara a SQL e verifica errors
		insForm, err := db.Prepare("INSERT INTO names(name, email) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}

		// Insere valores do formulario com a SQL tratada e verifica errors
		insForm.Exec(name, email)

		// Exibe um log com os valores digitados no formulário
		log.Println("INSERT: Name: " + name + " | E-mail: " + email)
	}

	// Encerra a conexão do dbConn()
	defer db.Close()

	//Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

// Função Update, atualiza valores no banco de dados
func Update(w http.ResponseWriter, r *http.Request) {

	// Abre a conexão com o banco de dados usando a função: dbConn()
	db := dbConn()

	// Verifica o METHOD do formulário passado
	if r.Method == "POST" {

		// Pega os campos do formulário
		name := r.FormValue("name")
		email := r.FormValue("email")
		id := r.FormValue("uid")

		// Prepara a SQL e verifica errors
		insForm, err := db.Prepare("UPDATE names SET name=?, email=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		// Insere valores do formulário com a SQL tratada e verifica erros
		insForm.Exec(name, email, id)

		// Exibe um log com os valores digitados no formulario
		log.Println("UPDATE: Name: " + name + " |E-mail: " + email)
	}

	// Encerra a conexão do dbConn()
	defer db.Close()

	// Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

// Função Delete, deleta valores no banco de dados
func Delete(w http.ResponseWriter, r *http.Request) {

	// Abre conexão com banco de dados usando a função: dbConn()
	db := dbConn()

	nId := r.URL.Query().Get("id")

	// Prepara a SQL e verifica errors
	delForm, err := db.Prepare("DELETE FROM names WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	// Insere valores do form com a SQL tratada e verifica errors
	delForm.Exec(nId)

	// Exibe um log com os valores digitados no form
	log.Println("DELETE")

	// Encerra a conexão do dbConn()
	defer db.Close()

	// Retorna a HOME
	http.Redirect(w, r, "/", 301)
}

func main() {
  log.Println("Server started on: http://localhost:9000") // exibe mensagem que o servidor iniciou.

  //gerencia as URLs
  http.HandleFunc("/", Index)
  http.HandleFunc("/show", Show)
  http.HandleFunc("/new", New)
  http.HandleFunc("/edit", Edit)

  //Ações
  http.HandleFunc("/insert", Insert)
  http.HandleFunc("/update", Update)
  http.HandleFunc("/delete", Delete)

  //Inicia o servidor nas porta 9000
  http.ListenAndServe(":9000", nil)
}

//renderiza todos os templates da pasta tmpl" independente da extensão.
var tmpl=template.Must(template.ParseGlob("tmpl/*")) 