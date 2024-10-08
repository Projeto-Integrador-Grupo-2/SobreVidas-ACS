package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db = fazConexaoComBanco()
var templates = template.Must(template.ParseGlob("*.html"))
var store = sessions.NewCookieStore([]byte("super-secret-key"))

func init() {
	store.Options = &sessions.Options{
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   1800,
		HttpOnly: true,
	}
}

func main() {
	// Configuração do servidor para servir arquivos estáticos (HTML, CSS, JS, imagens, etc.)
	fs := http.FileServer(http.Dir("./"))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		fs.ServeHTTP(w, r)
	}))

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	alimentaBancoDeDados()

	registerProtectedRoutes()

	log.Println("Server rodando na porta 8052")
	// Inicia o servidor na porta 8052
	err := http.ListenAndServe(":8052", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func perfilPacienteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	paciente := Paciente{}
	query := `SELECT data_cadastro, nome, nome_da_mae, cpf, sexo, data_nascimento, email, telefone_celular, bebe, fuma, possui_feridas_boca, cidade, cep, bairro, rua, num_casa FROM paciente WHERE id=$1`
	row := db.QueryRow(query, id)
	err := row.Scan(&paciente.DataCadastro, &paciente.Nome, &paciente.NomeMae, &paciente.Cpf, &paciente.Sexo, &paciente.DataNascimento, &paciente.Email, &paciente.Telefone, &paciente.Bebe, &paciente.Fuma, &paciente.PossuiFeridasBoca, &paciente.Cidade, &paciente.CEP, &paciente.Bairro, &paciente.Rua, &paciente.Numero)
	data_cad := strings.Split(paciente.DataCadastro, "/")
	paciente.DataCadastro = data_cad[2] + "/" + data_cad[1] + "/" + data_cad[0]
	data_nasc := strings.Split(paciente.DataNascimento, "/")
	paciente.DataNascimento = data_nasc[2] + "/" + data_nasc[1] + "/" + data_nasc[0]
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("perfil_paciente.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, paciente)
}

func pacientes(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		return
	}
	busca := strings.TrimSpace(r.Form.Get("busca"))

	Pacientes := buscaPacientePorNome(busca)

	templates.ExecuteTemplate(w, "listaPacientes.html", Pacientes)
}

func cadastroPacienteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.ServeFile(w, r, "cadastro.html")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao processar o formulário", http.StatusBadRequest)
		log.Println("Erro ao processar o formulário:", err)
		return
	}

	idStr := r.FormValue("id")
	var id uint64
	if idStr != "" {
		id, err = strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			log.Println("ID inválido:", err)
			return
		}
	}

	numero, _ := strconv.ParseUint(r.FormValue("numero"), 10, 64)

	dataCadastro := normalizeDate(r.FormValue("data_cadastro"))
	dataNascimento := normalizeDate(r.FormValue("data_nascimento"))
	bebe := r.FormValue("bebe") == "on"
	fuma := r.FormValue("fuma") == "on"
	possui_feridas_boca := r.FormValue("possui_feridas_boca") == "on"

	paciente := Paciente{
		Id:                id,
		DataCadastro:      dataCadastro,
		Nome:              r.FormValue("nome"),
		NomeMae:           r.FormValue("nome_mae"),
		Cpf:               r.FormValue("cpf"),
		Sexo:              r.FormValue("sexo"),
		Email:             r.FormValue("email"),
		Telefone:          r.FormValue("telefone"),
		DataNascimento:    dataNascimento,
		Cidade:            r.FormValue("cidade"),
		CEP:               r.FormValue("cep"),
		Bairro:            r.FormValue("bairro"),
		Rua:               r.FormValue("logradouro"),
		Numero:            numero,
		Bebe:              bebe,
		Fuma:              fuma,
		PossuiFeridasBoca: possui_feridas_boca,
	}

	if id > 0 {
		// Atualiza paciente existente
		_, err = db.Exec(`UPDATE paciente SET data_cadastro=$1, nome=$2, nome_da_mae=$3, cpf=$4, sexo=$5, email=$6, telefone_celular=$7, data_nascimento=$8, cidade=$9, cep=$10, bairro=$11, rua=$12, num_casa=$13, bebe=$14, fuma=$15, possui_feridas_boca=$16 WHERE id=$17`,
			paciente.DataCadastro, paciente.Nome, paciente.NomeMae, paciente.Cpf, paciente.Sexo, paciente.Email, paciente.Telefone, paciente.DataNascimento, paciente.Cidade, paciente.CEP, paciente.Bairro, paciente.Rua, paciente.Numero, paciente.Bebe, paciente.Fuma, paciente.PossuiFeridasBoca, paciente.Id)
		if err != nil {
			http.Error(w, "Erro ao atualizar paciente", http.StatusInternalServerError)
			log.Println("Erro ao atualizar paciente:", err)
			return
		}
	} else {
		// Insere novo paciente
		_, err = db.Exec(`INSERT INTO paciente (data_cadastro, nome, nome_da_mae, cpf, sexo, email, telefone_celular, data_nascimento, cidade, cep, bairro, rua, num_casa, bebe, fuma, possui_feridas_boca) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
			paciente.DataCadastro, paciente.Nome, paciente.NomeMae, paciente.Cpf, paciente.Sexo, paciente.Email, paciente.Telefone, paciente.DataNascimento, paciente.Cidade, paciente.CEP, paciente.Bairro, paciente.Rua, paciente.Numero, paciente.Bebe, paciente.Fuma, paciente.PossuiFeridasBoca)
		if err != nil {
			http.Error(w, "Erro ao salvar paciente", http.StatusInternalServerError)
			log.Println("Erro ao salvar paciente:", err)
			return
		}
		_, err = db.Exec(`UPDATE graphs SET novos_pacientes = novos_pacientes + 1`)
		if err != nil {
			http.Error(w, "Erro ao somar", http.StatusInternalServerError)
			log.Println("Erro ao somar:", err)
			return
		}
	}

	http.Redirect(w, r, "/listaPacientes", http.StatusSeeOther)
}

func deletePacienteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/listaPacientes", http.StatusSeeOther)
		return
	}

	var requestData struct {
		ID uint64 `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Erro ao processar a solicitação", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM paciente WHERE id = $1", requestData.ID)
	if err != nil {
		http.Error(w, "Erro ao excluir paciente", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func fazConexaoComBanco() *sql.DB {
	// carrega arquivo .env com dados de conexão com o banco
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo .env")
	}

	// faz a busca dos atributos no arquivo .env para usa-las na conexão com banco
	usuarioBancoDeDados := os.Getenv("USUARIO")
	senhaDoUsuario := os.Getenv("SENHA")
	nomeDoBancoDeDados := os.Getenv("NOME_BANCO_DE_DADOS")
	dadosParaConexao := "user=" + usuarioBancoDeDados + " dbname=" + nomeDoBancoDeDados + " password=" + senhaDoUsuario + " host=localhost port=5432 sslmode=disable"
	database, err := sql.Open("postgres", dadosParaConexao)
	if err != nil {
		fmt.Println(err)
	}

	// cria tabela paciente com atributos como: id, nome, cpf, data de nascimento, telefone, sexo e booleanos referente a situação fisica
	_, err = database.Query("CREATE TABLE IF NOT EXISTS paciente (id SERIAL PRIMARY KEY, data_cadastro varchar(10) NOT NULL, nome VARCHAR(255) NOT NULL, nome_da_mae varchar(255) NOT NULL,cpf VARCHAR(15) UNIQUE NOT NULL, sexo VARCHAR(10) NOT NULL, email varchar(255), telefone_celular VARCHAR(20), data_nascimento VARCHAR(12) NOT NULL, cidade varchar(255) NOT NULL, cep varchar(9) NOT NULL, bairro varchar(255) NOT NULL, rua varchar(255) NOT NULL, num_casa int, bebe bool, fuma bool, possui_feridas_boca bool)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Query(`CREATE TABLE IF NOT EXISTS agente (
		id SERIAL PRIMARY KEY,
		nome VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		regiao VARCHAR(100) NOT NULL,
		cpf VARCHAR(15) UNIQUE NOT NULL,
		ine VARCHAR(20) NOT NULL,
		cnes VARCHAR(20) NOT NULL,
		senha VARCHAR(255) NOT NULL
		)`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = database.Query(`CREATE TABLE IF NOT EXISTS graphs (
		novos_pacientes INT,
		pacientes_cadastrados INT
		)`)

	if err != nil {
		log.Fatal(err)
	}

	return database
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	if session.Values["authenticated"] == true {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		agente := buscaAgentePorEmailESenha(email, password)
		if agente != nil {

			session, _ := store.Get(r, "session-name")
			session.Values["authenticated"] = true
			session.Values["agenteID"] = agente.ID
			session.Save(r, w)

			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		log.Println("Credenciais inválidas")
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, "Erro ao renderizar template", http.StatusInternalServerError)
		log.Println("Erro ao renderizar template:", err)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	delete(session.Values, "agenteID")
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func buscaAgentePorEmailESenha(email, password string) *Agente {

	row := db.QueryRow("SELECT id, nome, email, regiao, cpf, ine, cnes FROM agente WHERE email = $1 AND senha = $2", email, password)

	var agente Agente
	err := row.Scan(&agente.ID, &agente.Nome, &agente.Email, &agente.Regiao, &agente.CPF, &agente.INE, &agente.CNES)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Nenhum agente encontrado com as credenciais fornecidas")
			return nil
		}
		log.Println("Erro ao buscar agente:", err)
		return nil
	}
	return &agente
}

func perfilHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	agenteID, ok := session.Values["agenteID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	agente := buscaAgentePorID(uint64(agenteID))
	if agente == nil {
		http.Error(w, "Agente não encontrado", http.StatusNotFound)
		return
	}

	templates.ExecuteTemplate(w, "perfil.html", agente)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		authenticated, ok := session.Values["authenticated"].(bool)
		if !ok || !authenticated {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func registerProtectedRoutes() {
	protectedRoutes := map[string]http.HandlerFunc{
		"/listaPacientes.html":  listaPacientesHandler,
		"/map_view.html":        map_viewHandler,
		"/graphs":               graphsHandler,
		"/perfil_paciente.html": perfilpacienteHandler,
		"/home":                 homepageHandler,
		"/perfil":               perfilHandler,
		"/perfil.html":          perfilhtmlHandler,
		"/listaPacientes":       pacientes,
		"/cadastro":             cadastroPacienteHandler,
		"/deletePaciente":       deletePacienteHandler,
		"/getPaciente":          getPacienteHandler,
		"/mapa":                 mapHandler,
		"/perfil_paciente":      perfilPacienteHandler,
	}

	for route, handler := range protectedRoutes {
		http.Handle(route, authMiddleware(http.HandlerFunc(handler)))
	}
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "home_page.html", nil)
}

func listaPacientesHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "listaPacientes.html", nil)
}

func map_viewHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "map_view.html", nil)
}

func perfilpacienteHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "perfil_paciente.html", nil)
}

func perfilhtmlHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "perfil.html", nil)
}

func buscaAgentePorID(id uint64) *Agente {

	row := db.QueryRow("SELECT id, nome, email, regiao, cpf, ine, cnes FROM agente WHERE id = $1", id)

	var agente Agente
	err := row.Scan(&agente.ID, &agente.Nome, &agente.Email, &agente.Regiao, &agente.CPF, &agente.INE, &agente.CNES)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Nenhum agente encontrado com o ID fornecido")
			return nil
		}
		log.Println("Erro ao buscar agente:", err)
		return nil
	}
	return &agente
}

type Agente struct {
	ID     int
	Nome   string
	Email  string
	Regiao string
	CPF    string
	INE    string
	CNES   string
	Senha  string
}

func normalizeDate(dateStr string) string {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date, err = time.Parse("02-01-2006", dateStr)
		if err != nil {
			log.Println("Erro ao parsear data:", err)
			return dateStr
		}
	}
	return date.Format("2006/01/02")
}

func cadastraPaciente(paciente Paciente) {
	// insere paciente no banco de dados
	_, err := db.Exec(`INSERT INTO paciente (data_cadastro, nome, nome_da_mae, cpf, sexo, email, telefone_celular, data_nascimento, cidade, cep, bairro, rua, num_casa, bebe, fuma, possui_feridas_boca) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) on conflict do nothing`, paciente.DataCadastro, paciente.Nome, paciente.NomeMae, paciente.Cpf, paciente.Sexo, paciente.Email, paciente.Telefone, paciente.DataNascimento, paciente.Cidade, paciente.CEP, paciente.Bairro, paciente.Rua, paciente.Numero, paciente.Bebe, paciente.Fuma, paciente.PossuiFeridasBoca)
	if err != nil {
		fmt.Println(err)
	}
}

func buscaPacientePorNome(nome string) Pacientes {
	// retorna pacientes por nome
	busca, err := db.Query(`SELECT * FROM paciente WHERE LOWER(nome) LIKE LOWER(concat(text($1), '%'))`, nome)
	if err != nil {
		fmt.Println(err)
	}

	var pacientes Pacientes

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for busca.Next() {

		var paciente Paciente

		// Armazena os valores em variáveis
		var Bebe, Fuma, Possui_feridas_boca bool
		var Id, Num_casa uint64
		var Data, Nome, Nome_mae, Cpf, Sexo, email, Telefone, DataNascimento, Cidade, CEP, Bairro, Rua string

		// Faz o Scan do SELECT
		err = busca.Scan(&Id, &Data, &Nome, &Nome_mae, &Cpf, &Sexo, &email, &Telefone, &DataNascimento, &Cidade, &CEP, &Bairro, &Rua, &Num_casa, &Bebe, &Fuma, &Possui_feridas_boca)
		if err != nil {
			panic(err.Error())
		}

		// Envia os resultados para a struct
		paciente.Id = Id
		paciente.DataCadastro = Data
		paciente.Nome = Nome
		paciente.NomeMae = Nome_mae
		paciente.Cpf = Cpf
		paciente.Sexo = Sexo
		paciente.DataNascimento = DataNascimento
		paciente.Email = email
		paciente.Telefone = Telefone
		paciente.Cidade = Cidade
		paciente.CEP = CEP
		paciente.Bairro = Bairro
		paciente.Rua = Rua
		paciente.Numero = Num_casa
		paciente.Bebe = Bebe
		paciente.Fuma = Fuma
		paciente.PossuiFeridasBoca = Possui_feridas_boca

		// Junta a Struct com Array
		pacientes.Pacientes = append(pacientes.Pacientes, paciente)
	}

	return pacientes
}

func getPacienteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID não fornecido", http.StatusBadRequest)
		log.Println("ID não fornecido")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		log.Println("ID inválido:", err)
		return
	}

	var paciente Paciente
	err = db.QueryRow("SELECT id, data_cadastro, nome, nome_da_mae, cpf, sexo, email, telefone_celular, data_nascimento, cidade, cep, bairro, rua, num_casa, bebe, fuma, possui_feridas_boca FROM paciente WHERE id = $1", id).Scan(
		&paciente.Id, &paciente.DataCadastro, &paciente.Nome, &paciente.NomeMae, &paciente.Cpf, &paciente.Sexo, &paciente.Email, &paciente.Telefone, &paciente.DataNascimento, &paciente.Cidade, &paciente.CEP, &paciente.Bairro, &paciente.Rua, &paciente.Numero, &paciente.Bebe, &paciente.Fuma, &paciente.PossuiFeridasBoca,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Paciente não encontrado", http.StatusNotFound)
			log.Println("Paciente não encontrado")
		} else {
			http.Error(w, "Erro ao buscar paciente", http.StatusInternalServerError)
			log.Println("Erro ao buscar paciente:", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(paciente)
	if err != nil {
		log.Println("Erro ao codificar resposta JSON:", err)
	}
}

func alimentaBancoDeDados() {
	var Pacientes Pacientes

	jsonFile, _ := os.Open("paciente.json")
	byteJson, _ := io.ReadAll(jsonFile)

	err := json.Unmarshal(byteJson, &Pacientes)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(Pacientes.Pacientes); i++ {
		cadastraPaciente(Pacientes.Pacientes[i])
	}
}

type Paciente struct {
	Id                uint64
	DataCadastro      string `json:"Data_cad"`
	Nome              string `json:"nome"`
	NomeMae           string `json:"Nome_mae"`
	Cpf               string `json:"cpf"`
	Sexo              string `json:"sexo"`
	Email             string `json:"email"`
	Telefone          string `json:"celular"`
	DataNascimento    string `json:"data_nasc"`
	Cidade            string `json:"Cidade"`
	CEP               string `json:"CEP"`
	Bairro            string `json:"Bairro"`
	Rua               string `json:"Rua"`
	Numero            uint64 `json:"Num_casa"`
	Fuma              bool   `json:"Fuma"`
	Bebe              bool   `json:"Bebe"`
	PossuiFeridasBoca bool   `json:"Possui_feridas_boca"`
}

type Pacientes struct {
	Pacientes []Paciente `json:"pacientes"`
}

type Endereco struct {
	Nome   string `json:"Nome"`
	Rua    string `json:"Rua"`
	Numero string `json:"Numero"`
	Bairro string `json:"Bairro"`
	Cidade string `json:"Cidade"`
	CEP    string `json:"CEP"`
}

type Enderecos struct {
	Enderecos []Endereco `json:"Enderecos"`
}

func mapHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.URL.Query().Get("json") == "true" {
		busca, err := db.Query("SELECT nome, cidade, num_casa, cep, rua, bairro FROM paciente")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer busca.Close()

		var enderecos Enderecos

		for busca.Next() {
			var endereco Endereco
			err = busca.Scan(&endereco.Nome, &endereco.Cidade, &endereco.Numero, &endereco.CEP, &endereco.Rua, &endereco.Bairro)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			enderecos.Enderecos = append(enderecos.Enderecos, endereco)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(enderecos)
		return
	}

	err := templates.ExecuteTemplate(w, "map_view.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func graphsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT novos_pacientes FROM graphs")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	var graph Graph
	for rows.Next() {
		err := rows.Scan(&graph.NovosPacientes)
		if err != nil {
			http.Error(w, "Erro ao ler dados", http.StatusInternalServerError)
			return
		}
	}

	rows, err = db.Query("SELECT count(*) AS exact_count FROM paciente;")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&graph.PacientesCadastrados)
		if err != nil {
			http.Error(w, "Erro ao ler dados", http.StatusInternalServerError)
			return
		}
	}

	tmpl, err := template.ParseFiles("graphs.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, graph)
}

type Graph struct {
	NovosPacientes       int
	PacientesCadastrados int
}
