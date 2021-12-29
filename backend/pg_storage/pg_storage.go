package pg_storage

import (
	"database/sql"
	"encoding/json"
	"github.com/asynkron/CallMeLater/server"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"time"
)

type PgStorage struct {
	db *sql.DB
}

type PgJob struct {
	Id                 string
	ScheduledTimestamp time.Time
	CreatedTimestamp   time.Time
	CompletedTimestamp sql.NullTime
	Data               string
	RetryCount         int
	ParentJobId        string
}

func New(connectionString string) *PgStorage {
	log.
		Info().
		Str("connectionString", connectionString).
		Msg("Connecting to PostgreSQL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.
			Err(err).
			Msg("Failed to connect to Postgres")

		panic(err)
	}
	log.
		Info().
		Str("connectionString", connectionString).
		Msg("Connected to PostgreSQL")

	return &PgStorage{db: db}
}

func (p *PgStorage) Pull(count int) ([]*server.RequestData, error) {
	rows, err := p.db.Query(
		`
		SELECT * 
		FROM "Jobs" 
		WHERE "CompletedTimestamp" is null  
		ORDER BY "ScheduledTimestamp" DESC 
		LIMIT $1`,
		count,
	)
	if err != nil {
		log.
			Err(err).
			Msg("Failed to get requests")

		return nil, err
	}

	var r []*server.RequestData
	//loop over rows and add to slice
	for rows.Next() {
		pgRow := &PgJob{}
		err := rows.Scan(
			&pgRow.Id,
			&pgRow.ScheduledTimestamp,
			&pgRow.CreatedTimestamp,
			&pgRow.CompletedTimestamp,
			&pgRow.Data,
			&pgRow.RetryCount,
			&pgRow.ParentJobId,
		)
		if err != nil {
			log.
				Err(err).
				Msg("Failed to scan row")

			return nil, err
		}

		rr := &server.RequestData{}
		var d = []byte(pgRow.Data)
		err = json.Unmarshal(d, rr)
		if err != nil {
			log.
				Err(err).
				Msg("Failed to unmarshal row")

			return nil, err
		}

		r = append(r, rr)
	}

	return r, nil
}

func (p *PgStorage) Push(data *server.RequestData) error {

	j, err := json.Marshal(data)
	if err != nil {
		log.
			Err(err).
			Msg("Failed to marshal data")

		return err
	}

	pgRow := &PgJob{
		Id:                 data.Id,
		ScheduledTimestamp: data.ScheduledTimestamp,
		CreatedTimestamp:   time.Now(),
		Data:               string(j),
	}

	_, err = p.db.Exec(
		`
		INSERT INTO "Jobs" 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		pgRow.Id,
		pgRow.ScheduledTimestamp,
		pgRow.CreatedTimestamp,
		nil,
		j,
		0,
		pgRow.ParentJobId,
	)

	if err != nil {
		log.
			Err(err).
			Msg("Failed to insert row")

		return err
	}
	log.Info().
		Str("Id", data.Id).
		Str("Url", data.RequestUrl).
		Msg("Inserted new request")

	return nil
}

func (p *PgStorage) Complete(requestId string) error {
	var _, err = p.db.Exec(
		`
		UPDATE "Jobs" 
		SET "CompletedTimestamp" = $1 
		WHERE "Id" = $2 `,
		time.Now(),
		requestId,
	)

	if err != nil {
		log.
			Err(err).
			Msg("Failed to update row")

		return err
	}

	log.Info().
		Str("requestId", requestId).
		Msg("Deleted request")
	return err
}
