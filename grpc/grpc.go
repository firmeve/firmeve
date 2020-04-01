package grpc

import (
	"context"
	"google.golang.org/grpc/metadata"
)

type GRPC struct {
	original context.Context
}

func (G GRPC) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (G GRPC) Write(p []byte) (n int, err error) {
	panic("implement me")
}

func (G GRPC) Name() string {
	return `grpc`
}

func (G GRPC) Metadata() map[string][]string {
	panic("implement me")
}

func (G GRPC) Message() ([]byte, error) {
	panic("implement me")
}

func (G *GRPC) Values() map[string][]string {
	v, _ := metadata.FromIncomingContext(G.original)
	return v
}
