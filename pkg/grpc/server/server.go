package grpc

import (
	"Desktop/golangProjects/CRUD/pkg/database"
	pb "Desktop/golangProjects/CRUD/pkg/grpc/proto"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
)

var (
	ErrInCreate         = errors.New("error in writing to the database")
	ErrInRead           = errors.New("error in reading from the database")
	ErrInUpdate         = errors.New("error in updating the database")
	ErrInDelete         = errors.New("error in deleting from the database")
	errMessageAgeOrName = "age or name was not specified"
	errMessageName      = "name was not specified"
	CreateMessage       = "user %s age %d has been successfully created"
	ReadMessage         = "user %s's age is %s"
	UpdateMessage       = "user %s's new age is %d"
	DeleteMessage       = "deleting user %s from the database"
	UserNotFoundMessage = "user %s not found in our database"
)

type Server struct {
	db database.Database
	pb.CRUDServer
}

func New(db database.Database) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Create(ctx context.Context, in *pb.UserWriteReq) (*pb.DatabaseResp, error) {
	name := in.GetName()
	age := in.GetAge()
	// TODO change age to something that can be nil
	if name == "" || in.Age == nil {
		return &pb.DatabaseResp{
			Success:    false,
			Message:    "",
			ErrMessage: errMessageAgeOrName,
		}, nil
	}
	err := s.db.Write(name, strconv.Itoa(int(age)))
	if err != nil {
		log.Printf("error in writing to the database: %v", err)
		return &pb.DatabaseResp{
			Success:    false,
			Message:    "",
			ErrMessage: ErrInCreate.Error(),
		}, nil
	}
	log.Printf("succesfully create user %s age %d", name, age)
	return &pb.DatabaseResp{Success: true, Message: fmt.Sprintf(CreateMessage, name, age), ErrMessage: ""}, nil
}

func (s *Server) Read(ctx context.Context, in *pb.UserReadReq) (*pb.DatabaseResp, error) {
	name := in.GetName()
	if name == "" {
		return &pb.DatabaseResp{
			Success:    false,
			Message:    "",
			ErrMessage: errMessageName,
		}, nil
	}
	age, err := s.db.Read(name)
	if err == database.ErrUserNotFound {
		return &pb.DatabaseResp{Success: false, Message: "", ErrMessage: fmt.Sprintf(UserNotFoundMessage, name)}, nil

	}
	if err != nil {
		log.Printf("error in reading from the database: %v", err)
		return &pb.DatabaseResp{
			Success:    false,
			Message:    "",
			ErrMessage: ErrInRead.Error(),
		}, nil
	}
	log.Printf("user %s is age %s", name, age)
	return &pb.DatabaseResp{Success: true, Message: fmt.Sprintf(ReadMessage, name, age), ErrMessage: ""}, nil
}

func (s *Server) Update(ctx context.Context, in *pb.UserWriteReq) (*pb.DatabaseResp, error) {
	name := in.GetName()
	age := in.GetAge()
	if name == "" || in.Age == nil {
		return &pb.DatabaseResp{
			Success:    false,
			Message:    "",
			ErrMessage: errMessageAgeOrName,
		}, nil
	}
	err := s.db.Write(name, strconv.Itoa(int(age)))
	if err != nil {
		log.Printf("error in updating the database: %v", err)
		return &pb.DatabaseResp{
			Success:    false,
			Message:    "",
			ErrMessage: ErrInUpdate.Error(),
		}, nil
	}
	log.Printf("succesfully updated user %s's age to %d", name, age)
	return &pb.DatabaseResp{Success: true, Message: fmt.Sprintf(UpdateMessage, name, age), ErrMessage: ""}, nil
}

func (s *Server) Delete(ctx context.Context, in *pb.UserReadReq) (*pb.DatabaseResp, error) {
	name := in.GetName()
	if name == "" {
		return &pb.DatabaseResp{
			Success:    false,
			Message:    "",
			ErrMessage: errMessageName,
		}, nil
	}
	err := s.db.Delete(name)
	if err != nil {
		log.Printf("error in deleting from the database: %v", err)
		return &pb.DatabaseResp{
			Success:    false,
			Message:    "",
			ErrMessage: ErrInDelete.Error(),
		}, nil
	}
	log.Printf("sucessfully deleted user %s from the database", name)
	return &pb.DatabaseResp{Success: true, Message: fmt.Sprintf(DeleteMessage, name), ErrMessage: ""}, nil
}
