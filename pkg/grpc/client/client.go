package grpc

import (
	pb "Desktop/golangProjects/CRUD/pkg/grpc/proto"
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc"
)

var (
	noNameOrAgeMsg = "need a name and an age"
	noNameMsg      = "need a name"
)

type Client struct {
	client pb.CRUDClient
}

func New(conn *grpc.ClientConn) (*Client, error) {
	client := pb.NewCRUDClient(conn)
	return &Client{
		client: client,
	}, nil
}

func (c *Client) Create(name string, age int) (string, error) {
	if name == "" || age == -1 {
		return noNameOrAgeMsg, errors.New("no age or name")
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	int32Age := int32(age)
	req := &pb.UserWriteReq{
		Name: name,
		Age:  &int32Age,
	}
	resp, err := c.client.Create(ctx, req)
	if err != nil {
		log.Printf("error during user create: %v", err)
		return "", err
	}
	if resp.Success {
		return resp.Message, nil
	}
	return resp.ErrMessage, nil
}

func (c *Client) Read(name string) (string, error) {
	if name == "" {
		return noNameMsg, errors.New("no name")
	}
	req := &pb.UserReadReq{
		Name: name,
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	resp, err := c.client.Read(ctx, req)
	if err != nil {
		log.Printf("error during read: %v", err)
		return "", err
	}
	if resp.Success {
		return resp.Message, nil
	}
	return resp.ErrMessage, nil
}

func (c *Client) Update(name string, age int) (string, error) {
	if name == "" || age == -1 {
		return noNameOrAgeMsg, errors.New("no age or name")
	}
	int32Age := int32(age)
	req := &pb.UserWriteReq{
		Name: name,
		Age:  &int32Age,
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	resp, err := c.client.Update(ctx, req)
	if err != nil {
		log.Printf("error during user update: %v", err)
		return "", err
	}
	if resp.Success {
		return resp.Message, nil
	}
	return resp.ErrMessage, nil
}

func (c *Client) Delete(name string) (string, error) {
	if name == "" {
		return noNameMsg, errors.New("no name")
	}
	req := &pb.UserReadReq{
		Name: name,
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	resp, err := c.client.Delete(ctx, req)
	if err != nil {
		log.Printf("error during read: %v", err)
		return "", err
	}
	if resp.Success {
		return resp.Message, nil
	}
	return resp.ErrMessage, nil
}
