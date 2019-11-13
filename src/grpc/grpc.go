package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"strings"

	"google.golang.org/grpc/metadata"

	pb "github.com/lain-m21/go-reverseproxy-test/src/grpc/proto"
)

type (
	testServerImpl struct{}
)

func NewBackendServer() *grpc.Server {
	server := grpc.NewServer()
	testServer := newTestServer()
	pb.RegisterTestServer(server, testServer)
	return server
}

func newTestServer() pb.TestServer {
	return &testServerImpl{}
}

func (s *testServerImpl) Say(ctx context.Context, r *pb.TestRequest) (*pb.TestResponse, error) {
	message := fmt.Sprintf("Test service received a message: %s", r.Message)
	response := &pb.TestResponse{Message: message}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		flags, ok := md["flags"]
		if ok {
			flagSet := strings.Split(flags[0], ",")
			response.Flags = flagSet
		}
	}

	return response, nil
}
