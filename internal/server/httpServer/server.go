package httpServer

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var DefaultFuncs []HandlerCommand

func AddNewFunction(f HandlerCommand) {
	DefaultFuncs = append(DefaultFuncs, f)
}

type HttpHandler struct {
	Router *mux.Router
}

func NewHttpHandler() *HttpHandler {
	router := mux.NewRouter()
	handler := &HttpHandler{Router: router}
	for _, v := range DefaultFuncs {
		fmt.Printf("add function name: %s and method name: %s\n", v.Name, v.Method)
		router.HandleFunc(v.Name, v.HttpFunc).Methods(v.Method)
	}
	return handler
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}
