package api

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/product_service/util"
	"github.com/DragonPow/Server-for-Ecommerce/library/encode/gzip"
	"github.com/DragonPow/Server-for-Ecommerce/library/server"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
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
	r := mux.NewRouter().PathPrefix(httpPattern).Subrouter()
	r.HandleFunc(httpPattern+"/", func(w http.ResponseWriter, r *http.Request) {
		_, cancel := context.WithCancel(context.Background())
		defer cancel()
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
	}).Methods(GET, http.MethodOptions)
	r.HandleFunc("/products/{id}", getDetailProductHandler(s)).Methods(GET, http.MethodOptions)
	r.HandleFunc("/products", getListProductHandler(s)).Methods(GET, http.MethodOptions)
	r.Use(mux.CORSMethodMiddleware(r))
	return r
}

func getDetailProductHandler(s HttpServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		productId, err := strconv.ParseInt(mux.Vars(r)["id"], util.Base10Int, util.BitSize64)
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
		resp, err := s.GetDetailProduct(ctx, &GetDetailProductRequest{
			Id: productId,
		})
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
		cacheSince := time.Now().Format(http.TimeFormat)
		cacheUntil := time.Now().Add(60 * time.Second).Format(http.TimeFormat)
		w.Header().Set("Cache-Control", "max-age:60, public")
		w.Header().Set("Last-Modified", cacheSince)
		w.Header().Set("Expires", cacheUntil)
		server.ForwardResponseMessage(ctx, gzip.NewGzipEncoder(), w, r, resp)
	}
}

func getListProductHandler(s HttpServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		query := r.URL.Query()
		page, err := strconv.ParseInt(query.Get("page"), util.Base10Int, util.BitSize64)
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
		pageSize, err := strconv.ParseInt(query.Get("page_size"), util.Base10Int, util.BitSize64)
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
		key := query.Get("key")
		req := &GetListProductRequest{
			Page:     page,
			PageSize: pageSize,
			Key:      key,
		}

		resp, err := s.GetListProduct(ctx, req)
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
		cacheSince := time.Now().Format(http.TimeFormat)
		cacheUntil := time.Now().Add(120 * time.Second).Format(http.TimeFormat)
		w.Header().Set("Cache-Control", "max-age:120, public")
		w.Header().Set("Last-Modified", cacheSince)
		w.Header().Set("Expires", cacheUntil)
		server.ForwardResponseMessage(ctx, gzip.NewGzipEncoder(), w, r, resp)
	}
}
