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

func (client *EventSourcingHttpClient) AddEvents(aggregateId string, events []models.Event) error {
	bodyBytes, err := json.Marshal(events)
	if err != nil {
		log.Info().Err(err).Msg("Could not marshal events")
		return err
	}
	buf := bytes.NewBuffer(bodyBytes)
	addEventsUrl, err := url.JoinPath(client.url, fmt.Sprintf("/%s/events", aggregateId))
	if err != nil {
		log.Info().Err(err).Msg("Could not use url")
		return err
	}

	resp, err := client.httpClient.Post(addEventsUrl, "application/json", buf)

	if err != nil {
		log.Info().Err(err).Msg("Error during the request")
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Info().Err(err).Msg("Got non 2XX header")
		return fmt.Errorf("UNSUCCESSFUL REQUEST")
	}
	return nil
}

func (client *EventSourcingHttpClient) GetEventsOrdered(aggregateId string) (*EventsIterator, error) {

	getEventsUrl, err := url.JoinPath(client.url, fmt.Sprintf("/%s/events", aggregateId))
	if err != nil {
		log.Info().Err(err).Msg("Could not use url")
		return nil, err
	}

	resp, err := client.httpClient.Get(getEventsUrl)
	if err != nil {
		log.Info().Err(err).Msg("Error during the request")
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Info().Err(err).Msg("Got non 2XX header")
		return nil, fmt.Errorf("UNSUCCESSFUL REQUEST")
	}
	var events []models.Event
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Info().Err(err).Msg("Error during reading response body")
		return nil, err
	}
	err = json.Unmarshal(buf, &events)

	if err != nil {
		log.Info().Err(err).Msg("Error during unmarshalling body")
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
