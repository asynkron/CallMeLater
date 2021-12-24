package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	storage  RequestStorage = &NoopStorage{} // NewPgStorage("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	requests                = make(chan *requestData)
)

func consumeLoop() {
	log.
		Info().
		Msg("Consume loop started")

	for {
		rd := <-requests
		go sendRequestResponse(rd)
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	go consumeLoop()
	handleRequests()
}
