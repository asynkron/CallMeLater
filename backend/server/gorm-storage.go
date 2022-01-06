package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type GormStorage struct {
	db *gorm.DB
}

func (g GormStorage) Cancel(job Job) error {
	jobEntity := job.GetEntity()

	jobEntity.Status = JobStatusCancelled
	g.db.Save(jobEntity)

	return nil
}

func (g GormStorage) Fail(job Job) error {
	jobEntity := job.GetEntity()

	jobEntity.Status = JobStatusFailed

	result := newJobResultEntity(jobEntity)
	result.Status = "failed"
	result.Data = "somejson"
	jobEntity.Results = append(jobEntity.Results, result)
	g.db.Save(jobEntity)

	return nil
}

func newJobResultEntity(jobEntity *JobEntity) JobResultEntity {
	result := JobResultEntity{
		Id:                 uuid.New().String(),
		JobId:              jobEntity.Id,
		ExecutionTimestamp: jobEntity.ScheduledTimestamp,
		DataDiscriminator:  jobEntity.DataDiscriminator,
	}
	return result
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
		var unmarshalledJob = instanceFromDiscriminator(job)
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

func instanceFromDiscriminator(job JobEntity) Job {
	var unmarshalledJob Job

	switch job.DataDiscriminator {
	case httpRequest:
		unmarshalledJob = &HttpRequestJob{
			JobEntity: &job,
		}
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

	jobEntity := job.GetEntity()
	jobEntity.Data = string(j)

	g.db.Create(jobEntity)

	return nil
}

func (g GormStorage) Retry(job Job) error {
	j, err := json.Marshal(job)
	if err != nil {
		log.Err(err).Msg("Failed to marshal job")

		return err
	}

	jobEntity := job.GetEntity()
	jobEntity.Data = string(j)
	jobEntity.Status = JobStatusScheduled

	result := newJobResultEntity(jobEntity)
	result.Status = "retry"
	result.Data = "somejson"
	jobEntity.Results = append(jobEntity.Results, result)

	g.db.Save(jobEntity)

	return nil
}

func (g GormStorage) RescheduleCron(job Job) error {
	jobEntity := job.GetEntity()

	log.Info().Str("Job", job.DiagnosticsString()).Msg("Scheduling next job")
	res, err := cronParser.Parse(jobEntity.CronExpression)
	if err != nil {
		return err
	}
	next := res.Next(jobEntity.ScheduledTimestamp)
	jobEntity.ScheduledTimestamp = next
	jobEntity.Status = JobStatusScheduled
	result := newJobResultEntity(jobEntity)
	result.Status = "cron"
	result.Data = "somejson"
	jobEntity.Results = append(jobEntity.Results, result)

	g.db.Save(jobEntity)

	return nil
}

func (g GormStorage) Complete(job Job) error {
	jobEntity := job.GetEntity()

	jobEntity.Status = JobStatusCompletedSuccessfully
	result := newJobResultEntity(jobEntity)
	result.Status = "completed"
	result.Data = "somejson"
	jobEntity.Results = append(jobEntity.Results, result)

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

func (g GormStorage) Read(skip int, limit int) ([]JobEntity, error) {
	var jobs []JobEntity

	err := g.db.
		Select("id, data_discriminator, status, scheduled_timestamp, created_timestamp").
		Offset(skip).
		Limit(limit).
		Order("scheduled_timestamp asc").
		Find(&jobs).Error

	if err != nil {
		return nil, err
	}

	return jobs, nil
}
