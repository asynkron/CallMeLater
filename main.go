package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	rtr := mux.NewRouter()
	rtr.HandleFunc(`/later/{rest:.+/}`, homePage)
	http.Handle("/", rtr)
	err := http.ListenAndServe(":10000", nil)
	log.Fatal(err)
}

func main() {
	handleRequests()
}
