package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"

	"github.com/L4B0MB4/EVTSRC/pkg/models"
	"github.com/rs/zerolog/log"
)

type EventSourcingHttpClient struct {
	httpClient *http.Client
	url        string
}

func stripOldEvents(events []models.ChangeTrackedEvent) []models.Event {
	newEvents := []models.Event{}
	for _, e := range events {

		if e.IsNew {
			ev := models.Event{
				Version:       e.Version,
				Name:          e.Name,
				Data:          e.Data,
				AggregateId:   e.AggregateId,
				AggregateType: e.AggregateType,
			}
			newEvents = append(newEvents, ev)

		}
	}
	return newEvents
}

func NewEventSourcingHttpClient(urlStr string) (*EventSourcingHttpClient, error) {

	path, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	baseUrl := fmt.Sprintf("%s://%s", path.Scheme, path.Host)

	httpClient := http.Client{}
	return &EventSourcingHttpClient{
		httpClient: &httpClient,
		url:        baseUrl,
	}, nil
}

func (client *EventSourcingHttpClient) AddEvents(aggregateId string, events []models.ChangeTrackedEvent) error {
	if len(aggregateId) <= 0 {
		return fmt.Errorf("aggregateId empty")
	}
	for _, event := range events {
		if len(event.AggregateType) <= 0 {
			return fmt.Errorf("aggregateType empty")
		}
		if len(event.Data) == 0 {
			return fmt.Errorf("data empty")
		}
		if len(event.Name) == 0 {
			return fmt.Errorf("name empty")
		}
	}
	return client.AddEventsWithoutValidation(aggregateId, events)
}

func (client *EventSourcingHttpClient) AddEventsWithoutValidation(aggregateId string, events []models.ChangeTrackedEvent) error {

	newEvents := stripOldEvents(events)
	bodyBytes, err := json.Marshal(newEvents)
	if err != nil {
		log.Info().Err(err).Msg("could not marshal events")
		return err
	}
	buf := bytes.NewBuffer(bodyBytes)
	addEventsUrl, err := url.JoinPath(client.url, fmt.Sprintf("/%s/events", aggregateId))
	if err != nil {
		log.Info().Err(err).Msg("could not use url")
		return err
	}

	resp, err := client.httpClient.Post(addEventsUrl, "application/json", buf)

	if err != nil {
		log.Info().Err(err).Msg("error during the request")
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Info().Err(err).Msg("got non 2XX header")
		return fmt.Errorf("unsuccessful request")
	}
	return nil
}

func (client *EventSourcingHttpClient) GetEventsOrdered(aggregateId string) (*EventsIterator, error) {

	getEventsUrl, err := url.JoinPath(client.url, fmt.Sprintf("/%s/events", aggregateId))
	if err != nil {
		log.Info().Err(err).Msg("could not use url")
		return nil, err
	}

	resp, err := client.httpClient.Get(getEventsUrl)
	if err != nil {
		log.Info().Err(err).Msg("error during the request")
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Info().Err(err).Msg("got non 2XX header")
		return nil, fmt.Errorf("unsuccessful request")
	}
	var events []models.Event
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Info().Err(err).Msg("error during reading response body")
		return nil, err
	}
	err = json.Unmarshal(buf, &events)

	if err != nil {
		log.Info().Err(err).Msg("error during unmarshalling body")
		return nil, err
	}

	slices.SortFunc(events, func(i models.Event, j models.Event) int {
		//ascending
		delta := i.Version - j.Version
		if delta > 0 {
			return 1
		} else if delta < 0 {
			return -1
		}
		return 0
	})
	eventsIterator := NewEventIterator(events)
	return eventsIterator, nil
}

func (client *EventSourcingHttpClient) GetEventsSince(eventId string, limit int) ([]models.Event, error) {
	if len(eventId) == 0 {
		eventId = "0"
	}
	if limit <= 0 {
		return nil, fmt.Errorf("invalid limit value")
	}
	if limit > 100 {
		limit = 100
	}

	getEventsSinceUrl, err := url.JoinPath(client.url, "/events", url.PathEscape(eventId), "since")
	if err != nil {
		log.Info().Err(err).Msg("could not use url")
		return nil, err
	}

	query := url.Values{}
	query.Set("limit", fmt.Sprintf("%d", limit))
	getEventsSinceUrl = fmt.Sprintf("%s?%s", getEventsSinceUrl, query.Encode())

	resp, err := client.httpClient.Get(getEventsSinceUrl)
	if err != nil {
		log.Info().Err(err).Msg("error during the request")
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Info().Err(err).Msg("got non 2XX header")
		return nil, fmt.Errorf("unsuccessful request")
	}

	var events []models.Event
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Info().Err(err).Msg("error during reading response body")
		return nil, err
	}
	err = json.Unmarshal(buf, &events)
	if err != nil {
		log.Info().Err(err).Msg("error during unmarshalling body")
		return nil, err
	}

	return events, nil
}
