package server

import (
	"database/sql"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type GormStorage struct {
	db *gorm.DB
}

func (g GormStorage) Pull(count int) ([]Job, error) {
	var jobs []JobEntity
	err := g.db.Limit(count).Order("scheduled_timestamp asc").Find(&jobs, "completed_timestamp is null").Error
	if err != nil {
		return nil, err
	}

	var requests []Job
	for _, job := range jobs {
		var unmarshalledJob Job

		switch job.DataDiscriminator {
		case httpRequest:
			unmarshalledJob = &HttpRequestJob{}
		case kafkaPublish:
			//unmarshalledJob = &KafkaPublishJob{}
		}

		var d = []byte(job.Data)
		err = json.Unmarshal(d, unmarshalledJob)
		if err != nil {
			log.Err(err).Msg("Failed to unmarshal row")

			continue
		}
		requests = append(requests, unmarshalledJob)
	}

	return requests, nil
}

func (g GormStorage) Push(job Job) error {
	j, err := json.Marshal(job)
	if err != nil {
		log.Err(err).Msg("Failed to marshal job")

		return err
	}

	jobEntity := &JobEntity{
		Id:                 job.GetId(),
		ScheduledTimestamp: job.GetScheduledTimestamp(),
		CreatedTimestamp:   time.Now(),
		DataDiscriminator:  job.GetType(),
		Data:               string(j),
	}

	g.db.Create(jobEntity)
	return nil
}

func (g GormStorage) Complete(job Job) error {
	log.Info().Str("Id", job.GetId()).Msg("Completing Job")
	jobEntity := &JobEntity{}
	err := g.db.Where("id = ?", job.GetId()).First(jobEntity).Error
	if err != nil {
		return err
	}

	jobEntity.CompletedTimestamp = sql.NullTime{Time: time.Now(), Valid: true}
	g.db.Save(jobEntity)

	return nil
}

func NewStorage(dialector gorm.Dialector) *GormStorage {
	db, err := gorm.Open(dialector, &gorm.Config{})
	q := &GormStorage{db: db}
	if err != nil {
		log.Err(err).Msg("failed to connect to database")
	}
	err = db.AutoMigrate(&JobEntity{}, &JobResultEntity{})
	if err != nil {
		log.Err(err).Msg("failed to migrate database")
	}
	return q
}
