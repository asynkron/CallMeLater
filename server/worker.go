package callmelater

import (
	"github.com/rs/zerolog/log"
	"sort"
	"time"
)

type worker struct {
	storage  RequestStorage
	hasMore  bool
	requests chan *RequestData
	pending  []*RequestData
}

func New(storage RequestStorage) *worker {
	w := &worker{
		storage:  storage,
		requests: make(chan *RequestData),
	}
	go w.run()
	handleRequests(w)
	return w
}

// TODO: make this all less hacky
func (w *worker) run() {
	var err error
	w.pending, err = w.storage.Get()
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

			w.pending = append(w.pending, rd)

			sort.Slice(w.pending, func(i, j int) bool {
				w1 := w.pending[i].When
				w2 := w.pending[j].When
				return w1.Before(w2)
			})

			log.
				Info().
				Int("QueueLength", len(w.pending)).
				Msg("Worker received new request")

			err = w.sendExpiredRequests()
			if err != nil {
				log.
					Err(err).
					Msg("failed to send expired requests")
				return
			}

			if len(w.pending) > 100 {
				w.pending = w.pending[0:100]
				w.hasMore = true
			}

		case <-time.After(time.Second):
			log.
				Info().
				Int("QueueLength", len(w.pending)).
				Msg("Worker received no new messages")

			err = w.sendExpiredRequests()
			if err != nil {
				log.
					Err(err).
					Msg("failed to send expired requests")
				return
			}
		}
	}
}

func (w *worker) sendExpiredRequests() error {
	for _, erd := range w.pending {
		if erd.When.Before(time.Now()) {
			//delete the request from the DB.
			err := w.storage.Delete(erd.RequestId)
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

	err := w.loadMore()
	if err != nil {
		return err
	}

	return nil
}

func (w *worker) loadMore() error {
	if len(w.pending) > 0 {
		return nil
	}

	if !w.hasMore {
		return nil
	}

	log.
		Info().
		Msg("Loading more")

	pr, err := w.storage.Get()
	if err != nil {
		return err
	}
	w.pending = pr
	w.hasMore = len(pr) > 0
	return nil
}
