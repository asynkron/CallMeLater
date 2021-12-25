package server

import "github.com/rs/zerolog/log"

type RequestStorage interface {
	Get() ([]*RequestData, error)
	Set(data *RequestData) error
	Delete(requestId string) error
}

type NoopStorage struct{}

func (n NoopStorage) Get() ([]*RequestData, error) {
	log.
		Info().
		Msg("NoopStorage.Get")
	return nil, nil
}

func (n NoopStorage) Set(data *RequestData) error {
	log.
		Info().
		Str("id", data.RequestId).
		Msg("NoopStorage.Set")

	return nil
}

func (n NoopStorage) Delete(requestId string) error {
	log.
		Info().
		Str("id", requestId).
		Msg("NoopStorage.Delete")

	return nil
}
