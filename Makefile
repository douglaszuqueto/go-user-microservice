include .env

.EXPORT_ALL_VARIABLES:

GRPC_GW_PATH := $(shell go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)
APIS_PATH := "$(GRPC_GW_PATH)/../third_party/googleapis"

dev-server:
	GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info go run -race cmd/server/server.go

dev-client:
	go run -race cmd/client/client.go

dev-gw:
	go run -race cmd/gw/gw.go

prod:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/grpc-server cmd/server/server.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/grpc-client cmd/client/client.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/grpc-gw cmd/gw/gw.go

	upx bin/grpc-server
	upx bin/grpc-client
	upx bin/grpc-gw

test:
	go test -race -cover ./...
	
update:
	go get ./...
	go mod tidy

pb:
	@protoc -I ${APIS_PATH} -I proto/ proto/*.proto --go_out=plugins=grpc:proto
	@protoc -I ${APIS_PATH} -I proto/ proto/*.proto --grpc-gateway_out=logtostderr=true:./proto

docker-build:
	./docker.sh server
	./docker.sh gw

docker-compose:
	docker-compose up

.PHONY: dev-server dev-client dev-gw prod test update pb docker-build docker-compose
