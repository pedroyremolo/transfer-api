# Transfer API

API criada com o propósito de prover o gerenciamento das transferências entre contas
de um banco digital.

A arquitetura da aplicação visa aplicar os conceitos do [Domain-Driver Design](https://www.amazon.com.br/Domain-Driven-Design-Eric-Evans/dp/8550800651),
Eric Evans, e da [Hexagonal Architecture](https://fideloper.com/hexagonal-architecture), de Alistair Cockburn.

## Requerimentos / Dependências
A aplicação, feita em [Go](https://golang.org/), depende do próprio módulo, e de pelo
menos uma instância [MongoDB](https://docs.mongodb.com/v4.2/).

A mesma é distribuída através de containers [Docker](https://docs.docker.com/).

Todas as dependências de pacotes estão relacionadas em go.mod, que é utilizado para 
gerenciamento das mesmas.

- [Go](https://golang.org/dl/) 1.15
- [BRDoc](https://github.com/Nhanderu/brdoc) 1.1.2
- [mongo-driver](https://github.com/mongodb/mongo-go-driver) 1.4.2
- [httprouter](https://github.com/julienschmidt/httprouter) 1.3.0
- [crypto](https://golang.org/x/crypto) 0.0.0-20201016220609-9e8e0b390897
- [jwt](https://github.com/gbrlsnchs/jwt) 3.0.0

Para baixa-las, com Go instalado na sua máquina:
```bash
$ go mod download
```

## Como usar

A aplicação possui distribuição via [Docker](Dockerfile), e possui um arquivo
[docker-compose](docker-compose.yml), sendo este o modo mais fácil de executa-la 
localmente.

Além disso, possui uma especificação [OpenAPI 3](https://swagger.io/specification/)
através do arquivo [openapi.yml](openapi.yml).

A mesma é gerenciada via variáveis de ambiente, segue abaixo a tabela:

| Nome                      | Descrição                                                  |
|---------------------------|------------------------------------------------------------|
| APP_PORT                  | Porta a ser escutada pela aplicação para novas requisições |
| APP_LOG_LEVEL             | Nível de log estruturado da aplicação                      |
| APP_DOCUMENT_DB_HOST      | Host da instância do MongoDB                               |
| APP_DOCUMENT_DB_PORT      | Porta da instância do MongoDB                              |
| APP_DOCUMENT_DB_USERNAME  | Usuário da instância do MongoDB                            |
| APP_DOCUMENT_DB_SECRET    | Senha da instância do MongoDB                              |
| APP_DOCUMENT_DB_NAME      | Nome do banco default da instância do MongoDB              |
| APP_JWT_GATEKEEPER_SECRET | Segredo de geração do token JWT                            |
| APP_JWT_GATEKEEPER_ISSUER | Emissor do token JWT                                       |

### Docker-Compose

Para executar via [docker-compose](https://docs.docker.com/compose/)

```bash
$ docker-compose up --build -d
```

## Licença
 A aplicação está sobe a licença [MIT](https://choosealicense.com/licenses/mit/)
 
## Créditos

Agradecimento em especial a todos os autores das bibliotecas de terceiros utilizadas, e
citadas acima.

E aos conteúdos, e seus autores, em que me baseei para a construção dessa aplicação:

### Videos

- [How Do You Structure Your Go Apps?](https://www.youtube.com/watch?v=1rxDzs0zgcE&t=2152s)
- [Building Hexagonal Microservices with Go](https://www.youtube.com/watch?v=rQnTtQZGpg8)
- [Unit testing HTTP servers](https://www.youtube.com/watch?v=hVFEV-ieeew)
- [The Context Package](https://www.youtube.com/watch?v=LSzR0VEraWw)


### Artigos

- [Aprenda Go com Testes](https://larien.gitbook.io/aprenda-go-com-testes/)
- [DDD Lite in Go](https://threedots.tech/post/ddd-lite-in-go-introduction/)
- [Repository Pattern in Go](https://threedots.tech/post/repository-pattern-in-go/)
- [Round Float to 2 decimal places](https://yourbasic.org/golang/round-float-2-decimal-places/)
- [Using Domain-Driven Design(DDD) in Golang](https://dev.to/stevensunflash/using-domain-driven-design-ddd-in-golang-3ee5)
- [A theory of modern Go](http://peter.bourgon.org/blog/2017/06/09/theory-of-modern-go.html)