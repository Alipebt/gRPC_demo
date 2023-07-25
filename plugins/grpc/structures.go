package grpc

import (
	"context"
)

//skywalking:native google.golang.org/grpc/internal/transport Stream
type nativeStream struct {
	ctx    context.Context
	method string
}

func (s *nativeStream) Method() string {
	return s.method
}

func (s *nativeStream) Context() context.Context {
	return s.ctx
}

//skywalking:native google.golang.org/grpc clientStream
type nativeclientStream struct {
	ctx context.Context
}

func (cs *nativeclientStream) Context() context.Context {
	return cs.ctx
}

//skywalking:native google.golang.org/grpc serverStream
type nativeserverStream struct {
	ctx context.Context
}

func (cs *nativeserverStream) Context() context.Context {
	return cs.ctx
}
