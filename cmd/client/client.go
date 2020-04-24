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

var userService proto.UserServiceClient

func main() {
	fmt.Printf("\nGolang GRPC - Client\n\n")

	conn := connect()
	defer conn.Close()

	fmt.Printf("==> LIST <==\n\n")
	listUser()

	fmt.Printf("\n==> GET <==\n\n")
	getUser("1")

	user := &proto.User{
		Id:       "11",
		Username: "username_11",
		State:    1,
	}

	fmt.Printf("\n==> CREATE <==\n\n")

	createUser(user)

	fmt.Printf("\n==> UPDATE <==\n\n")

	user.State = 2

	updateUser(user)
	getUser(user.Id)

	fmt.Printf("\n==> DELETE <==\n\n")

	deleteUser(user.Id)
	getUser(user.Id)

	fmt.Println("\nFinish...")
}

func connect() *grpc.ClientConn {
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

	userService = proto.NewUserServiceClient(conn)

	return conn
}

func listUser() {
	req := &proto.ListUserRequest{}

	users, err := userService.List(context.Background(), req)
	if err != nil {
		fmt.Println("userService.List", err)
		return
	}

	for _, u := range users.User {
		fmt.Printf("ID: %v \t| username: %v  \t| state: %v\n", u.Id, u.Username, u.State)
	}
}

func getUser(id string) {
	req := &proto.GetUserRequest{
		Id: id,
	}

	user, err := userService.Get(context.Background(), req)
	if err != nil {
		fmt.Println("userService.List", err)
		return
	}

	fmt.Printf("ID: %v \t| username: %v  \t| state: %v\n", user.User.Id, user.User.Username, user.User.State)
}

func createUser(user *proto.User) {
	req := &proto.CreateUserRequest{
		User: user,
	}

	res, err := userService.Create(context.Background(), req)
	if err != nil {
		fmt.Println("userService.Create", err)
		return
	}

	fmt.Println("userService.Create:", res.Result)
}

func updateUser(user *proto.User) {
	req := &proto.UpdateUserRequest{
		User: user,
	}

	res, err := userService.Update(context.Background(), req)
	if err != nil {
		fmt.Println("userService.Update", err)
		return
	}

	fmt.Println("userService.Update:", res.Result)
}

func deleteUser(id string) {
	req := &proto.DeleteUserRequest{
		Id: id,
	}

	res, err := userService.Delete(context.Background(), req)
	if err != nil {
		fmt.Println("userService.Delete", err)
		return
	}

	fmt.Println("userService.Delete:", res.Result)
}
