package store

import (
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/google/uuid"
)

// eventEntity represents an event with additional metadata.
type eventEntity struct {
	models.Event
	timestamp time.Time
	id        uuid.UUID
}
