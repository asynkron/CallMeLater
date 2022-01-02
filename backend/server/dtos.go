package server

import "time"

type JobDto struct {
	Id                 string    `json:"id"`
	ScheduledTimestamp time.Time `json:"scheduledTimestamp"`
	CreatedTimestamp   time.Time `json:"createdTimestamp"`
	DataDiscriminator  string    `json:"dataDiscriminator"`
	ParentJobId        string    `json:"parentJobId"`
	Status             JobStatus `json:"status"`
}
