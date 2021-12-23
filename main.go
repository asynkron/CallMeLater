package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	//X-Later-Request-Url 	- url to call
	requestUrl, err := url.Parse(r.Header.Get("X-Later-Request-Url"))
	if err != nil {
		return
	}
	//X-Later-When 			- UTC timestamp
	layout := "2006-01-02 15:04:05 -0700 MST"
	time, err := time.Parse(layout, r.Header.Get("X-Later-When"))
	if err != nil {
		return
	}
	//X-Later-Response-Url 	- webhook to send results to
	responseUrl, err := url.Parse(r.Header.Get("X-Later-Response-Url"))
	if err != nil {
		return
	}

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
