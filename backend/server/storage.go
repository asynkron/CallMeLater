package server

import (
	"database/sql"
	"time"
)

type JobStorage interface {
	Pull(count int) ([]Job, error)
	Push(job Job) error
	Complete(job Job) error
	//GetResults(requestId string) ([]*JobResultEntity, error)
}

type JobEntity struct {
	Id                 string    `gorm:"primaryKey"`
	ScheduledTimestamp time.Time `gorm:"index"`
	CreatedTimestamp   time.Time
	CompletedTimestamp sql.NullTime `gorm:"index"`
	Data               string
	RetryCount         int
	ParentJobId        string
	Results            []JobResultEntity
}

type JobResultEntity struct {
	Id                 string `gorm:"primaryKey"`
	JobEntityId        string
	ExecutionTimestamp time.Time
	Status             string
	Data               string
}
