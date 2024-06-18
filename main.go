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

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db = fazConexaoComBanco()
var templates = template.Must(template.ParseGlob("listaPacientes.html"))

func main() {
	// Configuração do servidor para servir arquivos estáticos (HTML, CSS, JS, imagens, etc.)
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/listaPacientes", pacientes)
	http.HandleFunc("/cadastro", cadastroPacienteHandler)
	http.HandleFunc("/deletePaciente", deletePacienteHandler)
	http.HandleFunc("/getPaciente", getPacienteHandler)

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

func cadastroPacienteHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
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

    numero, err := strconv.ParseUint(r.FormValue("numero"), 10, 64)
    if err != nil {
        http.Error(w, "Número inválido", http.StatusBadRequest)
        log.Println("Número inválido:", err)
        return
    }
    
    dataCadastro := normalizeDate(r.FormValue("data_cadastro"))
    dataNascimento := normalizeDate(r.FormValue("data_nascimento"))

    paciente := Paciente{
        Id:             id,
        DataCadastro:   dataCadastro,
        Nome:           r.FormValue("nome"),
        NomeMae:        r.FormValue("nome_mae"),
        Cpf:            r.FormValue("cpf"),
        Sexo:           r.FormValue("sexo"),
        Email:          r.FormValue("email"),
        Telefone:       r.FormValue("telefone"),
        DataNascimento: dataNascimento,
        Cidade:         r.FormValue("cidade"),
        CEP:            r.FormValue("cep"),
        Rua:            r.FormValue("logradouro"),
        Numero:         numero,
    }

    log.Printf("Paciente recebido: %+v", paciente)

    if id > 0 {
        // Atualiza paciente existente
        _, err = db.Exec(`UPDATE paciente SET data_cadastro=$1, nome=$2, nome_da_mae=$3, cpf=$4, sexo=$5, email=$6, telefone_celular=$7, data_nascimento=$8, cidade=$9, cep=$10, rua=$11, num_casa=$12 WHERE id=$13`,
            paciente.DataCadastro, paciente.Nome, paciente.NomeMae, paciente.Cpf, paciente.Sexo, paciente.Email, paciente.Telefone, paciente.DataNascimento, paciente.Cidade, paciente.CEP, paciente.Rua, paciente.Numero, paciente.Id)
        if err != nil {
            http.Error(w, "Erro ao atualizar paciente", http.StatusInternalServerError)
            log.Println("Erro ao atualizar paciente:", err)
            return
        }
    } else {
        // Insere novo paciente
        log.Println(idStr)
        _, err = db.Exec(`INSERT INTO paciente (data_cadastro, nome, nome_da_mae, cpf, sexo, email, telefone_celular, data_nascimento, cidade, cep, rua, num_casa) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
            paciente.DataCadastro, paciente.Nome, paciente.NomeMae, paciente.Cpf, paciente.Sexo, paciente.Email, paciente.Telefone, paciente.DataNascimento, paciente.Cidade, paciente.CEP, paciente.Rua, paciente.Numero)
        if err != nil {
            http.Error(w, "Erro ao salvar paciente", http.StatusInternalServerError)
            log.Println("Erro ao salvar paciente:", err)
            return
        }
    }

    http.Redirect(w, r, "/listaPacientes", http.StatusSeeOther)
}

func deletePacienteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
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
	_, err = database.Query("CREATE TABLE IF NOT EXISTS paciente (id SERIAL PRIMARY KEY, data_cadastro varchar(10) NOT NULL, nome VARCHAR(255) NOT NULL, nome_da_mae varchar(255) NOT NULL,cpf VARCHAR(15) UNIQUE NOT NULL, sexo VARCHAR(10) NOT NULL, email varchar(255), telefone_celular VARCHAR(20), data_nascimento VARCHAR(12) NOT NULL, cidade varchar(255) NOT NULL, cep varchar(9) NOT NULL, rua varchar(255) NOT NULL, num_casa int)")
	if err != nil {
		log.Fatal(err)
	}

	return database
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
	return date.Format("2006-01-02")
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
	busca, err := db.Query(`SELECT * FROM paciente WHERE LOWER(nome) LIKE LOWER(concat(text($1), '%'))`, nome)
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
	err = db.QueryRow("SELECT id, data_cadastro, nome, nome_da_mae, cpf, sexo, email, telefone_celular, data_nascimento, cidade, cep, rua, num_casa FROM paciente WHERE id = $1", id).Scan(
		&paciente.Id, &paciente.DataCadastro, &paciente.Nome, &paciente.NomeMae, &paciente.Cpf, &paciente.Sexo, &paciente.Email, &paciente.Telefone, &paciente.DataNascimento, &paciente.Cidade, &paciente.CEP, &paciente.Rua, &paciente.Numero,
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

	log.Printf("Dados do paciente: %+v", paciente)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(paciente)
	if err != nil {
		log.Println("Erro ao codificar resposta JSON:", err)
	}
}

func alimentaBancoDeDados() {
	var Pacientes Pacientes

	//Erro ao salvar pacientelê o arquivo paciente.json e passa para o Go
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
