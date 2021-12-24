package main

import (
	"database/sql"
	"github.com/rs/zerolog/log"
)

type PgStorage struct {
	db *sql.DB
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
	//TODO implement me
	panic("implement me")
}

func (p *PgStorage) Set(id string, data *requestData) error {
	_, err := p.db.Exec("INSERT INTO requests (id, data) VALUES ($1, $2)", id, data)

	log.Info().
		Str("id", id).
		Msg("Inserted new request")
	return err
}
