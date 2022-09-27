package server

import (
	"ethereum-monitor/internal/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func InitServer(dataHandler *handlers.Data) {
	//Create a mux router
	r := mux.NewRouter()

	// We will define endpoints
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	}).Methods("GET")
	r.HandleFunc("/address/add", dataHandler.AddAddress)
	r.HandleFunc("/address/get", dataHandler.GetAddress)
	r.HandleFunc("/address/getBalance/{address:[a-zA-Z0-9]*}", dataHandler.GetBalance)

	log.Fatal(http.ListenAndServe(":8080", r))
}
