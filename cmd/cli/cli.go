package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/douglaszuqueto/go-grpc-user/proto"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	grpcServerHost = os.Getenv("GRPC_SERVER_HOST")
	grpcServerPort = os.Getenv("GRPC_SERVER_PORT")
)

var userService proto.UserServiceClient

func main() {
	conn := connect()
	defer conn.Close()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "user",
				Aliases:     []string{"u"},
				Usage:       "options for user",
				Subcommands: userCommands(),
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func userCommands() []*cli.Command {
	commands := []*cli.Command{
		{
			Name:  "list",
			Usage: "list a users",
			Action: func(c *cli.Context) error {
				listUser()
				return nil
			},
		},
		{
			Name:  "get",
			Usage: "get a user",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "id", Aliases: []string{"i"}},
			},
			Action: func(c *cli.Context) error {
				getUser(c.String("id"))
				return nil
			},
		},
		{
			Name:  "create",
			Usage: "get a user",
			Action: func(c *cli.Context) error {
				fmt.Println("create a user")

				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Username: ")
				text, _ := reader.ReadString('\n')
				fmt.Println(text)

				fmt.Print("Password: ")
				text, _ = reader.ReadString('\n')
				fmt.Println(text)

				return nil
			},
		},
		{
			Name:  "update",
			Usage: "update a user",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "id", Aliases: []string{"i"}},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("update a user: ", c.String("id"))
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove an existing",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "id", Aliases: []string{"i"}},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("removed user: ", c.String("id"))
				return nil
			},
		},
	}

	return commands
}

func connect() *grpc.ClientConn {
	uri := fmt.Sprintf("%s:%s", grpcServerHost, grpcServerPort)

	creds, err := credentials.NewClientTLSFromFile("./certs/server.crt", "")
	if err != nil {
		panic("could not load tls cert: %s" + err.Error())
	}

	options := []grpc.DialOption{
		// grpc.WithInsecure(),
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(false),
		),
	}

	conn, err := grpc.Dial(
		uri,
		options...,
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

func createUser(user *proto.User) (string, error) {
	req := &proto.CreateUserRequest{
		User: user,
	}

	res, err := userService.Create(context.Background(), req)
	if err != nil {
		fmt.Println("userService.Create", err)
		return "", err
	}

	fmt.Println("userService.Create:", res.Id)

	return res.Id, nil
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
