package main

import (
	"bytes"
	json2 "encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type RequestStorage interface {
	Get(id string) (*requestData, error)
	Set(id string, data *requestData) error
}

type NoopStorage struct{}

func (n NoopStorage) Get(id string) (*requestData, error) {
	log.Printf("NoopStorage.Get(%s)", id)
	return nil, nil
}

func (n NoopStorage) Set(id string, data *requestData) error {
	log.Printf("NoopStorage.Set(%s)", id)
	return nil
}

type requestData struct {
	RequestId   string
	Method      string
	Header      map[string][]string
	Form        map[string][]string
	RequestUrl  string
	ResponseUrl string
	When        time.Time
	Body        []byte
}

type responseData struct {
	Header      map[string][]string
	Body        []byte
	ResponseUrl string
	Method      string
}

var (
	storage  RequestStorage = NoopStorage{}
	requests                = make(chan *requestData)
)

func consumeLoop() {
	log.Print("Consume loop started")
	for {
		rd := <-requests
		go sendRequestResponse(rd)
	}
}

func sendRequestResponse(rd *requestData) {
	response, err := sendRequest(rd)
	if err != nil {
		log.Printf("Error making request: %s", err)
		return
	}
	err = sendResponse(response)
	if err != nil {
		log.Printf("Error sending response: %s", err)
		return
	}
}

func sendResponse(rd *responseData) error {
	log.Printf("Sending response to %s", rd.ResponseUrl)
	var r io.Reader
	request, err := http.NewRequest(rd.Method, rd.ResponseUrl, r)
	if err != nil {
		return err
	}
	request.Header = rd.Header
	request.Body = ioutil.NopCloser(bytes.NewReader(rd.Body))
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	log.Printf("Response: %s", string(body))
	return nil
}

func sendRequest(p *requestData) (*responseData, error) {
	log.Printf("Sending request: %s", p.RequestUrl)
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
		Header:      response.Header,
		Body:        body,
		ResponseUrl: p.ResponseUrl,
		Method:      "POST",
	}

	return res, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	//X-Later-Request-Url 	- url to sendRequest
	requestUrl, err := url.Parse(r.Header.Get("X-Later-Request-Url"))
	if err != nil {
		return
	}
	//X-Later-When 			- UTC timestamp
	layout := "2006-01-02 15:04:05 -0700 MST"
	when, err := time.Parse(layout, r.Header.Get("X-Later-When"))
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
		RequestId:   uuid.New().String(),
		RequestUrl:  requestUrl.String(),
		ResponseUrl: responseUrl.String(),
		When:        when,
		Header:      r.Header,
		Form:        r.Form,
		Body:        body,
		Method:      r.Method,
	}

	saveRequest(p)

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, "OK")
	log.Printf("Endpoint Hit: handler %v", p)
}

func saveRequest(rd *requestData) {
	j, err := json2.Marshal(rd)
	if err != nil {
		return
	}
	json := string(j)

	err = storage.Set(rd.RequestId, rd)
	if err != nil {
		log.Printf("Error saving request: %s", err)
		return
	}
	log.Printf("Saved request: %s", json)
	requests <- rd
}

func handleRequests() {
	http.HandleFunc("/later", handler)
	err := http.ListenAndServe(":10000", nil)
	log.Fatal(err)
}

func main() {
	go consumeLoop()
	handleRequests()
}
