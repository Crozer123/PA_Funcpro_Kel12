package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
)

func main() {
	// konfig comsole
	log.SetFlags(0)


	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func (w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		data := map[string]string{"data":"ping"}
		json.NewEncoder(w).Encode(data)
	}) 

	var finalHandler http.Handler = mux
	finalHandler = middleware.Logging(finalHandler)


	http.ListenAndServe(":8080", finalHandler);
}