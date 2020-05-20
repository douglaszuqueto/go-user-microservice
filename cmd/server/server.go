package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/douglaszuqueto/go-grpc-user/pkg/grpc/api"
	"github.com/douglaszuqueto/go-grpc-user/pkg/grpc/server"
	"github.com/douglaszuqueto/go-grpc-user/pkg/storage"
	"github.com/douglaszuqueto/go-grpc-user/pkg/util"
	"github.com/douglaszuqueto/go-grpc-user/pkg/util/graceful"

	_ "github.com/lib/pq"
)

var (
	grpcServerHost = os.Getenv("GRPC_SERVER_HOST")
	grpcServerPort = os.Getenv("GRPC_SERVER_PORT")
)

var db storage.UserStorage

func main() {
	grace := graceful.New()

	db = storage.GetStorageType()

	if storageType := os.Getenv("APP_STORAGE"); storageType == "memory" {
		insertData()
	}

	uri := fmt.Sprintf("%s:%s", grpcServerHost, grpcServerPort)

	rpcServer := server.NewServer(uri)

	api.NewUserService(rpcServer.Grpc, db)

	go func() {
		err := rpcServer.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	grace.Wait()
	rpcServer.Stop()

	log.Println("Finalizando...")
}

func insertData() {
	for i := 1; i <= 10; i++ {
		idString := strconv.Itoa(i)

		log.Println("Inserindo user:", idString)

		password, _ := util.GeneratePassword("password_" + idString)

		user := storage.User{
			Username:  "username_" + idString,
			Password:  password,
			State:     1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now().Add(time.Hour),
		}

		_, err := db.CreateUser(user)
		if err != nil {
			log.Println("CreateUser err", err)
		}
	}
}
