package pg_storage

import (
	"database/sql"
	"github.com/asynkron/CallMeLater/server"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"time"
)

type PgStorage struct {
	db *sql.DB
}

type PgRow struct {
	RequestId string
	Timestamp time.Time
	Data      server.RequestData
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

func (p *PgStorage) Get() ([]*server.RequestData, error) {
	//gets the top 1000 requests
	rows, err := p.db.Query(`SELECT * FROM "Requests" ORDER BY "Timestamp" DESC LIMIT 100`)
	if err != nil {
		log.
			Err(err).
			Msg("Failed to get requests")

		return nil, err
	}

	var r []*server.RequestData
	//loop over rows and add to slice
	for rows.Next() {
		pgRow := &PgRow{}
		err := rows.Scan(&pgRow.RequestId, &pgRow.Timestamp, &pgRow.Data)
		if err != nil {
			log.
				Err(err).
				Msg("Failed to scan row")

			return nil, err
		}
		r = append(r, &pgRow.Data)
	}

	return r, nil
}

func (p *PgStorage) Set(data *server.RequestData) error {
	var _, err = p.db.Exec(
		`INSERT INTO "Requests" VALUES ($1, $2, $3)`,
		data.RequestId,
		data.When,
		data,
	)

	log.Info().
		Str("id", data.RequestId).
		Str("Url", data.RequestUrl).
		Msg("Inserted new request")
	return err
}

func (p *PgStorage) Delete(requestId string) error {
	var _, err = p.db.Exec(
		`DELETE FROM "Requests" WHERE "RequestId" = $1`,
		requestId,
	)

	log.Info().
		Str("requestId", requestId).
		Msg("Deleted request")
	return err
}
