package main

import "github.com/rs/zerolog/log"

type RequestStorage interface {
	Get(id string) (*requestData, error)
	Set(id string, data *requestData) error
}

type NoopStorage struct{}

func (n NoopStorage) Get(id string) (*requestData, error) {
	log.
		Info().
		Str("id", id).
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
