Pitch de apresentação do projeto: https://www.youtube.com/watch?v=nFMU7ks_OvE

Apresentação do código SobreVidas ACS: https://youtu.be/p4hflOwnVWM

Como configurar o SobreVidas ACS para uso.


1. Baixando o projeto.
   
   1.1. Na página do repositório, clique em “Code” e em seguida, “Download ZIP”.
   
   1.2. Extraia o ZIP para a pasta SobreVidas-ACS-main.
2. Baixando o Go.

   2.1. No site do Go https://go.dev/dl/, faça download da linguagem de programação Go e instale-a.
3. Configurando o Banco de Dados.

   3.1. Baixe o SGBD PostgreSQL pelo link https://www.postgresql.org/download/ e instale-o, anotando a senha escolhida.

   3.2. No arquivo “.env” do projeto, insira seu usuário (por padrão postgres), sua senha definida na instalação do PostgreSQL e o nome do banco de dados (por padrão postgres).
  
   3.3. Pelo terminal ou cmd, execute o arquivo main.go para que a Tabela de Informações seja gerada no banco de dados utilizando:

       go run main.go
   
   3.4. Preencha a tabela "agente" com o login do agente. Para isso, acesse o banco de dados e clique em “Query Tool”. Ao acessar o query, digite o seguinte comando substituindo os nomes entre aspas pelos seus respectivos valores:

       INSERT INTO agente (nome, email, regiao, cpf, ine, cnes, senha)VALUES ('Nome do agente', 'email do Agente', 'Região atendida pelo Agente', 'CPF do Agente', 'ine do Agente', 'cnes do Agente', 'senha');  
4. Executando o SobreVidas ACS.

   4.1. Utilizando um compilador (cmd, terminal ou outro) e com a linguagem Go instalada, navegue até a pasta do projeto (SobreVidas-ACS-main) e utilize o comando “go run main.go” para iniciar a hospedagem do projeto localmente.

   4.2. Pelo seu navegador, acesse o site pela URL http://localhost:8052.

   4.3. Realize o login com as credenciais cadastradas.


