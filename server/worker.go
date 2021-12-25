package server

import (
	"github.com/rs/zerolog/log"
	"sort"
	"time"
)

type worker struct {
	storage   RequestStorage
	hasMore   bool
	requests  chan *RequestData
	pending   []*RequestData
	pullCount int
}

func New(storage RequestStorage) *worker {
	w := &worker{
		storage:   storage,
		requests:  make(chan *RequestData),
		pullCount: 100,
	}
	go w.run()
	handleRequests(w)
	return w
}

// TODO: make this all less hacky
func (w *worker) run() {
	var err error
	w.pending, err = w.storage.Pull(w.pullCount)
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
		case rd := <-w.requests:

			w.appendRequest(rd)

			log.
				Info().
				Int("QueueLength", len(w.pending)).
				Msg("Worker received new request")

		case <-time.After(time.Second):
			log.
				Info().
				Int("QueueLength", len(w.pending)).
				Msg("Worker received no new messages")
		}
		_ = w.sendExpiredRequests()
		_ = w.loadMoreRequests()
	}
}

func (w *worker) appendRequest(rd *RequestData) {
	w.pending = append(w.pending, rd)

	sort.Slice(w.pending, func(i, j int) bool {
		w1 := w.pending[i].When
		w2 := w.pending[j].When
		return w1.Before(w2)
	})
}

func (w *worker) sendExpiredRequests() error {
	for _, erd := range w.pending {
		if erd.When.Before(time.Now()) {
			//delete the request from the DB.
			err := w.storage.Complete(erd.RequestId)
			if err != nil {
				log.
					Err(err).
					Msg("Error deleting request")
			}

			go sendRequestResponse(erd)
			w.pending = w.pending[1:]
		} else {
			break
		}
	}

	if len(w.pending) > 100 {
		w.pending = w.pending[0:100]
		w.hasMore = true
	}

	return nil
}

func (w *worker) loadMoreRequests() error {
	if len(w.pending) > 0 {
		return nil
	}

	if !w.hasMore {
		return nil
	}

	log.
		Info().
		Msg("Loading more")

	pr, err := w.storage.Pull(w.pullCount)
	if err != nil {
		if err != nil {
			log.
				Err(err).
				Msg("Error loading more")
		}
		return err
	}
	w.pending = pr
	w.hasMore = len(pr) > 0
	return nil
}
