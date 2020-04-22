package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/douglaszuqueto/go-grpc-user/pkg/grpc/api"
	"github.com/douglaszuqueto/go-grpc-user/pkg/grpc/server"
	"github.com/douglaszuqueto/go-grpc-user/pkg/storage"
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

	port := ":8001"
	rpcServer := server.NewServer(port)
	if rpcServer == nil {
		log.Println("Nao consigo escutar na porta:", port)
	}

	api.NewUserService(rpcServer.Grpc)

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
			ID:        idString,
			Username:  "username_" + idString,
			Email:     "username_" + idString + "@mail.com",
			State:     1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now().Add(time.Hour),
		}

		err := storage.CreateUser(user)
		if err != nil {
			log.Println("CreateUser err", err)
		}
	}

	<-doneCh
	rpcServer.Stop()
	log.Println("Finalizando...")
}
