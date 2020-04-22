include .env

.EXPORT_ALL_VARIABLES:

GRPC_GW_PATH := $(shell go list -f '{{ .Dir }}' github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway)
APIS_PATH := "$(GRPC_GW_PATH)/../third_party/googleapis"

dev-server:
	GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info go run cmd/server/server.go

dev-client:
	go run cmd/client/client.go

dev-gw:
	go run cmd/gw/gw.go

prod:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/grpc-server cmd/server/server.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/grpc-client cmd/client/client.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/grpc-gw cmd/gw/gw.go

	upx bin/grpc-server
	upx bin/grpc-client
	upx bin/grpc-gw

update:
	go get all
	go mod tidy

pb:
	@protoc -I ${APIS_PATH} -I proto/ proto/*.proto --go_out=plugins=grpc:proto
	@protoc -I ${APIS_PATH} -I proto/ proto/*.proto --grpc-gateway_out=logtostderr=true:./proto
