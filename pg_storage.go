package main

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"time"
)

type PgStorage struct {
	db *sql.DB
}

type PgRow struct {
	id   int
	when time.Time
	data *requestData
}

func NewPgStorage(connectionString string) *PgStorage {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.
			Err(err).
			Msg("Failed to connect to Postgres")

		panic(err)
	}
	return &PgStorage{db: db}
}

func (p *PgStorage) Get() ([]*requestData, error) {
	//gets the top 1000 requests
	rows, err := p.db.Query("SELECT * FROM requests ORDER BY id DESC LIMIT 1000")
	if err != nil {
		log.
			Err(err).
			Msg("Failed to get requests")

		return nil, err
	}

	var r []*requestData
	//loop over rows and add to slice
	for rows.Next() {
		pgRow := &PgRow{}
		err := rows.Scan(&pgRow.id, &pgRow.when, &pgRow.data)
		if err != nil {
			log.
				Err(err).
				Msg("Failed to scan row")

			return nil, err
		}
		r = append(r, pgRow.data)
	}

	return r, nil
}

func (p *PgStorage) Set(data *requestData) error {
	_, err := p.db.Exec("INSERT INTO requests (id, when, data) VALUES ($1, $2, $3)", data.RequestId, data.When, data)

	log.Info().
		Str("id", data.RequestId).
		Str("Url", data.RequestUrl).
		Msg("Inserted new request")
	return err
}
