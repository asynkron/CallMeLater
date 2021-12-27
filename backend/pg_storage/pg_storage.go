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
	RequestId          string
	ScheduledTimestamp time.Time
	CreatedTimestamp   time.Time
	CompletedTimestamp sql.NullTime
	Data               string
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
	rows, err := p.db.Query(`SELECT * FROM "Requests" ORDER BY "ScheduledTimestamp" DESC LIMIT $1`, count)
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
			&pgRow.RequestId,
			&pgRow.ScheduledTimestamp,
			&pgRow.CreatedTimestamp,
			&pgRow.CompletedTimestamp,
			&pgRow.Data,
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
		RequestId:          data.RequestId,
		ScheduledTimestamp: data.ScheduledTimestamp,
		CreatedTimestamp:   time.Now(),
		Data:               string(j),
	}

	_, err = p.db.Exec(
		`INSERT INTO "Requests" VALUES ($1, $2, $3, $4, $5)`,
		pgRow.RequestId,
		pgRow.ScheduledTimestamp,
		pgRow.CreatedTimestamp,
		nil,
		j,
	)

	log.Info().
		Str("id", data.RequestId).
		Str("Url", data.RequestUrl).
		Msg("Inserted new request")
	return err
}

func (p *PgStorage) Complete(requestId string) error {
	var _, err = p.db.Exec(
		`DELETE FROM "Requests" WHERE "RequestId" = $1`,
		requestId,
	)

	log.Info().
		Str("requestId", requestId).
		Msg("Deleted request")
	return err
}
