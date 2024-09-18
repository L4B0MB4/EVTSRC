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
)

func setup() (*client.EventSourcingHttpClient, *httphandler.HttpHandler, *store.DatabaseConnection) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	db := store.DatabaseConnection{}
	db.SetUp()
	conn, err := db.GetDbConnection()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessfull initalization of db")
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

	err := client.AddEventsWithoutValidation("myaggregate4444", []models.Event{{Version: 1, Name: "asdasd", Data: []byte{0, 1, 2}, AggregateType: "mytype"}, {Version: 2, Name: "asdasd2", Data: []byte{1, 2, 3}, AggregateType: "mytype"}})
	if err != nil {
		log.Error().Err(err).Msg("ERROR ADDING EVENTS")
		t.Fail()
	}
	err = client.AddEventsWithoutValidation("differentaggregate", []models.Event{{Version: 7, Name: "asdasd", Data: []byte{0, 1, 2}, AggregateType: "mytype"}, {Version: 8, Name: "asdasd2", Data: []byte{1, 2, 3}, AggregateType: "mytype"}})
	if err != nil {
		log.Error().Err(err).Msg("ERROR ADDING EVENTS FOR SECOND AGGREGATE")
		t.Fail()
	}
	evs, err := client.GetEventsOrdered("myaggregate4444")
	if err != nil {
		log.Error().Err(err).Msg("ERROR RETRIEVING EVENTS")
		t.Fail()
	}
	ev, ok := evs.Next()
	if !ok {
		log.Error().Err(err).Msg("DOES NOT HAVE ONE EVENT")
	}
	if ev.Version != 1 {
		log.Error().Err(err).Msg("DOES HAVE WRONG FIRST EVENT")
	}
	ev, ok = evs.Next()
	if !ok {
		log.Error().Err(err).Msg("DOES NOT HAVE TWO EVENTS")
	}
	if ev.Version != 2 {
		log.Error().Err(err).Msg("DOES HAVE WRONG SECOND EVENT")
	}

}
