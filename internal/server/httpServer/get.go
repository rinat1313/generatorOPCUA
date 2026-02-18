package httpServer

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprint(w, "This is new server!")

	vars := mux.Vars(r)
	v, err := strconv.Atoi(vars["v"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
	}
	fmt.Printf("\nПолученны данные: %v", v)
}

//func init() {
//	AddNewFunction(CreateHandlerCommand("/{v:[0-9]+}", "GET", GetHandler))
//}
