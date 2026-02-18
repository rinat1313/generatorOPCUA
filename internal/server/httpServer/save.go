package httpServer

import (
	"net/http"
)

func SaveHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	//if r.Method != "POST" {
	//	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	//	return
	//}
	//contentType := r.Header.Get("Content-Type")
	//if contentType != "application/json" {
	//	http.Error(w, "Content-Type not supported", http.StatusUnsupportedMediaType)
	//	return
	//}
	//r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	//body, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//defer r.Body.Close()
	////id, err := service.SavingMessage(string(body))
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//w.WriteHeader(http.StatusOK)
	//err = json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

}

//func init() {
//AddNewFunction(CreateHandlerCommand("/save", "POST", SaveHandler))
//}
