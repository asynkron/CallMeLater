package main

import (
	"github.com/asynkron/CallMeLater/pg_storage"
	"github.com/asynkron/CallMeLater/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// PSQL
	storage := pg_storage.New("postgres://postgres:postgres@localhost:5432/callme?sslmode=disable")
	// SQLite
	// storage = newSqLiteStorage("file:storage.db?cache=shared&mode=memory")

	_ = callmelater.New(storage)
}
