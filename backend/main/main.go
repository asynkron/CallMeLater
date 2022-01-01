package main

import (
	"github.com/asynkron/CallMeLater/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"os"
)

func main() {
	//zerolog.TimeFieldFormat = zerolog.tim
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	dialector := postgres.Open("postgres://postgres:postgres@localhost:5432/callme?sslmode=disable")
	storage := server.NewStorage(dialector)
	_ = server.New(storage)
}
