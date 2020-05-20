package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/douglaszuqueto/go-grpc-user/pkg/util/graceful"
	"github.com/douglaszuqueto/go-grpc-user/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	grpcServerHost = os.Getenv("GRPC_SERVER_HOST")
	grpcServerPort = os.Getenv("GRPC_SERVER_PORT")

	grpcGatewayHost = os.Getenv("GRPC_GW_HOST")
	grpcGatewayPort = os.Getenv("GRPC_GW_PORT")
)

func main() {
	grace := graceful.New()

	creds, err := credentials.NewClientTLSFromFile("./certs/server.crt", "")
	if err != nil {
		log.Fatalln("could not load tls cert:", err)
	}

	mux := runtime.NewServeMux()

	ctx, cancelFunc := context.WithCancel(context.Background())

	grpcURI := fmt.Sprintf("%s:%s", grpcServerHost, grpcServerPort)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	err = proto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcURI, opts)
	if err != nil {
		panic(err)
	}

	grpcGatewayURI := fmt.Sprintf("%s:%s", grpcGatewayHost, grpcGatewayPort)

	go func() {
		if err := http.ListenAndServe(grpcGatewayURI, mux); err != nil {
			panic(err)
		}
	}()

	grace.Wait()
	cancelFunc()

	log.Println("Finalizando...")
}
