package server

import (
	"time"
)

type JobStorage interface {
	Pull(count int) ([]Job, error)
	Push(job Job) error
	Complete(job Job) error
	Cancel(job Job) error
	//GetResults(requestId string) ([]*JobResultEntity, error)
}

type JobEntity struct {
	Id                 string    `gorm:"primaryKey;type:varchar"`
	ScheduledTimestamp time.Time `gorm:"index"`
	CreatedTimestamp   time.Time
	DataDiscriminator  string
	Data               string
	RetryCount         int
	ParentJobId        string
	Status             JobStatus         `gorm:"index"`
	Results            []JobResultEntity `gorm:"foreignKey:JobId"`
}

func (JobEntity) TableName() string {
	return "jobs"
}

type JobStatus int64

const (
	Scheduled             JobStatus = 0
	CompletedSuccessfully           = 1
	Cancelled                       = 2
	Failed                          = 3
)

type JobResultEntity struct {
	Id                 string `gorm:"primaryKey;type:varchar"`
	JobId              string
	ExecutionTimestamp time.Time
	Status             string
	DataDiscriminator  string
	Data               string
}

func (JobResultEntity) TableName() string {
	return "job_results"
}
