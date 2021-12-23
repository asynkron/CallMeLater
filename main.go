package main

import (
	json2 "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type payload struct {
	Header      map[string][]string
	Form        map[string][]string
	RequestUrl  string
	ResponseUrl string
	When        time.Time
	Body        []byte
}

func handler(w http.ResponseWriter, r *http.Request) {
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

	body, err := ioutil.ReadAll(r.Body)

	var p = &payload{
		RequestUrl:  requestUrl.String(),
		ResponseUrl: responseUrl.String(),
		When:        time,
		Header:      r.Header,
		Form:        r.Form,
		Body:        body,
	}

	j, err := json2.Marshal(p)

	fmt.Fprintf(w, "json "+string(j))
	fmt.Println("Endpoint Hit: handler")
}

func handleRequests() {
	http.HandleFunc("/later", handler)
	err := http.ListenAndServe(":10000", nil)
	log.Fatal(err)
}

func main() {
	handleRequests()
}
