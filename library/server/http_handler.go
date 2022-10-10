package server

import "net/http"

type HTTPServerHandler func(*http.ServeMux)
