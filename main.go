package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	storage  RequestStorage
	requests = make(chan *requestData)
	hasMore  = false
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// PSQL
	storage = newPgStorage("postgres://postgres:postgres@localhost:5432/callme?sslmode=disable")
	// SQLite
	storage = newSqLiteStorage("file:storage.db?cache=shared&mode=memory")

	go consumeLoop()
	handleRequests()
}
