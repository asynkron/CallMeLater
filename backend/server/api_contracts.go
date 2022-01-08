package server

import "time"

type GetJobsResponse struct {
	Skip  int           `json:"skip"`
	Limit int           `json:"limit"`
	Count int           `json:"count"`
	Jobs  []JobResponse `json:"jobs"`
}

type JobResponse struct {
	Id                     string         `json:"id"`
	ParentJobId            string         `json:"parentJobId"`
	Description            string         `json:"description"`
	ScheduleTimestamp      *time.Time     `json:"scheduleTimestamp"`
	ScheduleCronExpression string         `json:"scheduleCronExpression"`
	DataDiscriminator      string         `json:"dataDiscriminator"`
	ExecutedTimestamp      *time.Time     `json:"executedTimestamp"`
	ExecutedStatus         ExecutedStatus `json:"executedStatus"`
	ExecutedCount          int            `json:"executedCount"`
}
