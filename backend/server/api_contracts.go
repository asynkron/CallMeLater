package server

import "time"

type GetJobsResponse struct {
	Skip  int           `json:"skip"`
	Limit int           `json:"limit"`
	Count int           `json:"count"`
	Jobs  []JobResponse `json:"jobs"`
}

type JobResponse struct {
	Id                 string    `json:"id"`
	ScheduledTimestamp time.Time `json:"scheduledTimestamp"`
	CreatedTimestamp   time.Time `json:"createdTimestamp"`
	DataDiscriminator  string    `json:"dataDiscriminator"`
	ParentJobId        string    `json:"parentJobId"`
	Status             JobStatus `json:"status"`
}
