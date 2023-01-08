package server

import (
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type HTTPServerHandler func(*http.ServeMux)

type Marshaller interface {
	Marshal(v any) ([]byte, error)
}

// ForwardResponseMessage forwards the message "resp" from gRPC server to REST client.
func ForwardResponseMessage(ctx context.Context, marshaler Marshaller, w http.ResponseWriter, req *http.Request, resp any) {
	contentType := "application/json"
	//acceptEncoding := "gzip"

	w.Header().Set("Content-Type", contentType)
	//w.Header().Set("Content-Encoding", acceptEncoding)

	buf, err := json.Marshal(resp)
	if err != nil {
		HTTPError(w, req, err)
		return
	}

	if _, err = w.Write(buf); err != nil {
		HTTPError(w, req, err)
		return
	}
}

func HTTPError(w http.ResponseWriter, req *http.Request, err error) {
	v := make(map[string]any)
	if statusErr, ok := status.FromError(err); ok {
		v["code"] = statusErr.Code()
		v["message"] = statusErr.Message()
		w.WriteHeader(runtime.HTTPStatusFromCode(statusErr.Code()))
	} else {
		v["code"] = codes.Internal
		v["message"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	b, _ := json.Marshal(v)
	w.Write(b)
}
