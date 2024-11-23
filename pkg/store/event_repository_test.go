package store_test

import (
	"os"
	"testing"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/L4B0MB4/EVTSRC/pkg/store"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setup() *store.DatabaseConnection {
	db := store.DatabaseConnection{}
	if _, err := os.Stat(store.GetDbFileLocation()); err == nil {
		db.Teardown()
	}
	db.SetUp()
	if !db.IsInitialized() {
		panic("Error during DB init")
	}
	return &db
}

func teardown(db *store.DatabaseConnection) {
	err := db.Teardown()
	if err != nil {
		panic(err)
	}
}

func TestAddEventSuccessful(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("Connection should be retrieved without a problem")
		t.Fail()
	}
	r := store.NewEventRepository(conn)
	ev := models.Event{
		Version:       1,
		Name:          "testevent",
		Data:          []byte{0, 1},
		AggregateId:   "anyaggregateId",
		AggregateType: "aggregateType",
	}
	err = r.AddEvents([]models.Event{ev})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	q, _ := conn.Query("SELECT * FROM events")
	if q.Next() != true {
		t.Error("There should be one entry in the DB")
		t.Fail()
	}
	if q.Next() == true {
		t.Error("There are more than one entry in the DB")
		t.Fail()
	}

	evs, _ := r.GetEventsForAggregate(ev.AggregateId)
	evToComp := evs[0]
	if evToComp.AggregateId != ev.AggregateId || evToComp.Data[0] != ev.Data[0] || evToComp.Data[1] != ev.Data[1] || evToComp.Name != ev.Name || evToComp.Version != ev.Version {
		t.Error("Something went wrong during serialization or deserialization")
		t.Fail()
	}
}

func TestAddEventDuplicate(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("Connection should be retrieved without a problem")
		t.Fail()
	}
	r := store.NewEventRepository(conn)
	ev := models.Event{
		Version:       1,
		Name:          "testevent",
		Data:          []byte{0, 1},
		AggregateId:   "anyaggregateId",
		AggregateType: "aggregateType",
	}
	err = r.AddEvents([]models.Event{ev})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	err = r.AddEvents([]models.Event{ev})
	if err == nil {
		t.Error("No error when adding the same event twice")
		t.Fail()
	}
}

func TestAddTwoFollowingEvents(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("Connection should be retrieved without a problem")
		t.Fail()
	}
	r := store.NewEventRepository(conn)
	ev := models.Event{
		Version:       1,
		Name:          "testevent",
		Data:          []byte{0, 1},
		AggregateId:   "anyaggregateId",
		AggregateType: "aggregateType",
	}
	err = r.AddEvents([]models.Event{ev})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	ev.Version++
	err = r.AddEvents([]models.Event{ev})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestAddThreeEventsOfTwoAggregates(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("Connection should be retrieved without a problem")
		t.Fail()
	}
	r := store.NewEventRepository(conn)
	oldAggType := "anyaggregateId"
	ev := models.Event{
		Version:     1,
		Name:        "testevent",
		Data:        []byte{0, 1},
		AggregateId: oldAggType,
	}
	r.AddEvents([]models.Event{ev})
	ev.Version++
	r.AddEvents([]models.Event{ev})
	ev.Version = 0
	newAggType := "aggregateId2"
	ev.AggregateId = newAggType
	r.AddEvents([]models.Event{ev})

	i, _ := r.GetEventsForAggregate(oldAggType)
	if len(i) != 2 {
		t.Error("Should have 2 events for this aggregate")
		t.Fail()
	}
	i, _ = r.GetEventsForAggregate(newAggType)
	if len(i) != 1 {
		t.Error("Should have 1 event for this aggregate")
		t.Fail()
	}
}

func TestAddTwoFollowingEventsInOneArray(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("Connection should be retrieved without a problem")
		t.Fail()
	}
	r := store.NewEventRepository(conn)
	ev := models.Event{
		Version:       1,
		Name:          "testevent",
		Data:          []byte{0, 1},
		AggregateId:   "anyaggregateId",
		AggregateType: "aggregateType",
	}
	ev1 := models.Event{
		Version:       2,
		Name:          "testevent",
		Data:          []byte{0, 1},
		AggregateId:   "anyaggregateId",
		AggregateType: "aggregateType",
	}
	err = r.AddEvents([]models.Event{ev, ev1})
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	q, _ := conn.Query("SELECT * FROM events")
	if q.Next() != true {
		t.Error("There should be one entry in the DB")
		t.Fail()
	}
	if q.Next() != true {
		t.Error("There should be a second entry in the DB")
		t.Fail()
	}
	if q.Next() == true {
		t.Error("There are more than two entries in the DB")
		t.Fail()
	}
}

func TestAddTwoEventsWithSameVersionInOneArray(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("Connection should be retrieved without a problem")
		t.Fail()
	}
	r := store.NewEventRepository(conn)
	ev := models.Event{
		Version:       1,
		Name:          "testevent",
		Data:          []byte{0, 1},
		AggregateId:   "anyaggregateId",
		AggregateType: "aggregateType",
	}
	ev1 := models.Event{
		Version:       1,
		Name:          "otherevent",
		Data:          []byte{0, 1},
		AggregateId:   "anyaggregateId",
		AggregateType: "aggregateType",
	}
	err = r.AddEvents([]models.Event{ev, ev1})
	if err == nil {
		t.Error("Should have failed due to version clash")
		t.Fail()
	}
	q, _ := conn.Query("SELECT * FROM events")
	if q.Next() == true {
		t.Error("There should be no entry in the DB")
		t.Fail()
	}
}

func TestGetEventsSinceEvent(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("Connection should be retrieved without a problem")
		t.Fail()
	}
	repo := store.NewEventRepository(conn)

	// Insert test data
	event1 := models.Event{
		AggregateId:   "agg1",
		Name:          "Event1",
		Version:       1,
		Data:          []byte("data1"),
		AggregateType: "type1",
	}
	event2 := models.Event{
		AggregateId:   "agg1",
		Name:          "Event2",
		Version:       2,
		Data:          []byte("data2"),
		AggregateType: "type1",
	}
	event3 := models.Event{
		AggregateId:   "agg1",
		Name:          "Event3",
		Version:       3,
		Data:          []byte("data3"),
		AggregateType: "type1",
	}
	event4 := models.Event{
		AggregateId:   "agg2",
		Name:          "Event4",
		Version:       1,
		Data:          []byte("data4"),
		AggregateType: "type2",
	}
	event5 := models.Event{
		AggregateId:   "agg2",
		Name:          "Event5",
		Version:       2,
		Data:          []byte("data5"),
		AggregateType: "type2",
	}

	err = repo.AddEvents([]models.Event{event1, event2, event3, event4, event5})
	assert.NoError(t, err)

	events, err := repo.GetEventsSinceEvent(event1.Id, 2)
	assert.NoError(t, err)
	assert.Len(t, events, 4)
	assert.Equal(t, event1.Name, events[0].Name)
	assert.Equal(t, event2.Name, events[1].Name)
	events, err = repo.GetEventsSinceEvent(events[1].Id, 2)
	assert.NoError(t, err)
	assert.Len(t, events, 4)
	assert.Equal(t, event3.Name, events[0].Name)
	assert.Equal(t, event4.Name, events[1].Name)
}
