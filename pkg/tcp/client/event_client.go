package client

import (
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type TcpEventClient struct {
	conn net.Conn
}

func NewTcpEventClient() (*TcpEventClient, error) {
	tcpEv := TcpEventClient{}
	conn, err := tcpEv.setup()
	if err != nil {
		return nil, err
	}
	tcpEv.conn = conn
	return &tcpEv, nil
}

func (tcpEv *TcpEventClient) setup() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:5521")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to server")
		return nil, err
	}
	return conn, nil

}

func (tcpEv *TcpEventClient) ListenForEvents(channel chan string) {
	for {
		buffer := make([]byte, 128)
		n, err := tcpEv.conn.Read(buffer)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read from connection")
			time.Sleep(1 * time.Second)
			tcpEv.conn.Close()
			tcpEv.setup()
			continue
		}
		message := strings.TrimRight(string(buffer[:n]), "\x00")
		log.Info().Msg(message)
		channel <- message
	}
}
