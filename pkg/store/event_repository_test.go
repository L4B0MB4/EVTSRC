package store_test

import (
	"os"
	"testing"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/L4B0MB4/EVTSRC/pkg/store"
)

func setup() *store.DatabaseConnection {
	db := store.DatabaseConnection{}
	if _, err := os.Stat(store.GetDbFileLocation()); err == nil {
		db.Teardown()
	}
	db.SetUp()
	if !db.IsInitialized() {
		panic("ERROR DURING DB INIT ")
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
		t.Error("CONNECTION SHOULD BE RETRIEVED WITHOUT A PROBLEM")
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
		t.Error("THERE SHOULD BE ONE ENTRY IN THE DB")
		t.Fail()
	}
	if q.Next() == true {
		t.Error("THERE ARE MORE THEN ONE ENTRY IN THE DB")
		t.Fail()
	}

	evs, _ := r.GetEventsForAggregate(ev.AggregateId)
	evToComp := evs[0]
	if evToComp.AggregateId != ev.AggregateId || evToComp.Data[0] != ev.Data[0] || evToComp.Data[1] != ev.Data[1] || evToComp.Name != ev.Name || evToComp.Version != ev.Version {
		t.Error("SOMETHING WENT WRONG DURING SERIALIZATION OR DESERIALIZATION")
		t.Fail()
	}

}

func TestAddEventDuplicate(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("CONNECTION SHOULD BE RETRIEVED WITHOUT A PROBLEM")
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
		t.Error("NO ERROR WHEN ADDING THE SAME EVENT TWICE")
		t.Fail()
	}
}

func TestAddTwoFollowingEvents(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("CONNECTION SHOULD BE RETRIEVED WITHOUT A PROBLEM")
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
		t.Error("CONNECTION SHOULD BE RETRIEVED WITHOUT A PROBLEM")
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
		t.Error("SHOULD HAVE 2 EVENTS FOR THIS AGGREGATE")
		t.Fail()
	}
	i, _ = r.GetEventsForAggregate(newAggType)
	if len(i) != 1 {
		t.Error("SHOULD HAVE 1 EVENT FOR THIS AGGREGATE")
		t.Fail()
	}
}

func TestAddTwoFollowingEventsInOneArray(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("CONNECTION SHOULD BE RETRIEVED WITHOUT A PROBLEM")
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
		t.Error("THERE SHOULD BE ONE ENTRY IN THE DB")
		t.Fail()
	}
	if q.Next() != true {
		t.Error("THERE SHOULD BE A SECOND ENTRY IN THE DB")
		t.Fail()
	}
	if q.Next() == true {
		t.Error("THERE ARE MORE THEN TWO ENTRIES IN THE DB")
		t.Fail()
	}
}

func TestAddTwoEventsWithSameVersionInOneArray(t *testing.T) {
	db := setup()
	defer teardown(db)
	conn, err := db.GetDbConnection()
	if err != nil {
		t.Error("CONNECTION SHOULD BE RETRIEVED WITHOUT A PROBLEM")
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
		Name:          "othervent",
		Data:          []byte{0, 1},
		AggregateId:   "anyaggregateId",
		AggregateType: "aggregateType",
	}
	err = r.AddEvents([]models.Event{ev, ev1})
	if err == nil {
		t.Error("SHOULD HAVE FAILED DUE TO VERSION CLASH")
		t.Fail()
	}
	q, _ := conn.Query("SELECT * FROM events")
	if q.Next() == true {
		t.Error("THERE SHOULD BE NO ENTRY IN THE DB")
		t.Fail()
	}
}
