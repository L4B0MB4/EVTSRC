package client

import (
	"os"

	"github.com/rs/zerolog/log"
)

func RetrieveEventSourcingClientUrl() string {
	clientURL := os.Getenv("EVENT_SOURCING_CLIENT_URL")
	if clientURL == "" {
		clientURL = "http://localhost:5515"
		log.Debug().Msgf("EVENT_SOURCING_CLIENT_URL not set, defaulting to %s", clientURL)
	} else {
		log.Debug().Msgf("Using EVENT_SOURCING_CLIENT_URL: %s", clientURL)
	}

	return clientURL
}
