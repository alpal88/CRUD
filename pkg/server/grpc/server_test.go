package grpc

import (
	"Desktop/golangProjects/CRUD/pkg/server/database"
	pb "Desktop/golangProjects/CRUD/proto"
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

type deps struct {
	fileName string
	db       database.Database
}

func initDeps(t *testing.T) deps {
	tmpFile, err := os.CreateTemp("", "testdb-*.db")
	require.NoError(t, err)
	db, err := database.New(tmpFile.Name(), 0666, nil)
	require.NoError(t, err)
	return deps{
		fileName: tmpFile.Name(),
		db:       db,
	}
}

func initGRPCServer(t *testing.T) (*grpc.ClientConn, func()) {
	deps := initDeps(t)
	lis = bufconn.Listen(bufSize)
	server := New(deps.db)
	s := grpc.NewServer()
	pb.RegisterCRUDServer(s, server)
	go func() {
		require.NoError(t, s.Serve(lis))
	}()

	// conn, err := grpc.NewClient("bufnet", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(bufDialer))
	// require.NoError(t, err)

	conn, err := grpc.Dial(
		"bufnet", // dummy name that matches the custom dialer
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	cleanup := func() {
		deps.db.Close()
		os.Remove(deps.fileName)
		s.Stop()
		conn.Close()
	}

	return conn, cleanup
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGRPCServerCreate(t *testing.T) {
	grpcConn, cleanup := initGRPCServer(t)
	defer cleanup()

	client := pb.NewCRUDClient(grpcConn)
	ctx := context.Background()
	age := int32(21)
	createResp, err := client.Create(ctx, &pb.UserWriteReq{
		Name: "jason",
		Age:  &age,
	})
	require.NoError(t, err)
	require.True(t, createResp.Success)
	require.Equal(t, fmt.Sprintf(CreateMessage, "jason", 21), createResp.Message)
	require.Equal(t, "", createResp.ErrMessage)

	createResp, err = client.Create(ctx, &pb.UserWriteReq{
		Name: "jason",
	})
	require.NoError(t, err)
	require.False(t, createResp.Success)
	require.Equal(t, "", createResp.Message)
	require.Equal(t, errMessageAgeOrName, createResp.ErrMessage)
}

func TestGRPCServerRead(t *testing.T) {
	grpcConn, cleanup := initGRPCServer(t)
	defer cleanup()

	client := pb.NewCRUDClient(grpcConn)
	ctx := context.Background()
	age := int32(21)
	_, err := client.Create(ctx, &pb.UserWriteReq{
		Name: "jason",
		Age:  &age,
	})
	require.NoError(t, err)

	readResponse, err := client.Read(ctx, &pb.UserReadReq{
		Name: "jason",
	})
	require.NoError(t, err)
	require.True(t, readResponse.Success)
	require.Equal(t, fmt.Sprintf(ReadMessage, "jason", "21"), readResponse.Message)
	require.Equal(t, "", readResponse.ErrMessage)

	// no name
	readResponse, err = client.Read(ctx, &pb.UserReadReq{})
	require.NoError(t, err)
	require.False(t, readResponse.Success)
	require.Equal(t, "", readResponse.Message)
	require.Equal(t, errMessageName, readResponse.ErrMessage)

	// not in database
	readResponse, err = client.Read(ctx, &pb.UserReadReq{
		Name: "steven",
	})
	require.NoError(t, err)
	require.False(t, readResponse.Success)
	require.Equal(t, "", readResponse.Message)
	require.Equal(t, ErrInRead.Error(), readResponse.ErrMessage)

}

func TestGRPCServerUpdate(t *testing.T) {
	grpcConn, cleanup := initGRPCServer(t)
	defer cleanup()

	client := pb.NewCRUDClient(grpcConn)
	ctx := context.Background()
	age := int32(21)
	_, err := client.Create(ctx, &pb.UserWriteReq{
		Name: "jason",
		Age:  &age,
	})
	require.NoError(t, err)

	age = int32(27)
	updateResp, err := client.Update(ctx, &pb.UserWriteReq{
		Name: "jason",
		Age:  &age,
	})
	require.NoError(t, err)
	require.True(t, updateResp.Success)
	require.Equal(t, fmt.Sprintf(UpdateMessage, "jason", 27), updateResp.Message)
	require.Equal(t, "", updateResp.ErrMessage)
	updateResp, err = client.Create(ctx, &pb.UserWriteReq{})
	require.NoError(t, err)
	require.False(t, updateResp.Success)
	require.Equal(t, "", updateResp.Message)
	require.Equal(t, errMessageAgeOrName, updateResp.ErrMessage)
}

func TestGRPCServerDelete(t *testing.T) {
	grpcConn, cleanup := initGRPCServer(t)
	defer cleanup()

	client := pb.NewCRUDClient(grpcConn)
	ctx := context.Background()
	age := int32(21)
	_, err := client.Create(ctx, &pb.UserWriteReq{
		Name: "jason",
		Age:  &age,
	})
	require.NoError(t, err)

	// no name
	delResponse, err := client.Delete(ctx, &pb.UserReadReq{})
	require.NoError(t, err)
	require.False(t, delResponse.Success)
	require.Equal(t, "", delResponse.Message)
	require.Equal(t, errMessageName, delResponse.ErrMessage)

	// not in database
	delResponse, err = client.Read(ctx, &pb.UserReadReq{
		Name: "steven",
	})
	require.NoError(t, err)
	require.False(t, delResponse.Success)
	require.Equal(t, "", delResponse.Message)
	require.Equal(t, ErrInRead.Error(), delResponse.ErrMessage)

	delResponse, err = client.Delete(ctx, &pb.UserReadReq{
		Name: "jason",
	})
	require.NoError(t, err)
	require.True(t, delResponse.Success)
	require.Equal(t, fmt.Sprintf(DeleteMessage, "jason"), delResponse.Message)
	require.Equal(t, "", delResponse.ErrMessage)

}
