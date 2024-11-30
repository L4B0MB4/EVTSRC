package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/L4B0MB4/EVTSRC/pkg/models/customerrors"
	"github.com/L4B0MB4/EVTSRC/pkg/store"
	"github.com/L4B0MB4/EVTSRC/pkg/tcp/server"
	"github.com/gin-gonic/gin"
)

// EventController handles HTTP requests for events.
type EventController struct {
	repo      *store.EventRepository
	tcpServer *server.TcpEventServer
}

// NewEventController creates a new EventController.
func NewEventController(repo *store.EventRepository, tcpServer *server.TcpEventServer) *EventController {
	return &EventController{
		repo:      repo,
		tcpServer: tcpServer,
	}
}

// GetEventsForAggregate handles the retrieval of events for a given aggregate ID.
func (ctrl *EventController) GetEventsForAggregate(c *gin.Context) {

	aggregateId := c.Param("aggregateId")

	if len(strings.TrimSpace(aggregateId)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path param cant be empty or null"})
		return
	}

	resp, err := ctrl.repo.GetEventsForAggregate(aggregateId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unkown error occured"})
		return
	}
	if len(resp) == 0 {
		c.JSON(http.StatusOK, []models.Event{})
		return
	}
	c.JSON(http.StatusOK, &resp)
}

// AddEventToAggregate handles the addition of events to a given aggregate ID.
func (ctrl *EventController) AddEventToAggregate(c *gin.Context) {
	var events []models.Event
	aggregateId := c.Param("aggregateId")
	if len(strings.TrimSpace(aggregateId)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path param cant be empty or null"})
		return
	}
	if err := c.ShouldBindJSON(&events); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := range events {
		events[i].AggregateId = aggregateId
	}
	err := ctrl.repo.AddEvents(events)
	if err != nil {
		_, ok := err.(*customerrors.DuplicateVersionError)
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error trying to add the same event multiple times"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unkown error occured"})
		return
	}
	ctrl.tcpServer.SendEvent("NewEvent")
}

// GetEventsSince handles the retrieval of events since a given event ID with a limit.
func (ctrl *EventController) GetEventsSince(c *gin.Context) {
	eventId := c.Param("eventId")
	if len(strings.TrimSpace(eventId)) == 0 {
		eventId = "0"
	}
	limitStr := c.Query("limit")
	limit := 100
	if len(strings.TrimSpace(limitStr)) > 0 {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
			return
		}
		if limit > 100 {
			limit = 100
		}
	}
	resp, err := ctrl.repo.GetEventsSinceEvent(eventId, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unkown error occured"})
		return
	}
	if len(resp) == 0 {
		c.JSON(http.StatusOK, []models.Event{})
		return
	}
	c.JSON(http.StatusOK, &resp)
}
