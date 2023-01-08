package api

import (
	"context"
	"github.com/DragonPow/Server-for-Ecommerce/app_v2/db_manager_service/util"
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

	productRouter  = "/products"
	categoryRouter = "/categories"
	uomRouter      = "/uoms"
	sellerRouter   = "/sellers"
	userRouter     = "/users"
)

type HttpServer interface {
	GetTimeOutHttpInSecond() int
	AddProduct(ctx context.Context, req *AddProductRequest) (res *AddProductResponse, err error)
	UpdateProduct(ctx context.Context, req *UpdateProductRequest) (res *UpdateProductResponse, err error)
	DeleteProduct(ctx context.Context, req *DeleteProductRequest) (res *DeleteProductResponse, err error)
}

type myRouter struct {
	*mux.Router
	service          HttpServer
	timeOutInSeconds time.Duration
}

func NewHttpHandler(httpPattern string, s HttpServer) http.Handler {
	router := &myRouter{
		Router:           mux.NewRouter().PathPrefix(httpPattern).Subrouter(),
		service:          s,
		timeOutInSeconds: time.Duration(s.GetTimeOutHttpInSecond()) * time.Second,
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, cancel := context.WithTimeout(context.Background(), router.timeOutInSeconds)
		defer cancel()
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			server.HTTPError(w, r, err)
			return
		}
	}).Methods(GET)

	// Register some route
	router.RegisterProductHandler()

	return router
}

func ParseIdFromRequest(w http.ResponseWriter, req *http.Request, fieldName string) (int64, error) {
	productId, err := strconv.ParseInt(mux.Vars(req)[fieldName], util.Base10Int, util.BitSize64)
	if err != nil {
		return 0, err
	}
	return productId, nil
}
