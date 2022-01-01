package server

import (
	"github.com/rs/zerolog/log"
	"sort"
	"sync"
	"time"
)

type worker struct {
	storage        JobStorage
	hasMore        bool
	executableJobs chan Job
	sortedJobs     []Job
	pullCount      int
}

func New(storage JobStorage) *worker {
	w := &worker{
		storage:        storage,
		executableJobs: make(chan Job),
		pullCount:      100,
	}
	go w.run()
	handleRequests(w)
	return w
}

func (w *worker) run() {
	var err error
	w.sortedJobs, err = w.storage.Pull(w.pullCount)
	if err != nil {
		log.Err(err).Msg("failed to get sortedJobs executableJobs")
		return
	}

	log.Info().Msg("Consume loop started")

	for {
		select {
		case rd := <-w.executableJobs:
			w.appendJob(rd)
		case <-time.After(time.Second):
		}
		if len(w.sortedJobs) > 0 {
			first := w.sortedJobs[0]
			log.Info().Time("NextExecution", first.GetScheduledTimestamp()).Msg("Next job to be executed")
		}
		_ = w.executeJobs()
		_ = w.getMoreJobs()
	}
}

func (w *worker) appendJob(rd Job) {
	w.sortedJobs = append(w.sortedJobs, rd)

	sort.Slice(w.sortedJobs, func(i, j int) bool {
		w1 := w.sortedJobs[i].GetScheduledTimestamp()
		w2 := w.sortedJobs[j].GetScheduledTimestamp()
		return w1.Before(w2)
	})
}

func (w *worker) executeJobs() error {
	var wg sync.WaitGroup

	var count int
	for _, job := range w.sortedJobs {
		if job.GetScheduledTimestamp().Before(time.Now()) {
			wg.Add(1)
			go w.executeJob(job, &wg)
			w.sortedJobs = w.sortedJobs[1:]
			count++
		} else {
			break
		}
	}

	if count == 0 {
		return nil
	}

	log.Info().Msg("Waiting for executableJobs to finish")
	wg.Wait()
	log.Info().Msg("All executableJobs finished")

	if len(w.sortedJobs) > 100 {
		w.sortedJobs = w.sortedJobs[0:100]
		w.hasMore = true
	}

	return nil
}

func (w *worker) executeJob(job Job, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info().Str("Id", job.GetId()).Msg("Executing job")
	err := job.Execute(w.storage, w.executableJobs)
	if err != nil {
		log.Err(err).Msg("Error executing job")
	}
	log.Info().Str("Id", job.GetId()).Msg("Completed job")
	err = w.storage.Complete(job)
	if err != nil {
		log.Err(err).Msg("Error completing job")
	}
}

func (w *worker) getMoreJobs() error {
	if len(w.sortedJobs) > 0 {
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
	w.sortedJobs = pr
	w.hasMore = len(pr) > 0
	return nil
}
