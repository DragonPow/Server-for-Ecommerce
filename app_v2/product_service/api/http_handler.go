package api

import (
	"context"
	"encoding/json"
	"github.com/DragonPow/Server-for-Ecommerce/library/encode/gzip"
	"google.golang.org/grpc/codes"
	"net/http"
	"strconv"

	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/util"
	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type HttpServer interface {
	GetDetailProduct(ctx context.Context, req *GetDetailProductRequest) (res *GetDetailProductResponse, err error)
	GetListProduct(ctx context.Context, req *GetListProductRequest) (res *GetListProductResponse, err error)
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
	r.HandleFunc(httpPattern+"/products/{id}", getDetailProductHandler(s)).Methods(GET)
	r.HandleFunc(httpPattern+"/products", getListProductHandler).Methods(GET)

	return r
}

func getDetailProductHandler(s HttpServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		productId, err := strconv.ParseInt(mux.Vars(r)["id"], util.Base10Int, util.BitSize64)
		if err != nil {
			HTTPError(ctx, w, r, err)
			return
		}
		resp, err := s.GetDetailProduct(ctx, &GetDetailProductRequest{
			Id: productId,
		})
		if err != nil {
			HTTPError(ctx, w, r, err)
			return
		}
		ForwardResponseMessage(ctx, gzip.NewGzipEncoder(), w, r, resp)
	}
}

func getListProductHandler(w http.ResponseWriter, r *http.Request) {

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
