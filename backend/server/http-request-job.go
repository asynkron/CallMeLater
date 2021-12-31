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

type HttpRequestJob struct {
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

func (h *HttpRequestJob) GetScheduledTimestamp() time.Time {
	return h.ScheduledTimestamp
}

func (h *HttpRequestJob) GetId() string {
	return h.Id
}

func (h *HttpRequestJob) Execute(storage JobStorage, expired chan Job) {
	response, err := sendRequest(h)

	if err != nil {
		log.Err(err).Msg("Error sending request")
		return
	}

	if h.ResponseUrl != "" {
		log.Info().Str("Id", response.Id).Str("Url", response.RequestUrl).Msg("Response Job created")
		_ = storage.Push(response)
		expired <- response
	}
}

func sendRequest(job *HttpRequestJob) (*HttpRequestJob, error) {
	log.Info().Str("Id", job.Id).Str("Url", job.RequestUrl).Msg("Sending request")

	var r io.Reader
	request, err := http.NewRequest(job.RequestMethod, job.RequestUrl, r)
	if err != nil {
		return nil, err
	}
	request.Header = job.Header
	request.Form = job.Form
	request.Body = ioutil.NopCloser(bytes.NewReader(job.Body))
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

	var res = &HttpRequestJob{
		Id:                 uuid.New().String(),
		Header:             response.Header,
		Body:               body,
		RequestUrl:         job.ResponseUrl,
		RequestMethod:      job.ResponseMethod,
		ScheduledTimestamp: job.ScheduledTimestamp,
		ParentId:           job.Id,
	}

	return res, nil
}