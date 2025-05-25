package main

import (
	"Desktop/golangProjects/CRUD/pkg/client"
	"flag"
	"log"
)

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
}

func main() {
	var user user

	flag.StringVar(&user.Name, "name", "", "this is the name of the user")
	flag.IntVar(&user.Age, "age", -1, "this is the age of the user")
	flag.BoolVar(&user.Create, "create", false, "this is the operation that creates a user")
	flag.BoolVar(&user.Create, "read", false, "this is the operation that reads a user's data")
	flag.BoolVar(&user.Create, "update", false, "this is the operation that updates a user's data")
	flag.BoolVar(&user.Create, "delete", false, "this is the operation that deletes a user")

	flag.Parse()

	if !validateBooleanFlags(user) {
		log.Panic("can only call one of the following at a time: create, read, update, or delete")
	}

	if user.Name == "" {
		log.Panic("no name was inputted")
	}
	client := client.New("")
	if user.Age == -1 && (user.Create || user.Update) {
		log.Panic("must input age (non-negative) as well")
	}
	name := user.Name
	age := user.Age
	if user.Create {
		client.CreateUser(name, age)
	} else if user.Read {
		client.ReadUser(name)
	} else if user.Update {
		client.UpdateUser(name, age)
	} else if user.Delete {
		client.DeleteUser(name)
	} else {
		log.Panicf("no options (create, read, update, or delete) were selected for this script")
	}
}
