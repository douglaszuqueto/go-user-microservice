# Serviço de usuário utilizando Go e GRPC

## Serviços

* Server
* Gateway
* CLI

## Requisitos

* [Go](https://golang.org/dl/)
* [Protobuf](https://github.com/protocolbuffers/protobuf)
* [gRPC-Go](https://github.com/grpc/grpc-go)
* [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)

## Storage

* Memória
* PostgreSQL

## Métodos

* list
* get
* create
* update
* delete

## Deploy
-
-

### Standalone
-
-

### Docker
-

### Docker compose
-

### Kubernetes
-

## API

* Endpoints

| Version | HTTP | Path |
|:--:|:--:|:--|
| /v1 | GET | /user | 
| /v1 | GET | /user/1 | 
| /v1 | POST | /user |
| /v1 | PUT | /user/1 | 
| /v1 | DELETE | /user/1 |

### All
```bash
curl --request GET \
  --url http://127.0.0.1:8081/v1/user
```

### Get
```bash
curl --request GET \
  --url http://127.0.0.1:8081/v1/user/1
```

### Create
```bash
curl --request POST \
  --url http://127.0.0.1:8081/v1/user \
  --header 'content-type: application/json' \
  --data '{
	"user": {
		"username": "admin",
		"state": 2
	}
}'
```

### Update
```bash
curl --request PUT \
  --url http://127.0.0.1:8081/v1/user/1 \
  --header 'content-type: application/json' \
  --data '{
	"user": {
		"id": "1",
		"username": "admin",
		"state": 5
	}
}'
```

### Delete
```bash
curl --request DELETE \
  --url http://127.0.0.1:8081/v1/user/1
```

## Geração de certificados de segurança

```bash
openssl genrsa -out server.key
openssl req -new -sha256 -key server.key -out server.csr
openssl x509 -req -days 3650 -in server.csr -out server.crt -signkey server.key
```

## Geração do JWT Secret
```bash
openssl rand -base64 64
```

## Todo

* Status codes nos erros
* Validação de contexto - timeout, interrupt...
* JWT Manager
* Interceptors
  * Server
    * unary
    * stream
  * Client
    * unary
    * stream
* CLI
* Testes - client, server

## Changelog

## Referências

https://dev.to/plutov/docker-and-go-modules-3kkn
