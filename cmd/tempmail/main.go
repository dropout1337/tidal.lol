package main

import (
	"tidal.lol/internal/api"
	"tidal.lol/internal/database"
	"tidal.lol/internal/emails"
	"tidal.lol/internal/logging"
	"tidal.lol/internal/smtpd"
)

var (
	// smtpPort to listen on
	smtpPort  = 25
	smtpDebug = true

	// httpPort to listen on
	httpPort  = 80
	httpDebug = true

	// databaseURL, databaseUsername, and databasePassword are used to connect to the database
	databaseURL      = "ws://localhost:8000/rpc"
	databaseUsername = "root"
	databasePassword = "root"
)

func main() {
	database.DB = database.NewDatabase(databaseURL, databaseUsername, databasePassword)

	go func() {
		logging.Logger.Info().
			Int("port", smtpPort).
			Msg("Started SMTP listener")

		err := smtpd.NewServer(smtpPort, emails.CreateHandler(smtpDebug))
		if err != nil {
			logging.Logger.Fatal().
				Err(err).
				Msg("Failed to start the SMTP server")
		}
	}()

	logging.Logger.Info().
		Int("port", httpPort).
		Msg("Started the API server")

	err := api.NewServer(httpPort, httpDebug)
	if err != nil {
		logging.Logger.Fatal().
			Err(err).
			Msg("Failed to start the API server")
	}
}
