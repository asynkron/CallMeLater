package main

import "github.com/rs/zerolog/log"

type RequestStorage interface {
	Get() ([]*requestData, error)
	Set(id string, data *requestData) error
}

type NoopStorage struct{}

func (n NoopStorage) Get() ([]*requestData, error) {
	log.
		Info().
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
