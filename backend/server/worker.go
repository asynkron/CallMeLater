package server

import (
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"sort"
	"sync"
	"time"
)

var (
	cronParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
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
			dur := first.GetScheduledTimestamp().Sub(time.Now()).String()

			log.Info().Str("In", dur).Time("Scheduled", first.GetScheduledTimestamp()).Msg("Next job to be executed")
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

	log.Info().Msg("Worker waiting for jobs to finish")
	wg.Wait()
	log.Info().Msg("Worker finished waiting for jobs")

	if len(w.sortedJobs) > 100 {
		w.sortedJobs = w.sortedJobs[0:100]
		w.hasMore = true
	}

	return nil
}

func (w *worker) executeJob(job Job, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Info().Str("Job", job.DiagnosticsString()).Msg("Executing job")
	err := job.Execute(w.storage, w.executableJobs)
	if err != nil {
		log.Err(err).Str("Job", job.DiagnosticsString()).Msg("Failed to execute job")
		if job.ShouldRetry() {
			log.Warn().Str("Job", job.DiagnosticsString()).Time("Scheduled", job.GetScheduledTimestamp()).Msg("Retrying job")
			_ = job.Retry(w.storage, w.executableJobs)
		} else {
			log.Warn().Str("Job", job.DiagnosticsString()).Msg("Marking job as failed")
			_ = job.Fail(w.storage, w.executableJobs)
		}
		return
	}
	log.Info().Str("Job", job.DiagnosticsString()).Msg("Completed job")

	if job.GetEntity().ScheduleCronExpression != "" {
		err = w.storage.RescheduleCron(job)
		if err != nil {
			log.Err(err).Msg("Error rescheduling cron job")
		}
		return
	}

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
