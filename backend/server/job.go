package server

import "time"

type Job interface {
	Execute(storage JobStorage, pending chan Job)
	GetScheduledTimestamp() time.Time
	GetId() string
}
