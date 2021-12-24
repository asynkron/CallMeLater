package main

import (
	"github.com/rs/zerolog/log"
	"sort"
	"time"
)

func consumeLoop() {
	var pendingRequests []*requestData

	log.
		Info().
		Msg("Consume loop started")

	for {
		select {
		case rd := <-requests:
			log.
				Info().
				Msg("Worker received new request")

			pendingRequests = append(pendingRequests, rd)

			sort.Slice(pendingRequests, func(i, j int) bool {
				w1 := pendingRequests[i].When
				w2 := pendingRequests[j].When
				return w1.Before(w2)
			})

			pendingRequests = sendExpiredRequests(pendingRequests)

			if len(pendingRequests) > 100 {
				pendingRequests = pendingRequests[0:100]
			}

		case <-time.After(time.Second):
			log.
				Info().
				Msg("Worker received no new messages")

			pendingRequests = sendExpiredRequests(pendingRequests)
		}
	}
}

func sendExpiredRequests(pendingRequests []*requestData) []*requestData {
	for _, erd := range pendingRequests {
		if erd.When.Before(time.Now()) {
			go sendRequestResponse(erd)
			pendingRequests = pendingRequests[1:]
		} else {
			break
		}
	}
	return pendingRequests
}
