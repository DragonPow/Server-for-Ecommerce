package api

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"net/http"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type HttpServer interface {
}

func NewHttpHandler(httpPattern string, s HttpServer) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc(httpPattern+"/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			HTTPError(ctx, w, r, err)
			return
		}
	}).Methods(GET)

	return r
}

type Marshaler interface {
	Marshal(v any) ([]byte, error)
}

// ForwardResponseMessage forwards the message "resp" from gRPC server to REST client.
func ForwardResponseMessage(ctx context.Context, marshaler Marshaler, w http.ResponseWriter, req *http.Request, resp any) {
	contentType := "application/json"
	//acceptEncoding := "gzip"

	w.Header().Set("Content-Type", contentType)
	//w.Header().Set("Content-Encoding", acceptEncoding)

	buf, err := json.Marshal(resp)
	if err != nil {
		HTTPError(ctx, w, req, err)
		return
	}

	if _, err = w.Write(buf); err != nil {
		HTTPError(ctx, w, req, err)
		return
	}
}

func HTTPError(ctx context.Context, w http.ResponseWriter, req *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	v := make(map[string]any)
	v["code"] = codes.Internal
	v["message"] = err.Error()
	b, _ := json.Marshal(v)
	w.Write(b)
}
