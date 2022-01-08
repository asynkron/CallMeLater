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
	*JobEntity     `json:"-"`
	RequestMethod  string              `json:"request_method,omitempty"`
	Header         map[string][]string `json:"header,omitempty"`
	Form           map[string][]string `json:"form,omitempty"`
	RequestUrl     string              `json:"request_url,omitempty"`
	ResponseUrl    string              `json:"response_url,omitempty"`
	ResponseMethod string              `json:"response_method,omitempty"`
	Body           []byte              `json:"body,omitempty"`
	RetryCount     int                 `json:"retry_count,omitempty"`
	RetryMax       int                 `json:"retry_max,omitempty"`
	RetryDelay     time.Duration       `json:"retry_delay,omitempty"`
}

func (job *HttpRequestJob) GetEntity() *JobEntity {
	return job.JobEntity
}

func (job *HttpRequestJob) DiagnosticsString() string {
	return job.Id + " " + job.RequestMethod + " " + job.RequestUrl
}

func (job *HttpRequestJob) ShouldRetry() bool {
	return job.RetryCount < job.RetryMax
}

func (job *HttpRequestJob) InitDefaults() {
	if job.RetryMax == 0 {
		job.RetryMax = 3
	}

	if job.RetryDelay == 0 {
		job.RetryDelay = time.Minute * 1
	}

	if job.JobEntity.Description == "" {
		job.JobEntity.Description = "HTTP " + job.RequestMethod + " " + job.RequestUrl
	}
}

func (job *HttpRequestJob) GetScheduledTimestamp() *time.Time {
	return job.ScheduleTimestamp
}

func (job *HttpRequestJob) GetId() string {
	return job.Id
}

func (job *HttpRequestJob) Execute(storage JobStorage, expired chan Job) error {
	response, err := send(job)

	if err != nil {
		return err
	}

	if job.ResponseUrl != "" {
		job.respond(storage, expired, response)
	}
	return nil
}

func (job *HttpRequestJob) respond(storage JobStorage, expired chan Job, response *HttpRequestJob) {
	_ = storage.Create(response)
	log.Info().Str("Job", response.DiagnosticsString()).Msg("Response Job created")
	schedule(expired, response)
}

func (job *HttpRequestJob) Retry(storage JobStorage, expired chan Job) error {
	//todo: define backoff strategy
	job.RetryCount++

	job.ScheduleTimestamp = timeToPtr(time.Now().Add(time.Duration(job.RetryCount) * job.RetryDelay))
	err := storage.Retry(job)
	if err != nil {
		log.Err(err).Str("Job", job.DiagnosticsString()).Msg("Error updating job")
		return err
	}
	schedule(expired, job)
	return nil
}

func (job *HttpRequestJob) Fail(storage JobStorage, _ chan Job) error {
	err := storage.Fail(job)
	if err != nil {
		log.Err(err).Str("Job", job.DiagnosticsString()).Msg("Error marking job as failed")
		return err
	}
	return nil
}

func schedule(expired chan Job, job *HttpRequestJob) {
	go func() { expired <- job }()
}

func send(job *HttpRequestJob) (*HttpRequestJob, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	log.Info().Str("Job", job.DiagnosticsString()).Msg("Sending request")

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
		Header:        response.Header,
		Body:          body,
		RequestUrl:    job.ResponseUrl,
		RequestMethod: job.ResponseMethod,
		JobEntity:     newHttpRequestJobEntity(job.ScheduleTimestamp, job.Id),
	}
	res.InitDefaults()

	return res, nil
}

func newHttpRequestJobEntity(scheduledTimestamp *time.Time, parentJobId string) *JobEntity {
	return &JobEntity{
		Id:                uuid.New().String(),
		ScheduleTimestamp: scheduledTimestamp,
		ParentJobId:       parentJobId,
		CreatedTimestamp:  time.Now(),
		DataDiscriminator: httpRequest,
	}
}
