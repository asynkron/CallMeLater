package main

import (
	json2 "encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	//X-Later-Request-Url 	- url to sendRequest
	requestUrl, err := url.Parse(r.Header.Get("X-Later-Request-Url"))
	if err != nil {
		log.
			Err(err).
			Msg("failed to parse request url")

		return
	}
	//X-Later-When 			- UTC timestamp
	when, err := time.ParseDuration(r.Header.Get("X-Later-When"))
	if err != nil {
		log.
			Err(err).
			Msg("failed to parse when")

		return
	}
	//X-Later-Response-Url 	- webhook to send results to
	tmp := r.Header.Get("X-Later-Response-Url")
	var responseUrl *url.URL
	if tmp != "" {
		responseUrl, err = url.Parse(tmp)
		if err != nil {
			log.
				Err(err).
				Msg("failed to parse response url")
			return
		}
	}

	body, err := ioutil.ReadAll(r.Body)

	t := time.Now().Add(when)

	var p = &requestData{
		RequestId:   uuid.New().String(),
		RequestUrl:  requestUrl.String(),
		ResponseUrl: responseUrl.String(),
		When:        t,
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

	err = storage.Set(rd)
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
