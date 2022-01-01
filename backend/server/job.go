package server

import "time"

const (
	httpRequest  = "http_request"
	kafkaPublish = "kafka_publish"
)

type Job interface {
	Execute(storage JobStorage, pending chan Job) error
	GetScheduledTimestamp() time.Time
	GetId() string
	GetType() string
	InitDefaults()
}
