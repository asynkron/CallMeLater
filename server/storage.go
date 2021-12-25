package server

import (
	"github.com/rs/zerolog/log"
	"time"
)

type Job struct {
	Id        string
	Timestamp time.Time
	Data      []byte
}

type RequestStorage interface {
	Pull(count int) ([]*RequestData, error)
	Push(data *RequestData) error
	Complete(requestId string) error
}

type NoopStorage struct{}

func (n NoopStorage) Pull(int) ([]*RequestData, error) {
	log.
		Info().
		Msg("NoopStorage.Pull")

	return nil, nil
}

func (n NoopStorage) Push(data *RequestData) error {
	log.
		Info().
		Str("id", data.RequestId).
		Msg("NoopStorage.Push")

	return nil
}

func (n NoopStorage) Complete(requestId string) error {
	log.
		Info().
		Str("id", requestId).
		Msg("NoopStorage.Complete")

	return nil
}
