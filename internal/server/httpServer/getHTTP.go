package httpServer

import (
	"fmt"
	"net/http"
)

func GetIndexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Push index page to client")
	path := "index.html"
	http.ServeFile(w, r, path)
}

func init() {
	AddNewFunction(CreateHandlerCommand("/", "GET", GetIndexPage))
}
