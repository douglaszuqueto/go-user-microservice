package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/douglaszuqueto/go-grpc-user/pkg/grpc/api"
	"github.com/douglaszuqueto/go-grpc-user/pkg/grpc/server"
	"github.com/douglaszuqueto/go-grpc-user/pkg/storage"

	_ "github.com/lib/pq"
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

	var db storage.UserStorage = storage.GetStorageType()

	grpcServerHost := os.Getenv("GRPC_SERVER_HOST")
	grpcServerPort := os.Getenv("GRPC_SERVER_PORT")

	uri := fmt.Sprintf("%s:%s", grpcServerHost, grpcServerPort)

	rpcServer := server.NewServer(uri)

	api.NewUserService(rpcServer.Grpc, db)

	go func() {
		err := rpcServer.Start()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	for i := 1; i <= 10; i++ {
		idString := strconv.Itoa(i)

		log.Println("Inserindo user:", idString)

		user := storage.User{
			Username:  "username_" + idString,
			State:     1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now().Add(time.Hour),
		}

		err := db.CreateUser(user)
		if err != nil {
			log.Println("CreateUser err", err)
		}
	}

	<-doneCh
	rpcServer.Stop()
	log.Println("Finalizando...")
}
