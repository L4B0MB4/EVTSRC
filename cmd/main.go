package main

import (
	"os"

	"github.com/L4B0MB4/EVTSRC/pkg/httphandler"
	"github.com/L4B0MB4/EVTSRC/pkg/httphandler/controller"
	"github.com/L4B0MB4/EVTSRC/pkg/store"
	"github.com/L4B0MB4/EVTSRC/pkg/tcp/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	db := store.DatabaseConnection{}
	db.SetUp()
	conn, err := db.GetDbConnection()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessfull initalization of db")
		return
	}
	log.Debug().Msg("Db Connection was successful")
	repository := store.NewEventRepository(conn)

	tcpServer, err := server.NewTcpEventServer()
	if err != nil {
		log.Error().Err(err).Msg("Unsuccessfull initalization of tcp server")
		return
	}
	go tcpServer.Start()

	c := controller.NewEventController(repository, tcpServer)
	h := httphandler.NewHttpHandler(c)

	h.Start()
}
