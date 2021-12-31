package server

import (
	"github.com/rs/zerolog/log"
	"sort"
	"time"
)

type worker struct {
	storage   JobStorage
	hasMore   bool
	requests  chan Job
	pending   []Job
	pullCount int
}

func New(storage JobStorage) *worker {
	w := &worker{
		storage:   storage,
		requests:  make(chan Job),
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
		log.Err(err).Msg("failed to get pending requests")
		return
	}

	log.Info().Msg("Consume loop started")

	for {
		select {
		case rd := <-w.requests:

			w.appendRequest(rd)

			log.Info().Int("QueueLength", len(w.pending)).Msg("Worker received new request")

		case <-time.After(time.Second):
			log.Info().Int("QueueLength", len(w.pending)).Msg("Worker received no new messages")
		}
		_ = w.executeExpiredJobs()
		_ = w.loadMoreRequests()
	}
}

func (w *worker) appendRequest(rd Job) {
	w.pending = append(w.pending, rd)

	sort.Slice(w.pending, func(i, j int) bool {
		w1 := w.pending[i].GetScheduledTimestamp()
		w2 := w.pending[j].GetScheduledTimestamp()
		return w1.Before(w2)
	})
}

func (w *worker) executeExpiredJobs() error {
	for _, job := range w.pending {
		if job.GetScheduledTimestamp().Before(time.Now()) {
			//delete the request from the DB.
			err := w.storage.Complete(job)
			if err != nil {
				log.Err(err).Msg("Error completing request")
			}

			go job.Execute(w.storage, w.requests)
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

	log.Info().Msg("Loading more")

	pr, err := w.storage.Pull(w.pullCount)
	if err != nil {
		if err != nil {
			log.Err(err).Msg("Error loading more")
		}
		return err
	}
	w.pending = pr
	w.hasMore = len(pr) > 0
	return nil
}
