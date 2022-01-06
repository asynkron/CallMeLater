package server

import "time"

type GetJobsResponse struct {
	Skip  int           `json:"skip"`
	Limit int           `json:"limit"`
	Count int           `json:"count"`
	Jobs  []JobResponse `json:"jobs"`
}

type JobResponse struct {
	Id                 string          `json:"id"`
	ScheduledTimestamp time.Time       `json:"scheduledTimestamp"`
	Description        string          `json:"description"`
	DataDiscriminator  string          `json:"dataDiscriminator"`
	ParentJobId        string          `json:"parentJobId"`
	Status             ScheduledStatus `json:"status"`
	CronExpression     string          `json:"cronExpression"`
	ExecutedTimestamp  time.Time       `json:"executedTimestamp"`
	ExecutedStatus     ExecutedStatus  `json:"executedStatus"`
}
