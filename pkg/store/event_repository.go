package store

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/helper"
	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/L4B0MB4/EVTSRC/pkg/models/customerrors"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type EventRepository struct {
	store *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	if db == nil {
		return nil
	}
	return &EventRepository{store: db}
}

func (e *EventRepository) AddEvents(events []models.Event) error {

	tx, err := e.store.Begin()
	if err != nil {
		return err
	}

	for _, event := range events {

		eEvent := &eventEntity{
			Event:     event,
			timestamp: time.Now(),
			id:        uuid.New(),
		}
		err = e.addEvent(tx, eEvent)
		if err != nil {
			tx.Rollback()
			log.Info().Err(err).Msg("ABORTED TX")
			return err
		}
	}

	return tx.Commit()
}

func (e *EventRepository) addEvent(tx *sql.Tx, event *eventEntity) error {
	t0, t1, err := helper.SplitInt62(event.timestamp.UnixMicro())
	if err != nil {
		return err
	}

	v0, v1, err := helper.SplitInt62(event.Version)
	if err != nil {
		return err
	}

	stmt, err := e.store.Prepare(`
        INSERT INTO events (id, aggregateId,timestamp_0 ,timestamp_1, Name, version_0, version_1, data)
        VALUES (?,?,?,?,?,?,?,?)
    `)
	if err != nil {
		log.Info().Err(err).Msg("Preparing insert statement for events table")
		return err
	}
	defer stmt.Close()

	_, err = tx.Stmt(stmt).Exec(event.id, event.AggregateId, t0, t1, event.Name, v0, v1, event.Data)
	if err != nil {
		tx.Rollback()
		log.Info().Err(err).Msg("ABORTED TX")
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return &customerrors.DuplicateVersionError{}
		}
		return err
	}

	stmtAgg, err := e.store.Prepare(`
        INSERT INTO aggregate_state(id,version_0, version_1)
        VALUES (?,?,?)
    `)
	if err != nil {
		log.Info().Err(err).Msg("Preparing insert statement for events table")
		return err
	}
	defer stmt.Close()

	_, err = tx.Stmt(stmtAgg).Exec(event.AggregateId, v0, v1)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventRepository) GetEventsForAggregate(aggregateType string) ([]models.Event, error) {

	// Prepare the SQL query
	query := `
		SELECT events.Name, events.version_0, events.version_1, events.data,events.aggregateId 
		FROM events 
		JOIN aggregate_state 
			ON events.aggregateId = aggregate_state.id 
		 	and events.version_0 = aggregate_state.version_0 
			and events.version_1 = aggregate_state.version_1 
		WHERE aggregate_state.id = ?
		ORDER BY events.version_0 ASC, events.version_1 ASC
	`

	stmt, err := e.store.Prepare(query)
	if err != nil {
		log.Info().Err(err).Msg("Error preparing statement")
		return nil, errors.New("COULD NOT PREPARE STATMENT FOR QUERY EVENTS")
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(aggregateType)
	if err != nil {
		log.Info().Err(err).Msg("Error running query statement")
		return nil, errors.New("COULD NOT QUERY EVENTS")
	}
	defer rows.Close()

	// Initialize a slice to hold all events
	var events []models.Event

	// Iterate over the results and store them in the slice
	for rows.Next() {
		var event models.Event

		var v0 int32
		var v1 int32
		err = rows.Scan(&event.Name, &v0, &v1, &event.Data, &event.AggregateId)
		if err != nil {
			log.Info().Err(err).Msg("Error scanning rows")
			return nil, errors.New("COULD NOT RETRIEVE EVENT")
		}
		version, err := helper.MergeInt62(v0, v1)
		if err != nil {
			log.Info().Err(err).Msg("Error transforming version")
			return nil, errors.New("COULD NOT RETRIEVE EVENT")
		}
		event.Version = version

		// Append the event to the slice
		events = append(events, event)
	}

	// Check for any error that might have occurred during iteration
	if err = rows.Err(); err != nil {
		log.Info().Err(err).Msg("Error checking row errors")
		return nil, errors.New("COULD NOT RETRIEVE ALL EVENTS")
	}
	return events, nil
}
