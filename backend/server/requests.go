package server

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestData struct {
	Id                 string              `json:"request_id,omitempty"`
	RequestMethod      string              `json:"request_method,omitempty"`
	Header             map[string][]string `json:"header,omitempty"`
	Form               map[string][]string `json:"form,omitempty"`
	RequestUrl         string              `json:"request_url,omitempty"`
	ResponseUrl        string              `json:"response_url,omitempty"`
	ResponseMethod     string              `json:"response_method,omitempty"`
	ScheduledTimestamp time.Time           `json:"when"`
	Body               []byte              `json:"body,omitempty"`
	ParentId           string              `json:"parent_id,omitempty"`
}

func (w *worker) sendRequestResponse(rd *RequestData) {
	//if the request fails after this, it will be lost

	response, err := sendRequest(rd)

	if err != nil {
		log.
			Err(err).
			Msg("Error sending request")
		return
	}

	if rd.ResponseUrl != "" {
		_ = w.storage.Push(response)
		w.requests <- response
	}
}

func sendRequest(p *RequestData) (*RequestData, error) {
	log.
		Info().
		Str("Url", p.RequestUrl).
		Msg("Sending request")

	var r io.Reader
	request, err := http.NewRequest(p.RequestMethod, p.RequestUrl, r)
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

	var res = &RequestData{
		Id:                 uuid.New().String(),
		Header:             response.Header,
		Body:               body,
		RequestUrl:         p.ResponseUrl,
		RequestMethod:      p.ResponseMethod,
		ScheduledTimestamp: p.ScheduledTimestamp,
		ParentId:           p.Id,
	}

	return res, nil
}
