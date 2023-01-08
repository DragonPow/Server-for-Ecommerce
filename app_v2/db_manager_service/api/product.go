package api

import (
	"net/http"
)

func (r *myRouter) RegisterProductHandler() {
	p := r.PathPrefix(productRouter).Subrouter()
	p.HandleFunc("/", addProductHandler).Methods(POST)
	p.HandleFunc("/", updateProductHandler).Methods(PUT)
	p.HandleFunc("/", deleteProductHandler).Methods(DELETE)
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {

}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {

}

func deleteProductHandler(w http.ResponseWriter, r *http.Request) {

}