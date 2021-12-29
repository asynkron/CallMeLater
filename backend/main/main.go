package main

import (
	"github.com/asynkron/CallMeLater/server"
	storage2 "github.com/asynkron/CallMeLater/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// PSQL
	storage := storage2.NewPg("postgres://postgres:postgres@localhost:5432/callme?sslmode=disable")
	// SQLite
	// storage = newSqLiteStorage("file:storage.db?cache=shared&mode=memory")

	_ = server.New(storage)
}
