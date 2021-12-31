package server

import "time"

type Job interface {
	Execute(storage JobStorage, pending chan Job) error
	GetScheduledTimestamp() time.Time
	GetId() string
}
