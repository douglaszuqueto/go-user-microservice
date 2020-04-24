package main

import (
	"fmt"
	"log"
	"time"

	"github.com/douglaszuqueto/go-grpc-user/pkg/storage"
	"github.com/douglaszuqueto/go-grpc-user/pkg/util"
)

var db storage.UserStorage

func main() {
	user := storage.User{
		ID:        "1",
		Username:  "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db = util.GetStorageType()

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
	users, err := db.ListUser()
	if err != nil {
		log.Fatalln("UserList", err)
	}

	fmt.Println(users)
}

func getUser(id string) {
	user, err := db.GetUser(id)
	if err != nil {
		log.Fatalln("UserGet", err)
	}

	fmt.Println(user)
}

func createUser(u storage.User) {
	err := db.CreateUser(u)
	if err != nil {
		log.Fatalln("UserCreate", err)
	}
}

func updateUser(u storage.User) {
	err := db.UpdateUser(u)
	if err != nil {
		log.Fatalln("UserUpdate", err)
	}
}

func deleteUser(id string) {
	err := db.DeleteUser(id)
	if err != nil {
		log.Fatalln("UserDelete", err)
	}
}
