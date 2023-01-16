package api

import (
	"context"
	"encoding/json"
	"github.com/DragonPow/Server-for-Ecommerce/library/encode/gzip"
	"github.com/DragonPow/Server-for-Ecommerce/library/server"
	"io"
	"io/ioutil"
	"net/http"
)

func (r *myRouter) RegisterProductHandler() {
	p := r.PathPrefix(productRouter).Subrouter()
	p.HandleFunc("", r.addProductHandler).Methods(POST)
	p.HandleFunc("/{id}", r.updateProductHandler).Methods(PUT)
	p.HandleFunc("/{id}", r.deleteProductHandler).Methods(DELETE)
}

func (r *myRouter) addProductHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeOutInSeconds)
	defer cancel()

	request := &AddProductRequest{}
	// Read the request body
	decode := json.NewDecoder(req.Body)
	err := decode.Decode(request)
	if err != nil && err != io.EOF {
		server.HTTPError(w, req, err)
		return
	}

	resp, err := r.service.AddProduct(ctx, request)
	if err != nil {
		server.HTTPError(w, req, err)
		return
	}

	server.ForwardResponseMessage(ctx, gzip.NewGzipEncoder(), w, req, resp)
}

func (r *myRouter) updateProductHandler(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeOutInSeconds)
	defer cancel()

	// Get ProductId
	productId, err := ParseInt64FromReq(req, "id")
	if err != nil {
		server.HTTPError(w, req, err)
		return
	}
	// Read the request body
	variants, err := ioutil.ReadAll(req.Body)
	if err != nil {
		server.HTTPError(w, req, err)
		return
	}

	request := &UpdateProductRequest{
		Id:       productId,
		Variants: variants,
	}
	resp, err := r.service.UpdateProduct(ctx, request)
	if err != nil {
		server.HTTPError(w, req, err)
		return
	}

	server.ForwardResponseMessage(ctx, gzip.NewGzipEncoder(), w, req, resp)
}

func (r *myRouter) deleteProductHandler(w http.ResponseWriter, req *http.Request) {
	panic("Implement me")
}
