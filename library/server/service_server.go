package server

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
)

type ServiceServer interface {
	RegisterWithServer(*grpc.Server)
	RegisterWithHttpHandler(firstPattern string) (http.Handler, error)
	RegisterWithHandler(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
	Close(context.Context)
}
