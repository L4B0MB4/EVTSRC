package integrationtest

import (
	"os"
	"testing"

	"github.com/L4B0MB4/EVTSRC/pkg/client"
	"github.com/L4B0MB4/EVTSRC/pkg/httphandler"
	"github.com/L4B0MB4/EVTSRC/pkg/httphandler/controller"
	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/L4B0MB4/EVTSRC/pkg/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func setup() (*client.EventSourcingHttpClient, *httphandler.HttpHandler, *store.DatabaseConnection) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	db := store.DatabaseConnection{}
	db.SetUp()
	conn, err := db.GetDbConnection()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessful initialization of db")
		panic(err)
	}
	log.Debug().Msg("Db Connection was successful")
	repository := store.NewEventRepository(conn)

	c := controller.NewEventController(repository)
	h := httphandler.NewHttpHandler(c)

	go func() {
		h.Start()
	}()
	evclient, err := client.NewEventSourcingHttpClient("http://localhost:5515")
	if err != nil {
		panic(err)
	}
	return evclient, h, &db
}

func teardown(httpHandler *httphandler.HttpHandler, db *store.DatabaseConnection) {
	httpHandler.Stop()
	db.Teardown()
}

func TestClientAddingEventsAndRetrievingThemFromServer(t *testing.T) {
	client, httpHandler, db := setup()
	defer teardown(httpHandler, db)

	err := client.AddEventsWithoutValidation("myaggregate4444", []models.ChangeTrackedEvent{
		{IsNew: true, Event: models.Event{Version: 1, Name: "asdasd", Data: []byte{0, 1, 2}, AggregateType: "mytype"}},
		{IsNew: true, Event: models.Event{Version: 2, Name: "asdasd2", Data: []byte{1, 2, 3}, AggregateType: "mytype"}}})
	if err != nil {
		log.Error().Err(err).Msg("Error adding events")
		t.Fail()
	}
	err = client.AddEventsWithoutValidation("differentaggregate", []models.ChangeTrackedEvent{
		{IsNew: true, Event: models.Event{Version: 7, Name: "asdasd", Data: []byte{0, 1, 2}, AggregateType: "mytype"}},
		{IsNew: true, Event: models.Event{Version: 8, Name: "asdasd2", Data: []byte{1, 2, 3}, AggregateType: "mytype"}}})
	if err != nil {
		log.Error().Err(err).Msg("Error adding events for second aggregate")
		t.Fail()
	}
	evs, err := client.GetEventsOrdered("myaggregate4444")
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving events")
		t.Fail()
	}
	ev, ok := evs.Next()
	if !ok {
		log.Error().Err(err).Msg("Does not have one event")
	}
	if ev.Version != 1 {
		log.Error().Err(err).Msg("Does have wrong first event")
	}
	ev, ok = evs.Next()
	if !ok {
		log.Error().Err(err).Msg("Does not have two events")
	}
	if ev.Version != 2 {
		log.Error().Err(err).Msg("Does have wrong second event")
	}

}

func TestClientGetEventsSince(t *testing.T) {
	client, httpHandler, db := setup()
	defer teardown(httpHandler, db)

	err := client.AddEventsWithoutValidation("myaggregate4444", []models.ChangeTrackedEvent{
		{IsNew: true, Event: models.Event{Version: 1, Name: "event1", Data: []byte{0, 1, 2}, AggregateType: "mytype"}},
		{IsNew: true, Event: models.Event{Version: 2, Name: "event2", Data: []byte{1, 2, 3}, AggregateType: "mytype"}},
		{IsNew: true, Event: models.Event{Version: 3, Name: "event3", Data: []byte{2, 3, 4}, AggregateType: "mytype"}},
	})

	if err != nil {
		log.Error().Err(err).Msg("Error adding events")
		t.Fail()
	}

	events, err := client.GetEventsSince("", 2)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving events since eventId")
		t.Fail()
	}

	if len(events) != 2 {
		log.Error().Msgf("Expected 2 events, got %d", len(events))
		t.Fail()
	}

	if events[0].Version != 1 || events[1].Version != 2 {
		log.Error().Msg("Events are not in the correct order or incorrect events retrieved")
		t.Fail()
	}
	events, err = client.GetEventsSince(events[1].Id, 2)
	assert.NoError(t, err)
	assert.Len(t, events, 1)

	if events[0].Version != 3 {
		log.Error().Msg("Events are not in the correct order or incorrect events retrieved")
		t.Fail()
	}
}
