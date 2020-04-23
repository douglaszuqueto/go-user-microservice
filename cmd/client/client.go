package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/douglaszuqueto/go-grpc-user/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var keep = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func main() {
	fmt.Printf("\nGolang GRPC - Client\n\n")

	grpcServerHost := os.Getenv("GRPC_SERVER_HOST")
	grpcServerPort := os.Getenv("GRPC_SERVER_PORT")

	uri := fmt.Sprintf("%s:%s", grpcServerHost, grpcServerPort)

	opts := grpc.WaitForReady(false)

	creds, err := credentials.NewClientTLSFromFile("./certs/server.crt", "")
	if err != nil {
		panic("could not load tls cert: %s" + err.Error())
	}

	conn, err := grpc.Dial(
		uri,
		// grpc.WithInsecure(),
		grpc.WithTransportCredentials(creds),
		grpc.WithKeepaliveParams(keep),
		grpc.WithDefaultCallOptions(opts),
	)

	if err != nil {
		panic("Error: " + err.Error())
	}

	userService := proto.NewUserServiceClient(conn)

	users, err := userService.List(context.Background(), &proto.ListUserRequest{})
	if err != nil {
		panic(err)
	}

	for _, u := range users.User {
		fmt.Printf("ID: %v \t| username: %v  \t| state: %v\n", u.Id, u.Username, u.State)
	}

	fmt.Println("\nFinish...")
	conn.Close()
}
