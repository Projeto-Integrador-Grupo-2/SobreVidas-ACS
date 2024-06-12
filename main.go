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
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db = fazConexaoComBanco()
var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	// Configuração do servidor para servir arquivos estáticos (HTML, CSS, JS, imagens, etc.)
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/listaPacientes", pacientes)

	alimentaBancoDeDados()

	log.Println("Server rodando na porta 8052")
	// Inicia o servidor na porta 8052
	err := http.ListenAndServe(":8052", nil)
	if err != nil {
		log.Fatal(err)
	}
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
	_, err = database.Query("CREATE TABLE IF NOT EXISTS paciente (id SERIAL PRIMARY KEY, data_cadastro varchar(10) NOT NULL, nome VARCHAR(255) NOT NULL, nome_da_mae varchar(255) NOT NULL,cpf VARCHAR(15) UNIQUE NOT NULL, sexo VARCHAR(10) NOT NULL, email varchar(255), telefone_celular VARCHAR(20), data_nascimento VARCHAR(12) NOT NULL, cidade varchar(255) NOT NULL, cep varchar(9) NOT NULL, rua varchar(255) NOT NULL, num_casa int)")
	if err != nil {
		log.Fatal(err)
	}

	return database
}

func cadastraPaciente(paciente Paciente) {
	// insere paciente no banco de dados
	_, err := db.Exec(`INSERT INTO paciente (data_cadastro, nome, nome_da_mae, cpf, sexo, email, telefone_celular, data_nascimento, cidade, cep, rua, num_casa) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) on conflict do nothing`, paciente.DataCadastro, paciente.Nome, paciente.NomeMae, paciente.Cpf, paciente.Sexo, paciente.Email, paciente.Telefone, paciente.DataNascimento, paciente.Cidade, paciente.CEP, paciente.Rua, paciente.Numero)
	if err != nil {
		fmt.Println(err)
	}
}

func buscaPacientePorNome(nome string) Pacientes {
	// retorna pacientes por nome
	busca, err := db.Query(`SELECT * FROM paciente WHERE LOWER(nome) LIKE LOWER(concat('%', text($1), '%'))`, nome)
	if err != nil {
		fmt.Println(err)
	}

	var pacientes Pacientes

	// Realiza a estrutura de repetição pegando todos os valores do banco
	for busca.Next() {

		var paciente Paciente

		// Armazena os valores em variáveis
		var Id, Num_casa uint64
		var Data, Nome, Nome_mae, Cpf, Sexo, email, Telefone, DataNascimento, Cidade, CEP, Rua string

		// Faz o Scan do SELECT
		err = busca.Scan(&Id, &Data, &Nome, &Nome_mae, &Cpf, &Sexo, &email, &Telefone, &DataNascimento, &Cidade, &CEP, &Rua, &Num_casa)
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
		paciente.Rua = Rua
		paciente.Numero = Num_casa

		// Junta a Struct com Array
		pacientes.Pacientes = append(pacientes.Pacientes, paciente)
	}

	return pacientes
}

func alimentaBancoDeDados() {
	var Pacientes Pacientes

	// lê o arquivo paciente.json e passa para o Go
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
	Id             uint64
	DataCadastro   string `json:"Data_cad"`
	Nome           string `json:"nome"`
	NomeMae        string `json:"Nome_mae"`
	Cpf            string `json:"cpf"`
	Sexo           string `json:"sexo"`
	Email          string `json:"email"`
	Telefone       string `json:"celular"`
	DataNascimento string `json:"data_nasc"`
	Cidade         string `json:"Cidade"`
	CEP            string `json:"CEP"`
	Rua            string `json:"Rua"`
	Numero         uint64 `json:"Num_casa"`
}

type Pacientes struct {
	Pacientes []Paciente `json:"pacientes"`
}
