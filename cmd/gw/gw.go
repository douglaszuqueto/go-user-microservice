package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/douglaszuqueto/go-grpc-user/proto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func main() {
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
	defer cancel()

	err = proto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:8001", opts)
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":8081", mux); err != nil {
		panic(err)
	}
}
