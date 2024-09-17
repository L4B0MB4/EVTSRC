package client

import "github.com/L4B0MB4/EVTSRC/pkg/models"

type EventsIterator struct {
	events []models.Event
	index  int
}

func NewEventIterator(events []models.Event) *EventsIterator {
	return &EventsIterator{
		events: events,
		index:  -1,
	}
}

func (e *EventsIterator) Next() (*models.Event, bool) {
	e.index++
	if e.index >= len(e.events) || e.index < 0 {
		return nil, false
	}
	ev := &e.events[e.index]
	return ev, true
}

func (e *EventsIterator) Current() *models.Event {
	if e.index >= len(e.events) || e.index < 0 {
		return nil
	}
	return &e.events[e.index]
}

func (e *EventsIterator) Reset() {
	e.index = -1
}
