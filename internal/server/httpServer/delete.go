package httpServer

import "net/http"

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

//func init() {
//	AddNewFunction(CreateHandlerCommand("/delete", "DELETE", DeleteHandler))
//}
