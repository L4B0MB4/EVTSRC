package client

import "github.com/L4B0MB4/EVTSRC/pkg/models"

// EventsIterator iterates over a list of events.
type EventsIterator struct {
	events []models.Event
	index  int
}

// NewEventIterator creates a new EventsIterator.
func NewEventIterator(events []models.Event) *EventsIterator {
	return &EventsIterator{
		events: events,
		index:  -1,
	}
}

// Next returns the next event in the iterator.
func (e *EventsIterator) Next() (*models.Event, bool) {
	e.index++
	if e.index >= len(e.events) || e.index < 0 {
		return nil, false
	}
	ev := &e.events[e.index]
	return ev, true
}

// Current returns the current event in the iterator.
func (e *EventsIterator) Current() *models.Event {
	if e.index >= len(e.events) || e.index < 0 {
		return nil
	}
	return &e.events[e.index]
}

// Reset resets the iterator to the beginning.
func (e *EventsIterator) Reset() {
	e.index = -1
}
