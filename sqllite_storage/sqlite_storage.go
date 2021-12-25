package sqllite_storage

import (
	"database/sql"
	"github.com/asynkron/CallMeLater/pg_storage"
	"github.com/asynkron/CallMeLater/server"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"time"
)

type SqLiteStorage struct {
	db *sql.DB
}

type SqLiteRow struct {
	RequestId string
	Timestamp time.Time
	Data      server.RequestData
}

func New(connectionString string) *SqLiteStorage {
	log.
		Info().
		Str("connectionString", connectionString).
		Msg("Connecting to SQLite")

	db, err := sql.Open("sqlite3", "./storage.db")
	if err != nil {
		log.
			Err(err).
			Msg("Failed to connect to sqlite")

		panic(err)
	}
	log.
		Info().
		Str("connectionString", connectionString).
		Msg("Connected to SQLite")
	SqLiteStorage.CreateTable(SqLiteStorage{db: db})
	return &SqLiteStorage{db: db}
}

func (sl *SqLiteStorage) Pull(count int) ([]*server.RequestData, error) {
	rows, err := sl.db.Query(`SELECT * FROM "Requests" ORDER BY "Timestamp" DESC LIMIT $1`, count)
	if err != nil {
		log.
			Err(err).
			Msg("Failed to get requests")

		return nil, err
	}

	var r []*server.RequestData
	//loop over rows and add to slice
	for rows.Next() {
		pgRow := &pg_storage.PgRow{}
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

func (sl *SqLiteStorage) Push(data *server.RequestData) error {
	var _, err = sl.db.Exec(
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

func (sl *SqLiteStorage) Complete(requestId string) error {
	var _, err = sl.db.Exec(
		`DELETE FROM "Requests" WHERE "RequestId" = $1`,
		requestId,
	)

	log.Info().
		Str("requestId", requestId).
		Msg("Deleted request")
	return err
}

func (sl SqLiteStorage) CreateTable() error {
	var _, err = sl.db.Exec(
		`create table Requests (RequestId string not null primary key, Timestamp text, Data text);`)

	log.Info().
		Msg("Created Requests table")
	return err
}
