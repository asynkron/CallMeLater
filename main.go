package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"sort"
)

var (
	storage  RequestStorage
	requests = make(chan *requestData)
)

func consumeLoop() {
	var pendingRequests []*requestData

	log.
		Info().
		Msg("Consume loop started")

	for {
		rd := <-requests

		pendingRequests = append(pendingRequests, rd)

		sort.Slice(pendingRequests, func(i, j int) bool {
			w1 := pendingRequests[i].When
			w2 := pendingRequests[j].When
			return w1.Before(w2)
		})

		pendingRequests = pendingRequests[0:100]
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	storage = newPgStorage("postgres://postgres:postgres@localhost:5432/callme?sslmode=disable")
	go consumeLoop()
	handleRequests()
}
