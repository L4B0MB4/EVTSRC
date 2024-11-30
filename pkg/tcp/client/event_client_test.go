package client

import (
	"testing"
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/tcp/server"
	"github.com/stretchr/testify/assert"
)

func TestTcpEventClientServerIntegration(t *testing.T) {
	server, err := server.NewTcpEventServer()
	assert.NoError(t, err)
	go server.Start()

	client, err := NewTcpEventClient()
	assert.NoError(t, err)

	eventChannel := make(chan string)
	go client.ListenForEvents(eventChannel)

	testEvent := "Test Event"
	err = server.SendEvent(testEvent)
	assert.NoError(t, err)

	select {
	case receivedEvent := <-eventChannel:
		assert.Equal(t, testEvent, receivedEvent)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for event")
	}
}

func TestTcpEventClientServerMultipleReads(t *testing.T) {
	server, err := server.NewTcpEventServer()
	assert.NoError(t, err)
	go server.Start()

	client, err := NewTcpEventClient()
	assert.NoError(t, err)

	eventChannel := make(chan string)
	go client.ListenForEvents(eventChannel)

	testEvents := []string{"Event 1", "Event 2", "Event 3"}
	for _, event := range testEvents {
		err = server.SendEvent(event)
		assert.NoError(t, err)
	}

	for _, expectedEvent := range testEvents {
		select {
		case receivedEvent := <-eventChannel:
			assert.Equal(t, expectedEvent, receivedEvent)
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for event")
		}
	}
}
