package main

import (
	"fmt"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	//X-Later-Request-Url 	- url to call
	//X-Later-When 			- UTC timestamp
	//X-Later-Response-Url 	- webhook to send results to

	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/later", homePage)
	err := http.ListenAndServe(":10000", nil)
	log.Fatal(err)
}

func main() {
	handleRequests()
}
