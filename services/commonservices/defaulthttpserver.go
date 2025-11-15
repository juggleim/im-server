package commonservices

import (
	"fmt"
	"net/http"
)

var mux *http.ServeMux

func InitDefaultHttpServer() {
	mux = http.NewServeMux()
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"version":"%s"}`, "1.8.6")
	})
}

func GetDefaultHttpServeMux() *http.ServeMux {
	return mux
}
