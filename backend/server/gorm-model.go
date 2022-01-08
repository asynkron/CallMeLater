package server

import "time"

type JobEntity struct {
	Id                     string            `gorm:"primaryKey;type:varchar"`
	ParentJobId            string            `gorm:""`
	CreatedTimestamp       time.Time         `gorm:""`
	Description            string            `gorm:""`
	ScheduleTimestamp      *time.Time        `gorm:"index"`
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

type ExecutedStatus int

const (
	ExecutedStatusPending ExecutedStatus = 0
	ExecutedStatusSuccess                = 1
	ExecutedStatusFail                   = 2
	ExecutedStatusRetry                  = 3
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
