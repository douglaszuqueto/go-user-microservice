package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/douglaszuqueto/go-grpc-user/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	grpcServerHost = os.Getenv("GRPC_SERVER_HOST")
	grpcServerPort = os.Getenv("GRPC_SERVER_PORT")

	grpcGatewayHost = os.Getenv("GRPC_GW_HOST")
	grpcGatewayPort = os.Getenv("GRPC_GW_PORT")
)

func main() {
	signalCh := make(chan os.Signal, 1)
	doneCh := make(chan bool, 1)

	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-signalCh
		log.Println("Signal stop:", sig)

		doneCh <- true
	}()

	creds, err := credentials.NewClientTLSFromFile("./certs/server.crt", "")
	if err != nil {
		log.Fatalln("could not load tls cert:", err)
	}

	k := keepalive.ClientParameters{
		Time:                5 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,     // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,            // send pings even without active streams
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithKeepaliveParams(k),
	}

	mux := runtime.NewServeMux()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	grpcURI := fmt.Sprintf("%s:%s", grpcServerHost, grpcServerPort)

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

	<-doneCh
	cancel()
	log.Println("Finalizando...")
}
