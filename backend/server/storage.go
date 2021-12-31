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
	Id                 string    `gorm:"primaryKey;type:varchar"`
	ScheduledTimestamp time.Time `gorm:"index"`
	CreatedTimestamp   time.Time
	CompletedTimestamp sql.NullTime `gorm:"index"`
	DataDiscriminator  string
	Data               string
	RetryCount         int
	ParentJobId        string
	Results            []JobResultEntity `gorm:"foreignKey:JobId"`
}

func (JobEntity) TableName() string {
	return "jobs"
}

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
