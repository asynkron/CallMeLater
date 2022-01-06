package server

import "time"

type JobEntity struct {
	Id                 string            `gorm:"primaryKey;type:varchar"`
	ScheduledTimestamp time.Time         `gorm:"index"`
	CreatedTimestamp   time.Time         `gorm:""`
	ExecutedTimestamp  time.Time         `gorm:""`
	ExecutedStatus     ExecutedStatus    `gorm:""`
	Description        string            `gorm:""`
	DataDiscriminator  string            `gorm:""`
	Data               string            `gorm:""`
	RetryCount         int               `gorm:""`
	ParentJobId        string            `gorm:""`
	CronExpression     string            `gorm:""`
	Status             ScheduledStatus   `gorm:"index"`
	Results            []JobResultEntity `gorm:"foreignKey:JobId"`
}

func (JobEntity) TableName() string {
	return "jobs"
}

type ScheduledStatus int

const (
	JobStatusScheduled ScheduledStatus = 0
	JobStatusSuccess                   = 1
	JobStatusCancelled                 = 2
	JobStatusFailed                    = 3
)

type ExecutedStatus int

const (
	ExecutedStatusPending ExecutedStatus = 0
	ExecutedStatusFail                   = 1
	ExecutedStatusSuccess                = 2
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
