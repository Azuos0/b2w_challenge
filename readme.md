# B2W Digital - Star Wars Challenge

## Sobre o desafio

Requisitos:

- A API deve ser REST

- Para cada planeta, os seguintes dados devem ser obtidos do banco de dados da aplicação, sendo inserido manualmente:

      Nome
      Clima
      Terreno


- Para cada planeta também devemos ter a quantidade de aparições em filmes, que podem ser obtidas pela API pública do Star Wars: https://swapi.dev/about

Funcionalidades desejadas: 

- Adicionar um planeta (com nome, clima e terreno)

- Listar planetas

- Buscar por nome

- Buscar por ID

- Remover planeta

OBS: A linguagem para realização do desafio será correspondente a do anúncio da vaga.
Bancos que usamos: MongoDB, Cassandra, DynamoDB, Datomic, ELK.

E lembre-se! Um bom software é um software bem testado.

## Rodando o sistema com Docker

Caso você possua o docker instalado na sua máquina basta apenas clonar esse repositório e criar um arquivo .env na raiz do projeto, copiando as informações do arquivo .env.example (apresentadas abaixo) para o .env.

```docker
#.env.example
PORT=":8000"
MONGODB_URL="mongodb://mongodb:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false"
MONGODB_DATABASE="swapp"
MONGODB_TEST_DATABASE="swapp_test"
```

Feito isso, abra um terminal na raiz do projeto e digite o comando:
```docker
docker-compose up -d 
```

E está feito! Você está pronto para adicionar planetas à sua galáxia!!! 

### Rodando os testes dentro do container

Com os containers do sistema rodando, abra um terminal e digite

```docker
docker exec -it swapp /bin/sh      # comando para entrar no container

go test -v ./...                   # roda todos os testes automatizados após entrar no container
```

## Rodando sem o docker
Para rodar o sistema sem o docker é necessário que você possua o golang instalado em sua máquina (caso não o tenha, veja [aqui](https://golang.org/doc/install) como fazer) e tenha acesso ao mongoDB (seja pelo Atlas Database ou qualquer outra forma).

Feito isso, você precisa criar o arquivo .env na raiz do projeto com os dados do seu mongoDB 

```docker
#.env
PORT=":8000"
MONGODB_URL=            #Aqui vai a url do seu cluster
MONGODB_DATABASE=       #seu banco de dados
MONGODB_TEST_DATABASE=  #o banco de dados que será utilizado para os testes automatizados
```

Abrir um terminal na raiz do projeto e baixar as dependências de desenvolvimento e rodar sua aplicação

```docker
go mod download       # baixa as dependências
go run app/main.go    # roda a aplicação
```

Para rodar todos os testes da aplicação é necessário ter um terminal aberto na raiz do projeto e rodar o comando
```docker
go test -v ./...
```

## Endpoints

Clique [aqui](https://app.swaggerhub.com/apis-docs/Azuos0/b-2_w_star_wars/1.0.0) para ver os Endpoints pelo swagger

- localhost:8000/api/   
  - Method: GET | Mensagem de boas-vindas
- localhost:8000/api/planet 
  - Method: POST | Adiciona um novo planeta
  - Request body:
    - name: string - obrigatório
    - climate: string - obrigatório
    - terrain: string - obrigatório
- localhost:8000/api/planet/:id
  - Method: GET | busca um determinado planeta pelo id
- localhost:8000/api/planet/:id
  - Method: DELETE | deleta um determinado planeta pelo id
- localhost:8000/api/planets
  - Method: GET | lista e procura por planetas
  - Query params:
    - name: nome do planeta
    - page: página da lista 


# May the force be with you!