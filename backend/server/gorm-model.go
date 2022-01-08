package server

import "time"

type JobEntity struct {
	Id                     string            `gorm:"primaryKey;type:varchar"`
	ParentJobId            string            `gorm:""`
	CreatedTimestamp       time.Time         `gorm:""`
	Description            string            `gorm:""`
	ScheduleTimestamp      *time.Time        `gorm:"index"`
	ScheduleStatus         ScheduleStatus    `gorm:"index"`
	ScheduleCronExpression string            `gorm:""`
	ExecutedTimestamp      *time.Time        `gorm:""`
	ExecutedStatus         ExecutedStatus    `gorm:""`
	ExecutedCount          int               `gorm:""`
	DataDiscriminator      string            `gorm:""`
	Data                   string            `gorm:""`
	Results                []JobResultEntity `gorm:"foreignKey:JobId"`
}

func (JobEntity) TableName() string {
	return "jobs"
}

type ScheduleStatus int

const (
	JobStatusScheduled ScheduleStatus = 0
	JobStatusSuccess                  = 1
	JobStatusCancelled                = 2
	JobStatusFailed                   = 3
)

type ExecutedStatus int

const (
	ExecutedStatusPending ExecutedStatus = 0
	ExecutedStatusFail                   = 1
	ExecutedStatusSuccess                = 2
)

type JobResultEntity struct {
	Id                 string     `gorm:"primaryKey;type:varchar"`
	JobId              string     `gorm:""`
	ExecutionTimestamp *time.Time `gorm:""`
	Status             string     `gorm:""`
	DataDiscriminator  string     `gorm:""`
	Data               string     `gorm:""`
}

func (JobResultEntity) TableName() string {
	return "job_results"
}
