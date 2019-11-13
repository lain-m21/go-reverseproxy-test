package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/trusch/grpc-proxy/proxy"
	"github.com/trusch/grpc-proxy/proxy/codec"
)

func NewProxy(director func(context.Context, string) (context.Context, *grpc.ClientConn, error)) *grpc.Server {
	handler := grpc.UnknownServiceHandler(proxy.TransparentHandler(director))
	server := grpc.NewServer(handler)
	return server
}

func NewStreamDirector(host string, flags string) func(context.Context, string) (context.Context, *grpc.ClientConn, error) {
	return func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, "flags", flags)
		conn, err := grpc.DialContext(
			ctx,
			host,
			grpc.WithDefaultCallOptions(grpc.CallContentSubtype((&codec.Proxy{}).Name())),
			grpc.WithInsecure(),
		)
		if err != nil {
			fmt.Println(fmt.Sprintf("director failed: %s", err.Error()))
			return nil, nil, grpc.ErrServerStopped
		}
		return ctx, conn, nil
	}
}
