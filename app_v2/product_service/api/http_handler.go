package api

import (
	"context"
	"encoding/json"
	"github.com/DragonPow/Server-for-Ecommerce/library/encode/gzip"
	"github.com/DragonPow/Server-for-Ecommerce/library/server"
	"io/ioutil"
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
	r := mux.NewRouter().PathPrefix(httpPattern).Subrouter()
	r.HandleFunc(httpPattern+"/", func(w http.ResponseWriter, r *http.Request) {
		_, cancel := context.WithCancel(context.Background())
		defer cancel()
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
	}).Methods(GET)
	r.HandleFunc("/products/{id}", getDetailProductHandler(s)).Methods(GET)
	r.HandleFunc("/products", getListProductHandler(s)).Methods(GET)

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
		server.ForwardResponseMessage(ctx, gzip.NewGzipEncoder(), w, r, resp)
	}
}

func getListProductHandler(s HttpServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}

		var req *GetListProductRequest
		err = json.Unmarshal(body, req)
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
		resp, err := s.GetListProduct(ctx, req)
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
		server.ForwardResponseMessage(ctx, gzip.NewGzipEncoder(), w, r, resp)
	}
}
