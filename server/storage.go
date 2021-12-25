package callmelater

import "github.com/rs/zerolog/log"

type RequestStorage interface {
	Get() ([]*requestData, error)
	Set(data *requestData) error
	Delete(requestId string) error
}

type NoopStorage struct{}

func (n NoopStorage) Get() ([]*requestData, error) {
	log.
		Info().
		Msg("NoopStorage.Get")
	return nil, nil
}

func (n NoopStorage) Set(data *requestData) error {
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
