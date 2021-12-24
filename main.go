package main

import (
	"bytes"
	json2 "encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type RequestStorage interface {
	Get(id string) (*requestData, error)
	Set(id string, data *requestData) error
}

type NoopStorage struct{}

func (n NoopStorage) Get(id string) (*requestData, error) {
	log.
		Info().
		Str("id", id).
		Msg("NoopStorage.Get")
	return nil, nil
}

func (n NoopStorage) Set(id string, data *requestData) error {
	log.
		Info().
		Str("id", id).
		Msg("NoopStorage.Set")

	return nil
}

type requestData struct {
	RequestId   string              `json:"request_id,omitempty"`
	Method      string              `json:"method,omitempty"`
	Header      map[string][]string `json:"header,omitempty"`
	Form        map[string][]string `json:"form,omitempty"`
	RequestUrl  string              `json:"request_url,omitempty"`
	ResponseUrl string              `json:"response_url,omitempty"`
	When        time.Time           `json:"when"`
	Body        []byte              `json:"body,omitempty"`
}

type responseData struct {
	Header      map[string][]string `json:"header,omitempty"`
	Body        []byte              `json:"body,omitempty"`
	ResponseUrl string              `json:"response_url,omitempty"`
	Method      string              `json:"method,omitempty"`
}

var (
	storage  RequestStorage = NoopStorage{}
	requests                = make(chan *requestData)
)

func consumeLoop() {
	log.
		Info().
		Msg("Consume loop started")

	for {
		rd := <-requests
		go sendRequestResponse(rd)
	}
}

func sendRequestResponse(rd *requestData) {
	response, err := sendRequest(rd)
	if err != nil {
		log.
			Err(err).
			Msg("Error sending request")
		return
	}
	err = sendResponse(response)
	if err != nil {
		log.
			Err(err).
			Msg("Error sending response")

		return
	}
}

func sendResponse(rd *responseData) error {
	log.
		Info().
		Str("Url", rd.ResponseUrl).
		Msg("Sending response")

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
	log.
		Info().
		Str("Body", string(body)).Str("Url", rd.ResponseUrl).
		Msg("Response sent")

	return nil
}

func sendRequest(p *requestData) (*responseData, error) {
	log.
		Info().
		Str("Url", p.RequestUrl).
		Msg("Sending request")

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
	log.Info().
		Msg("Request accepted")
}

func saveRequest(rd *requestData) {
	j, err := json2.Marshal(rd)
	if err != nil {
		return
	}
	json := string(j)

	err = storage.Set(rd.RequestId, rd)
	if err != nil {
		log.
			Err(err).
			Msg("Error saving request")
		return
	}
	log.
		Info().
		Str("Json", json).Msg("Saved Request")
	requests <- rd
}

func handleRequests() {
	http.HandleFunc("/later", handler)
	err := http.ListenAndServe(":10000", nil)
	log.
		Err(err).
		Msg("Error listening")
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	go consumeLoop()
	handleRequests()
}
