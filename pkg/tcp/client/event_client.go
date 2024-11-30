package client

import (
	"net"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type TcpEventClient struct {
	conn      net.Conn
	clientURL string
}

func NewTcpEventClient() (*TcpEventClient, error) {
	clientURL := os.Getenv("EVENT_SOURCING_CLIENT_TCP")
	if clientURL == "" {
		clientURL = "localhost:5521"
		log.Debug().Msgf("EVENT_SOURCING_CLIENT_TCP not set, defaulting to %s", clientURL)
	} else {
		log.Debug().Msgf("Using EVENT_SOURCING_CLIENT_TCP: %s", clientURL)
	}
	tcpEv := TcpEventClient{
		clientURL: clientURL,
	}
	log.Debug().Msg("Setting up client")
	err := tcpEv.setup(10)
	if err != nil {
		return nil, err
	}

	return &tcpEv, nil
}

func (tcpEv *TcpEventClient) setup(retries int) error {
	if retries <= 0 {
		log.Fatal().Msg("Exceeded maximum reconnection attempts")
		panic("Exceeded maximum reconnection attempts")

	}
	if tcpEv.conn != nil {

		tcpEv.conn.Close()
	}

	conn, err := net.Dial("tcp", tcpEv.clientURL)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to connect to server, retries left: %d", retries-1)
		time.Sleep(1 * time.Second)
		return tcpEv.setup(retries - 1)
	}
	log.Debug().Msg("Connected to server")
	tcpEv.conn = conn
	return nil
}

func (tcpEv *TcpEventClient) ListenForEvents(channel chan string) {
	for {
		buffer := make([]byte, 128)
		n, err := tcpEv.conn.Read(buffer)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read from connection")
			tcpEv.setup(10)
			continue
		}
		message := strings.TrimRight(string(buffer[:n]), "\x00")
		log.Info().Msg(message)
		channel <- message
	}
}
