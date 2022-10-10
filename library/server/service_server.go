package server

import (
	"context"
	"google.golang.org/grpc"
)

type ServiceServer interface {
	RegisterWithServer(*grpc.Server)
	//RegisterWithHandler(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
	Close(context.Context)
}
