package util

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorGrpc(code int, message string) error {
	return status.Error(codes.Code(code), message)
}

func ErrorGrpcf(code int, formatMessage string, dataFormat ...interface{}) error {
	return status.Errorf(codes.Code(code), formatMessage, dataFormat...)
}

func Response[T any](err error, defaultResponse func(code int32, message string) *T) (*T, error) {
	if s, ok := status.FromError(err); ok {
		return defaultResponse(int32(s.Code()), s.Message()), nil
	}
	return nil, err
}
