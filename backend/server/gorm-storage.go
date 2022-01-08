package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"time"
)

type GormStorage struct {
	db *gorm.DB
}

var (
	zeroTime = time.Time{}
)

func (g GormStorage) Cancel(job Job) error {
	jobEntity := job.GetEntity()
	jobEntity.ScheduleTimestamp = zeroTime
	jobEntity.ScheduleStatus = JobStatusCancelled
	g.db.Save(jobEntity)

	return nil
}

func (g GormStorage) Fail(job Job) error {
	jobEntity := job.GetEntity()
	jobEntity.ExecutedTimestamp = time.Now()
	jobEntity.ScheduleTimestamp = zeroTime
	jobEntity.ExecutedStatus = ExecutedStatusFail
	jobEntity.ScheduleStatus = JobStatusFailed

	result := newJobResultEntity(jobEntity)
	result.Status = "failed"
	result.Data = "somejson"
	jobEntity.Results = append(jobEntity.Results, result)

	if jobEntity.Data == "" {
		panic("Job data is empty - failed")
	}

	g.db.Save(jobEntity)

	return nil
}

func newJobResultEntity(jobEntity *JobEntity) JobResultEntity {
	result := JobResultEntity{
		Id:                 uuid.New().String(),
		JobId:              jobEntity.Id,
		ExecutionTimestamp: jobEntity.ScheduleTimestamp,
		DataDiscriminator:  jobEntity.DataDiscriminator,
	}
	return result
}

func (g GormStorage) Pull(count int) ([]Job, error) {
	var jobs []JobEntity
	err := g.db.
		Limit(count).
		Order("schedule_timestamp asc").
		Find(&jobs, "schedule_status = ?", JobStatusScheduled).Error

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

	if jobEntity.Data == "" {
		panic("Job data is empty - create")
	}

	g.db.Create(jobEntity)

	return nil
}

func (g GormStorage) Retry(job Job) error {
	jobEntity := job.GetEntity()
	jobEntity.ExecutedTimestamp = time.Now()
	jobEntity.ExecutedStatus = ExecutedStatusFail
	jobEntity.ScheduleStatus = JobStatusScheduled
	jobEntity.ExecutedCount++
	result := newJobResultEntity(jobEntity)
	result.Status = "retry"
	result.Data = "somejson"
	jobEntity.Results = append(jobEntity.Results, result)

	if jobEntity.Data == "" {
		panic("Job data is empty - retry")
	}

	g.db.Save(jobEntity)

	return nil
}

func (g GormStorage) RescheduleCron(job Job) error {
	jobEntity := job.GetEntity()

	jobEntity.ExecutedTimestamp = time.Now()
	jobEntity.ExecutedStatus = ExecutedStatusSuccess

	log.Info().Str("Job", job.DiagnosticsString()).Msg("Scheduling next job")
	res, err := cronParser.Parse(jobEntity.ScheduleCronExpression)
	if err != nil {
		return err
	}

	next := res.Next(jobEntity.ScheduleTimestamp)
	jobEntity.ScheduleTimestamp = next
	jobEntity.ScheduleStatus = JobStatusScheduled
	result := newJobResultEntity(jobEntity)
	result.Status = "cron"
	result.Data = "somejson"
	jobEntity.Results = append(jobEntity.Results, result)

	if jobEntity.Data == "" {
		panic("Job data is empty - reschedule cron")
	}

	g.db.Save(jobEntity)

	return nil
}

func (g GormStorage) Complete(job Job) error {
	jobEntity := job.GetEntity()
	jobEntity.ExecutedTimestamp = time.Now()
	jobEntity.ScheduleTimestamp = zeroTime
	jobEntity.ExecutedStatus = ExecutedStatusSuccess

	jobEntity.ScheduleStatus = JobStatusSuccess
	result := newJobResultEntity(jobEntity)
	result.Status = "completed"
	result.Data = "somejson"
	jobEntity.Results = append(jobEntity.Results, result)

	if jobEntity.Data == "" {
		panic("Job data is empty - complete")
	}

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
		//	Select("id, data_discriminator, status, scheduled_timestamp, created_timestamp, executed_timestamp, description").
		Offset(skip).
		Limit(limit).
		Order("schedule_timestamp asc").
		Find(&jobs).Error

	if err != nil {
		return nil, err
	}

	return jobs, nil
}
