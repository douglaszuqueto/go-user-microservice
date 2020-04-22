package main

import (
	"fmt"
	"log"
	"time"

	"github.com/douglaszuqueto/go-grpc-user/pkg/storage"
)

func main() {
	user := storage.User{
		ID:        "1",
		Username:  "admin",
		Email:     "admin@mail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	fmt.Println("==> LIST <==")
	listUser()

	fmt.Println("==> CREATE <==")
	createUser(user)

	listUser()

	// update
	fmt.Println("==> UPDATE <==")
	user.State = 1

	updateUser(user)
	listUser()

	// get
	fmt.Println("==> GET <==")
	getUser(user.ID)

	// delete
	fmt.Println("==> DELETE <==")
	deleteUser(user.ID)

	listUser()
	getUser(user.ID)
}

func listUser() {
	users, err := storage.ListUser()
	if err != nil {
		log.Fatalln("UserList", err)
	}

	fmt.Println(users)
}

func getUser(id string) {
	user, err := storage.GetUser(id)
	if err != nil {
		log.Fatalln("UserGet", err)
	}

	fmt.Println(user)
}

func createUser(u storage.User) {
	err := storage.CreateUser(u)
	if err != nil {
		log.Fatalln("UserCreate", err)
	}
}

func updateUser(u storage.User) {
	err := storage.UpdateUser(u)
	if err != nil {
		log.Fatalln("UserUpdate", err)
	}
}

func deleteUser(id string) {
	err := storage.DeleteUser(id)
	if err != nil {
		log.Fatalln("UserDelete", err)
	}
}
