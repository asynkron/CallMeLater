package server

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type RequestData struct {
	RequestId      string              `json:"request_id,omitempty"`
	RequestMethod  string              `json:"request_method,omitempty"`
	Header         map[string][]string `json:"header,omitempty"`
	Form           map[string][]string `json:"form,omitempty"`
	RequestUrl     string              `json:"request_url,omitempty"`
	ResponseUrl    string              `json:"response_url,omitempty"`
	ResponseMethod string              `json:"response_method,omitempty"`
	When           time.Time           `json:"when"`
	Body           []byte              `json:"body,omitempty"`
}

type responseData struct {
	Header         map[string][]string `json:"header,omitempty"`
	Body           []byte              `json:"body,omitempty"`
	ResponseUrl    string              `json:"response_url,omitempty"`
	ResponseMethod string              `json:"response_method,omitempty"`
}

func sendRequestResponse(rd *RequestData) {
	//if the request fails after this, it will be lost

	response, err := sendRequest(rd)

	if err != nil {
		log.
			Err(err).
			Msg("Error sending request")
		return
	}

	if rd.ResponseUrl != "" {
		err = sendResponse(response)
		if err != nil {
			log.
				Err(err).
				Msg("Error sending response")

			return
		}
	} else {
		log.Info().Msg("No response url")
	}
}

func sendResponse(rd *responseData) error {
	log.
		Info().
		Str("Url", rd.ResponseUrl).
		Msg("Sending response")

	var r io.Reader
	request, err := http.NewRequest(rd.ResponseMethod, rd.ResponseUrl, r)
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
	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.
			Err(err).
			Msg("Error reading response body")

		return err
	}
	log.
		Info().
		Str("Url", rd.ResponseUrl).
		Msg("Response sent")

	return nil
}

func sendRequest(p *RequestData) (*responseData, error) {
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

	var res = &responseData{
		Header:         response.Header,
		Body:           body,
		ResponseUrl:    p.ResponseUrl,
		ResponseMethod: p.ResponseMethod,
	}

	return res, nil
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a *RequestData) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (p *RequestData) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return nil
	}

	err := json.Unmarshal(source, p)
	if err != nil {
		return err
	}

	return nil
}
