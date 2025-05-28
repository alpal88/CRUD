package grpc

import (
	"Desktop/golangProjects/CRUD/pkg/database"
	pb "Desktop/golangProjects/CRUD/pkg/grpc/proto"
	server "Desktop/golangProjects/CRUD/pkg/grpc/server"
	"context"
	"errors"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var _ database.Database = &databaseMock{}

type databaseMock struct {
	db map[string]int
}

func dbMockNew() *databaseMock {
	return &databaseMock{
		db: make(map[string]int),
	}
}

func (d *databaseMock) Write(name string, age string) error {
	ageAsInt, err := strconv.Atoi(age)
	if err != nil {
		return err
	}
	d.db[name] = ageAsInt
	return nil
}

func (d *databaseMock) Read(name string) (string, error) {
	age, ok := d.db[name]
	if !ok {
		return "", database.ErrUserNotFound
	}
	return strconv.Itoa(age), nil
}

func (d *databaseMock) Delete(name string) error {
	_, ok := d.db[name]
	if !ok {
		return errors.New("user not in database")
	}
	delete(d.db, name)
	return nil
}

func (d *databaseMock) Close() error {
	return nil
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func initGRPCServer(t *testing.T) *grpc.ClientConn {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	mockDb := dbMockNew()
	server := server.New(mockDb)
	pb.RegisterCRUDServer(s, server)

	go func() {
		require.NoError(t, s.Serve(lis))
	}()

	conn, err := grpc.Dial("bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	return conn
}

func TestClient(t *testing.T) {
	conn := initGRPCServer(t)
	defer conn.Close()

	c, err := New(conn)
	require.NoError(t, err)
	resp, err := c.CreateUser("jason", 21)
	require.NoError(t, err)
	require.Equal(t, "user jason age 21 has been successfully created", resp)

	resp, err = c.ReadUser("jason")
	require.NoError(t, err)
	require.Equal(t, "user jason's age is 21", resp)

	resp, err = c.ReadUser("j")
	require.NoError(t, err)
	require.Equal(t, "user j not found in our database", resp)

	resp, err = c.UpdateUser("jason", 27)
	require.NoError(t, err)
	require.Equal(t, "user jason's new age is 27", resp)

	resp, err = c.DeleteUser("jason")
	require.NoError(t, err)
	require.Equal(t, "deleting user jason from the database", resp)

	resp, err = c.DeleteUser("jason")
	require.NoError(t, err)
	require.Equal(t, "error in deleting from the database", resp)
}
