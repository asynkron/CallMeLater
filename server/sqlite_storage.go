package callmelater

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
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
	Data      requestData
}

func NewSqLiteStorage(connectionString string) *SqLiteStorage {
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

func (sl *SqLiteStorage) Get() ([]*requestData, error) {
	//gets the top 1000 requests
	rows, err := sl.db.Query(`SELECT * FROM "Requests" ORDER BY "Timestamp" DESC LIMIT 100`)
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

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a *requestData) Values() (driver.Value, error) {
	return json.Marshal(a)
}

func (sl *requestData) Scans(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return nil
	}

	err := json.Unmarshal(source, sl)
	if err != nil {
		return err
	}

	return nil
}

func (sl *SqLiteStorage) Set(data *requestData) error {
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

func (sl *SqLiteStorage) Delete(requestId string) error {
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
