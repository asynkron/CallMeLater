package server

import (
	"bytes"
	"context"
	"fmt"
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
	RetryCount         int                 `json:"retry_count,omitempty"`
}

func (job *HttpRequestJob) GetScheduledTimestamp() time.Time {
	return job.ScheduledTimestamp
}

func (job *HttpRequestJob) GetId() string {
	return job.Id
}

func (job *HttpRequestJob) Execute(storage JobStorage, expired chan Job) error {
	response, err := sendRequest(job)

	if err != nil {
		log.Err(err).Msg("Error sending request")
		job.RetryCount++
		if job.RetryCount > 10 {
			err = storage.Fail(job)
			if err != nil {
				log.Err(err).Msg("Error marking job as failed")
				return err
			}
		} else {
			//todo: define backoff strategy
			job.ScheduledTimestamp = job.ScheduledTimestamp.Add(time.Duration(job.RetryCount) * time.Second)
			err = storage.Update(job)
			if err != nil {
				log.Err(err).Msg("Error updating job")
				return err
			}
			scheduleJob(expired, job)
		}
		return err
	}

	if job.ResponseUrl != "" {
		log.Info().Str("Id", response.Id).Str("Url", response.RequestUrl).Msg("Response Job created")
		_ = storage.Create(response)
		log.Info().Str("Id", response.Id).Str("Url", response.RequestUrl).Msg("Response Job stored")
		scheduleJob(expired, response)
	}
	return nil
}

func scheduleJob(expired chan Job, job *HttpRequestJob) {
	go func() { expired <- job }()
}

func sendRequest(job *HttpRequestJob) (*HttpRequestJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Info().Str("Id", job.Id).Str("Url", job.RequestUrl).Msg("Sending request")

	var r io.Reader
	request, err := http.NewRequestWithContext(ctx, job.RequestMethod, job.RequestUrl, r)
	if err != nil {
		return nil, err
	}
	request.Header = job.Header
	request.Form = job.Form
	request.Body = ioutil.NopCloser(bytes.NewReader(job.Body))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("request failed with status code %d", response.StatusCode)
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

func (job *HttpRequestJob) GetType() string {
	return httpRequest
}
