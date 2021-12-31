package server

import (
	"github.com/rs/zerolog/log"
	"sort"
	"time"
)

type worker struct {
	storage   JobStorage
	hasMore   bool
	expired   chan Job
	cache     []Job
	pullCount int
}

func New(storage JobStorage) *worker {
	w := &worker{
		storage:   storage,
		expired:   make(chan Job),
		pullCount: 100,
	}
	go w.run()
	handleRequests(w)
	return w
}

// TODO: make this all less hacky
func (w *worker) run() {
	var err error
	w.cache, err = w.storage.Pull(w.pullCount)
	if err != nil {
		log.Err(err).Msg("failed to get cache expired")
		return
	}

	log.Info().Msg("Consume loop started")

	for {
		select {
		case rd := <-w.expired:

			w.appendRequest(rd)

			log.Info().Int("QueueLength", len(w.cache)).Msg("Worker received new request")

		case <-time.After(time.Second):
			log.Info().Int("QueueLength", len(w.cache)).Msg("Worker received no new messages")
		}
		_ = w.executeExpiredJobs()
		_ = w.loadMoreRequests()
	}
}

func (w *worker) appendRequest(rd Job) {
	w.cache = append(w.cache, rd)

	sort.Slice(w.cache, func(i, j int) bool {
		w1 := w.cache[i].GetScheduledTimestamp()
		w2 := w.cache[j].GetScheduledTimestamp()
		return w1.Before(w2)
	})
}

func (w *worker) executeExpiredJobs() error {
	for _, job := range w.cache {
		if job.GetScheduledTimestamp().Before(time.Now()) {
			//delete the request from the DB.
			err := w.storage.Complete(job)
			if err != nil {
				log.Err(err).Msg("Error completing request")
			}

			go job.Execute(w.storage, w.expired)
			w.cache = w.cache[1:]
		} else {
			break
		}
	}

	if len(w.cache) > 100 {
		w.cache = w.cache[0:100]
		w.hasMore = true
	}

	return nil
}

func (w *worker) loadMoreRequests() error {
	if len(w.cache) > 0 {
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
	w.cache = pr
	w.hasMore = len(pr) > 0
	return nil
}
