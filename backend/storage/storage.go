package storage

import (
	"database/sql"
	"encoding/json"
	"github.com/asynkron/CallMeLater/server"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Job struct {
	Id                 string `gorm:"primaryKey"`
	ScheduledTimestamp time.Time
	CreatedTimestamp   time.Time
	CompletedTimestamp sql.NullTime `gorm:"index"`
	Data               string
	RetryCount         int
	ParentJobId        string
	Results            []JobResult
}

type JobResult struct {
	Id                 string `gorm:"primaryKey"`
	JobId              string
	ExecutionTimestamp time.Time
	status             string
	Data               string
}

type GormStorage struct {
	db *gorm.DB
}

func (g GormStorage) Pull(count int) ([]*server.RequestData, error) {
	var jobs []Job
	err := g.db.Limit(count).Find(&jobs).Error
	if err != nil {
		return nil, err
	}

	var requests []*server.RequestData
	for _, job := range jobs {
		request := &server.RequestData{}
		var d = []byte(job.Data)
		err = json.Unmarshal(d, request)
		if err != nil {
			log.Err(err).Msg("Failed to unmarshal row")

			requests = append(requests, request)
		}
	}

	return requests, nil
}

func (g GormStorage) Push(data *server.RequestData) error {
	j, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Msg("Failed to marshal data")

		return err
	}

	job := &Job{
		Id:                 data.Id,
		ScheduledTimestamp: data.ScheduledTimestamp,
		CreatedTimestamp:   time.Now(),
		Data:               string(j),
	}

	g.db.Create(job)
	return nil
}

func (g GormStorage) Complete(requestId string) error {
	job := &Job{}
	err := g.db.Where("id = ?", requestId).First(job).Error
	if err != nil {
		return err
	}

	job.CompletedTimestamp = sql.NullTime{Time: time.Now(), Valid: true}
	g.db.Save(job)

	return nil
}

func NewPg(connectionString string) server.JobStorage {
	dialector := postgres.Open(connectionString)
	return New(dialector)
}

func New(dialector gorm.Dialector) *GormStorage {
	db, err := gorm.Open(dialector, &gorm.Config{})
	q := &GormStorage{db: db}
	if err != nil {
		log.Err(err).Msg("failed to connect to database")
	}
	err = db.AutoMigrate(&Job{})
	if err != nil {
		log.Err(err).Msg("failed to migrate database for Job")
	}
	err = db.AutoMigrate(&JobResult{})
	if err != nil {
		log.Err(err).Msg("failed to migrate database for JobResult")
	}
	return q
}
