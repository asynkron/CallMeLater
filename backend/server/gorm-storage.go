package server

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type GormStorage struct {
	db *gorm.DB
}

func (g GormStorage) Cancel(job Job) error {
	jobEntity := &JobEntity{}
	err := g.db.
		Where("id = ?", job.GetId()).
		First(jobEntity).Error

	if err != nil {
		return err
	}

	jobEntity.Status = JobStatusCancelled
	g.db.Save(jobEntity)

	return nil
}

func (g GormStorage) Fail(job Job) error {
	jobEntity := &JobEntity{}
	err := g.db.
		Where("id = ?", job.GetId()).
		First(jobEntity).Error

	if err != nil {
		return err
	}

	jobEntity.Status = JobStatusFailed
	g.db.Save(jobEntity)

	return nil
}

func (g GormStorage) Pull(count int) ([]Job, error) {
	var jobs []JobEntity
	err := g.db.
		Limit(count).
		Order("scheduled_timestamp asc").
		Find(&jobs, "status = ?", JobStatusScheduled).Error

	if err != nil {
		return nil, err
	}

	var requests []Job
	for _, job := range jobs {
		var unmarshalledJob = instanceFromDiscriminator(job.DataDiscriminator)
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

func instanceFromDiscriminator(discriminator string) Job {
	var unmarshalledJob Job

	switch discriminator {
	case httpRequest:
		unmarshalledJob = &HttpRequestJob{}
	case kafkaPublish:
		//unmarshalledJob = &KafkaPublishJob{}
	}
	return unmarshalledJob
}

func (g GormStorage) Create(job Job) error {
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

func (g GormStorage) Update(job Job) error {
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

	g.db.Save(jobEntity)

	return nil
}

func (g GormStorage) Complete(job Job) error {
	jobEntity := &JobEntity{}
	err := g.db.
		Where("id = ?", job.GetId()).
		First(jobEntity).Error

	if err != nil {
		return err
	}

	jobEntity.Status = JobStatusCompletedSuccessfully
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
