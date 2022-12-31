package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/util"
	"github.com/DragonPow/Server-for-Ecommerce/library/cache"
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

func NewHttpHandler(s HttpServer) error {
	r := mux.NewRouter()

	r.HandleFunc("/products/{id}", getDetailProductHandler(s)).Methods(GET)
	r.HandleFunc("/products", getListProductHandler).Methods(GET)

	return nil
}

func getDetailProductHandler(m *http.ServeMux, s HttpServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		productId, err := strconv.ParseInt(mux.Vars(r)["id"], util.Base10Int, util.BitSize64)
		if err != nil {
			HTTPError(ctx, m, w, r, err)
			return
		}
		resp, err := s.GetDetailProduct(ctx, &GetDetailProductRequest{
			Id: productId,
		})
		if err != nil {
			HTTPError(ctx, m, w, r, err)
			return
		}
		ForwardResponseMessage(ctx, mux, cache.Marshal(), w, r, resp)
	}
}

func getListProductHandler(w http.ResponseWriter, r *http.Request) {

}

type Marshaler interface {
	Marshal(v any) ([]byte, error)
}

// ForwardResponseMessage forwards the message "resp" from gRPC server to REST client.
func ForwardResponseMessage(ctx context.Context, mux *http.ServeMux, marshaler Marshaler, w http.ResponseWriter, req *http.Request, resp any) {
	contentType := ""
	w.Header().Set("Content-Type", contentType)

	buf, err := marshaler.Marshal(resp)
	if err != nil {
		HTTPError(ctx, mux, w, req, err)
		return
	}

	if _, err = w.Write(buf); err != nil {
		HTTPError(ctx, mux, w, req, err)
		return
	}
}

func HTTPError(ctx context.Context, mux *http.ServeMux, w http.ResponseWriter, req *http.Request, err error) {
	panic("unimplemented")
}
