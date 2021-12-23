package main

import (
	"bytes"
	json2 "encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type requestData struct {
	Method      string
	Header      map[string][]string
	Form        map[string][]string
	RequestUrl  string
	ResponseUrl string
	When        time.Time
	Body        []byte
}

type responseData struct {
	Header map[string][]string
	Body   []byte
}

func makeRequest(p *requestData) (*responseData, error) {
	var r io.Reader
	request, err := http.NewRequest(p.Method, p.RequestUrl, r)
	if err != nil {
		return nil, err
	}
	request.Header = p.Header
	request.Form = p.Form
	request.Body = ioutil.NopCloser(bytes.NewReader(p.Body))
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res = &responseData{
		Header: response.Header,
		Body:   body,
	}

	return res, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	//X-Later-Request-Url 	- url to makeRequest
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

	var p = &requestData{
		RequestUrl:  requestUrl.String(),
		ResponseUrl: responseUrl.String(),
		When:        time,
		Header:      r.Header,
		Form:        r.Form,
		Body:        body,
		Method:      r.Method,
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
