package api

import (
	"Server-for-Ecommerce/app_v2/product_service/util"
	"Server-for-Ecommerce/library/encode/gzip"
	"Server-for-Ecommerce/library/server"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
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
	DeleteCache(ctx context.Context, req *DeleteCacheRequest) (res *DeleteCacheResponse, err error)
}

func NewHttpHandler(httpPattern string, s HttpServer) *mux.Router {
	r := mux.NewRouter().PathPrefix(httpPattern).Subrouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	r.HandleFunc("/cache", deleteCacheHandler(s)).Methods(DELETE, http.MethodOptions)
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

		var (
			page     int64
			pageSize int64
			err      error
		)

		query := r.URL.Query()
		pageString := query.Get("page")
		if pageString != util.EmptyString {
			page, err = strconv.ParseInt(pageString, util.Base10Int, util.BitSize64)
			if err != nil {
				server.HTTPError(w, r, err)
				return
			}
		} else {
			page = 1
		}

		pageSizeString := query.Get("page_size")
		if pageSizeString != util.EmptyString {
			pageSize, err = strconv.ParseInt(pageSizeString, util.Base10Int, util.BitSize64)
			if err != nil {
				server.HTTPError(w, r, err)
				return
			}
		} else {
			pageSize = 20
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

func deleteCacheHandler(s HttpServer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		decode := json.NewDecoder(r.Body)
		req := &DeleteCacheRequest{}
		err := decode.Decode(req)
		if err != nil && err != io.EOF {
			server.HTTPError(w, r, err)
			return
		}

		resp, err := s.DeleteCache(ctx, req)
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
		server.ForwardResponseMessage(ctx, gzip.NewGzipEncoder(), w, r, resp)
	}
}
