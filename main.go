package main

import (
	"Desktop/golangProjects/CRUD/pkg"
	"Desktop/golangProjects/CRUD/pkg/server/database"
	"Desktop/golangProjects/CRUD/pkg/server/grpc"
	server "Desktop/golangProjects/CRUD/pkg/server/http"
	pb "Desktop/golangProjects/CRUD/proto"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	rpc "google.golang.org/grpc"
)

func grpcMain() {
	db, err := database.New("/Users/alexpaley/Desktop/golangProjects/CRUD/pkg/server/database/database.db", 0666, nil)
	if err != nil {
		log.Fatalf("unable to start database: %v", err)
	}
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("grpc server listening at %v", lis.Addr())
	server := grpc.New(db)
	s := rpc.NewServer()
	pb.RegisterCRUDServer(s, server)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("unable to serve due to err: %v", err)
	}

}

func httpMain() {
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

func main() {
	// httpMain()
	grpcMain()
}
