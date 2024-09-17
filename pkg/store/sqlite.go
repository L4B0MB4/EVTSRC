package store

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

type DatabaseConnection struct {
	initialized bool
	db          *sql.DB
}

var _DBFILE = "./eventstore.db"

func GetDbFileLocation() string {
	return _DBFILE
}

func (d *DatabaseConnection) Teardown() error {
	if d.db != nil {
		d.db.Close()
	}
	return os.Remove(_DBFILE)
}

func (d *DatabaseConnection) SetUp() {
	db, err := sql.Open("sqlite3", _DBFILE)
	if err != nil {

		log.Info().Err(err).Msg("Opening sqlite connection")
		return
	}
	if createEventTable(db) != nil {
		return
	}
	if createEventTableIndex(db) != nil {
		return
	}
	if createAggregateStateTable(db) != nil {
		return
	}
	if createAggregateTableIndex(db) != nil {
		return
	}
	if createAggregateSnapshotTable(db) != nil {
		return
	}
	d.db = db
	d.initialized = true
}

func (d *DatabaseConnection) IsInitialized() bool {
	return d.initialized
}

func createEventTable(db *sql.DB) error {
	//name = name of the event
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS events (id TEXT PRIMARY KEY, aggregateId TEXT, timestamp_0 INTEGER,timestamp_1 INTEGER,Name TEXT, version_0 INTEGER,version_1 INTEGER,data BLOB,UNIQUE(aggregateId,version_0, version_1) ON CONFLICT FAIL)")
	if err != nil {

		log.Info().Err(err).Msg("Preparing statement for events table")
		return err
	}
	_, err = stmt.Exec()
	if err != nil {

		log.Info().Err(err).Msg("Creating events table")
		return err
	}
	return nil
}

func createEventTableIndex(db *sql.DB) error {

	stmt, err := db.Prepare("CREATE INDEX IF NOT EXISTS IX_event__aggregateId ON events(aggregateId)")
	if err != nil {

		log.Info().Err(err).Msg("Preparing statement for events table")
		return err
	}
	_, err = stmt.Exec()
	if err != nil {

		log.Info().Err(err).Msg("Creating index on events table")
		return err
	}
	return nil
}

func createAggregateStateTable(db *sql.DB) error {
	//type = name of the aggregate
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS aggregate_state (id TEXT,version_0 INTEGER,version_1 INTEGER,UNIQUE(id,version_0, version_1) ON CONFLICT FAIL )")
	if err != nil {

		log.Info().Err(err).Msg("Preparing statement for aggregate_state table")
		return err
	}
	_, err = stmt.Exec()
	if err != nil {

		log.Info().Err(err).Msg("Creating aggregate_state table")
		return err
	}
	return nil
}

func createAggregateTableIndex(db *sql.DB) error {

	stmt, err := db.Prepare("CREATE INDEX IF NOT EXISTS IX_aggregate_state__id ON aggregate_state(id);")
	if err != nil {

		log.Info().Err(err).Msg("Preparing statement for aggregate_state table")
		return err
	}
	_, err = stmt.Exec()
	if err != nil {

		log.Info().Err(err).Msg("Creating index on aggregate_state table")
		return err
	}
	return nil
}

func createAggregateSnapshotTable(db *sql.DB) error {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS aggregate_snapshots (id TEXT PRIMARY KEY, name TEXT,version_0 INTEGER,version_1 INTEGER,UNIQUE(version_0, version_1) ON CONFLICT FAIL )")
	if err != nil {

		log.Info().Err(err).Msg("Preparing statement for aggregate_snapshots table")
		return err
	}
	_, err = stmt.Exec()
	if err != nil {

		log.Info().Err(err).Msg("Creating aggregate_snapshots table")
		return err
	}
	return nil
}

func (d *DatabaseConnection) GetDbConnection() (*sql.DB, error) {
	if !d.initialized {
		return nil, errors.New("DatabaseConnection not properly initialized")
	}
	return d.db, nil
}
