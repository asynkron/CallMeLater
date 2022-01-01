package server

import "time"

type JobEntity struct {
	Id                 string            `gorm:"primaryKey;type:varchar"`
	ScheduledTimestamp time.Time         `gorm:"index"`
	CreatedTimestamp   time.Time         `gorm:""`
	DataDiscriminator  string            `gorm:""`
	Data               string            `gorm:""`
	RetryCount         int               `gorm:""`
	ParentJobId        string            `gorm:""`
	Status             JobStatus         `gorm:"index"`
	Results            []JobResultEntity `gorm:"foreignKey:JobId"`
}

func (JobEntity) TableName() string {
	return "executableJobs"
}

type JobStatus int

const (
	JobStatusScheduled             JobStatus = 0
	JobStatusCompletedSuccessfully           = 1
	JobStatusCancelled                       = 2
	JobStatusFailed                          = 3
)

type JobResultEntity struct {
	Id                 string    `gorm:"primaryKey;type:varchar"`
	JobId              string    `gorm:""`
	ExecutionTimestamp time.Time `gorm:""`
	Status             string    `gorm:""`
	DataDiscriminator  string    `gorm:""`
	Data               string    `gorm:""`
}

func (JobResultEntity) TableName() string {
	return "job_results"
}
