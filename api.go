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

const (
	HeaderRequestUrl     = "X-Later-Request-Url"
	HeaderWhen           = "X-Later-When"
	HeaderResponseUrl    = "X-Later-Response-Url"
	HeaderResponseMethod = "X-Later-Response-Method"
)

func handler(w http.ResponseWriter, r *http.Request) {
	requestUrl, err := url.Parse(r.Header.Get(HeaderRequestUrl))
	if err != nil {
		log.
			Err(err).
			Msg("Failed to parse request url")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to parse "+HeaderRequestUrl)
		return
	}
	when, err := time.ParseDuration(r.Header.Get(HeaderWhen))
	if err != nil {
		log.
			Err(err).
			Msg("failed to parse when")

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Failed to parse "+HeaderWhen)
		return
	}
	tmp := r.Header.Get(HeaderResponseUrl)
	var responseUrlStr string
	if tmp != "" {
		responseUrl, err := url.Parse(tmp)
		if err != nil {
			log.
				Err(err).
				Msg("failed to parse response url")

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Failed to parse "+HeaderResponseUrl)
			return
		}
		responseUrlStr = responseUrl.String()
	}

	body, err := ioutil.ReadAll(r.Body)

	t := time.Now().Add(when)

	var p = &requestData{
		RequestId:      uuid.New().String(),
		RequestUrl:     requestUrl.String(),
		RequestMethod:  r.Method,
		ResponseUrl:    responseUrlStr,
		ResponseMethod: r.Header.Get(HeaderResponseMethod),
		When:           t,
		Header:         r.Header,
		Form:           r.Form,
		Body:           body,
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
