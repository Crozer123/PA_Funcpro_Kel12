package main

import (
	"log"
	"net/http"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/handler"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/platform/database"
	"github.com/joho/godotenv"
)

func main() {
	// Konfigurasi database
	err := godotenv.Load()
	if err != nil {
		log.Println("error saat load env")
	}

	db := database.NewConnection()
	db.AutoMigrate(&domain.User{})
	// konfig comsole
	log.SetFlags(0)


	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func (w http.ResponseWriter, r *http.Request){
		type User struct {
        ID   string 
        Name string
    }
    user := User{ID: "123", Name: "John Doe"}
		response.WriteJSON(w, http.StatusOK, "succes", user)
	}) 

	var finalHandler http.Handler = mux
	finalHandler = middleware.Logging(finalHandler)


	http.ListenAndServe(":8080", finalHandler);
}	