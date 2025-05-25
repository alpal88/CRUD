package main

import (
	"Desktop/golangProjects/CRUD/pkg"
	server "Desktop/golangProjects/CRUD/pkg/server"
	"Desktop/golangProjects/CRUD/pkg/server/database"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	db, err := database.New("/Users/alexpaley/Desktop/golangProjects/CRUD/pkg/server/database/database.db", 0666, nil)
	if err != nil {
		log.Panicf("unable to start database: %v", err)
	}
	serv := server.New(db)
	mux := mux.NewRouter()
	mux.HandleFunc(pkg.CREATEADDRROUTE, serv.HandleCreate)
	mux.HandleFunc(fmt.Sprintf("%s%s", pkg.USERADDROUTE, "{name}"), serv.HandleUsers)

	addr := ":8080"

	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("the server is up and running")
	log.Fatal(s.ListenAndServe())
}
