package server

import (
	"errors"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
)

type TcpEventServer struct {
	consumer []net.Conn
	mu       sync.Mutex
	listener net.Listener
}

func NewTcpEventServer() (*TcpEventServer, error) {

	tcpServer := &TcpEventServer{
		consumer: []net.Conn{},
		mu:       sync.Mutex{},
	}
	err := tcpServer.setup()
	if err != nil {
		return nil, err
	}
	return tcpServer, nil
}

func (tcpServer *TcpEventServer) setup() error {
	listener, err := net.Listen("tcp", "0.0.0.0:5521")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
		return err
	}
	if tcpServer.listener != nil {
		tcpServer.Stop()
	}
	tcpServer.listener = listener
	return nil
}

func (tcpServer *TcpEventServer) Stop() {
	err := tcpServer.listener.Close()
	if err != nil {
		log.Debug().Err(err).Msg("Failed to close listener")
	}
}

func (tcpServer *TcpEventServer) Start() {
	log.Debug().Msg("Starting tcp server")
	for {
		conn, err := tcpServer.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Err(err).Msg("Listener closed")
				log.Error().Err(err).Msg("Failed retry setup listener")
				return
			}
			log.Error().Err(err).Msg("Failed to accept connection")
			continue
		}
		tcpServer.mu.Lock()
		tcpServer.consumer = append(tcpServer.consumer, conn)
		tcpServer.mu.Unlock()
		log.Trace().Msg("Accepted connection")
	}
}

func removenullvalue(slice []net.Conn) []net.Conn {
	output := []net.Conn{}
	for _, element := range slice {
		if element != nil {
			output = append(output, element)
		}
	}
	return output
}

func (tcpServer *TcpEventServer) SendEvent(event string) error {
	if len([]byte(event)) > 128 {
		return errors.New("event size exceeds 1024 bytes")
	}
	eventBytes := make([]byte, 128)
	copy(eventBytes, []byte(event))
	for i, conn := range tcpServer.consumer {
		log.Trace().Msgf("Sending event to consumer %d", i)
		_, err := conn.Write(eventBytes)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send event - closing connection")
			conn.Close()
			tcpServer.consumer[i] = nil
		} else {

			log.Trace().Msgf("Sent event to consumer %d", i)
		}
	}
	tcpServer.mu.Lock()
	defer tcpServer.mu.Unlock()
	arr := removenullvalue(tcpServer.consumer)
	tcpServer.consumer = arr
	return nil
}
