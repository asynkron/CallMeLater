package main

import (
	"github.com/rs/zerolog/log"
	"sort"
	"time"
)

// TODO: make this all less hacky
func consumeLoop() {
	pendingRequests, err := storage.Get()
	if err != nil {
		log.
			Err(err).
			Msg("failed to get pending requests")
		return
	}

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

			pendingRequests, err = sendExpiredRequests(pendingRequests)
			if err != nil {
				log.
					Err(err).
					Msg("failed to send expired requests")
				return
			}

			if len(pendingRequests) > 100 {
				pendingRequests = pendingRequests[0:100]
				hasMore = true
			}

		case <-time.After(time.Second):
			log.
				Info().
				Msg("Worker received no new messages")

			pendingRequests, err = sendExpiredRequests(pendingRequests)
			if err != nil {
				log.
					Err(err).
					Msg("failed to send expired requests")
				return
			}
		}
	}
}

func sendExpiredRequests(pendingRequests []*requestData) ([]*requestData, error) {
	for _, erd := range pendingRequests {
		if erd.When.Before(time.Now()) {
			//delete the request from the DB.
			err := storage.Delete(erd.RequestId)
			if err != nil {
				log.
					Err(err).
					Msg("Error deleting request")
			}

			go sendRequestResponse(erd)
			pendingRequests = pendingRequests[1:]
		} else {
			break
		}
	}

	pendingRequests, err := loadMore(pendingRequests)
	if err != nil {
		return nil, err
	}

	return pendingRequests, nil
}

func loadMore(pendingRequests []*requestData) ([]*requestData, error) {
	if len(pendingRequests) > 0 {
		return pendingRequests, nil
	}

	if !hasMore {
		return pendingRequests, nil
	}

	log.
		Info().
		Msg("Loading more")

	pr, err := storage.Get()
	if err != nil {
		return nil, err
	}
	hasMore = len(pr) > 0
	return pr, nil
}
