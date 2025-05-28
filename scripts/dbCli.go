package main

import (
	grpClient "Desktop/golangProjects/CRUD/pkg/grpc/client"
	httpClient "Desktop/golangProjects/CRUD/pkg/http/client"
	"flag"
	"fmt"
	"log"

	grpc "google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	CreateUser(name string, age int) (string, error)
	ReadUser(name string) (string, error)
	UpdateUser(name string, age int) (string, error)
	DeleteUser(name string) (string, error)
}

var _ Client = &httpClient.Client{}
var _ Client = &grpClient.Client{}

func validateBooleanFlags(u user) bool {
	counter := 0
	if u.Create {
		counter += 1
	}
	if u.Read {
		counter += 1
	}
	if u.Update {
		counter += 1
	}
	if u.Delete {
		counter += 1
	}

	return counter == 1
}

type user struct {
	Name   string
	Age    int
	Create bool
	Read   bool
	Update bool
	Delete bool
	GRPC   bool
}

func main() {
	var user user

	flag.StringVar(&user.Name, "name", "", "this is the name of the user")
	flag.IntVar(&user.Age, "age", -1, "this is the age of the user")
	flag.BoolVar(&user.Create, "create", false, "this is the operation that creates a user")
	flag.BoolVar(&user.Read, "read", false, "this is the operation that reads a user's data")
	flag.BoolVar(&user.Update, "update", false, "this is the operation that updates a user's data")
	flag.BoolVar(&user.Delete, "delete", false, "this is the operation that deletes a user")
	flag.BoolVar(&user.GRPC, "grpc", false, "this uses grpc if true and http if false")

	flag.Parse()

	if !validateBooleanFlags(user) {
		log.Panic("can only call one of the following at a time: create, read, update, or delete")
	}

	if user.Name == "" {
		log.Panic("no name was inputted")
	}
	var c Client
	if user.GRPC {
		conn, err := grpc.NewClient(
			"localhost:8080", // Replace with your server address
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("error starting a grpc connection: %v", err)
		}
		c, err = grpClient.New(conn)
		if err != nil {
			log.Fatalf("error starting a grpc client: %v", err)
		}
	} else {
		c = httpClient.New("")
	}
	if user.Age == -1 && (user.Create || user.Update) {
		log.Panic("must input age (non-negative) as well")
	}
	name := user.Name
	age := user.Age
	if user.Create {
		resp, err := c.CreateUser(name, age)
		if err != nil {
			log.Panic("error in creating a user")
		}
		fmt.Println(resp)
	} else if user.Read {
		resp, err := c.ReadUser(name)
		if err != nil {
			log.Panicf("error in reading the user: %s: %v", name, err)
		}
		fmt.Println(resp)
	} else if user.Update {
		resp, err := c.UpdateUser(name, age)
		if err != nil {
			log.Panicf("error in updating the user: %s: %v", name, err)
		}
		fmt.Println(resp)
	} else if user.Delete {
		resp, err := c.DeleteUser(name)
		if err != nil {
			log.Panicf("error in creating the user: %s: %v", name, err)
		}
		fmt.Println(resp)
	} else {
		log.Panicf("no options (create, read, update, or delete) were selected for this script")
	}
}
