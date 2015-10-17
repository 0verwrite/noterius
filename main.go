package main

import (
	"github.com/Nyarum/noterius/core"
	"github.com/Nyarum/noterius/land"
	log "github.com/Sirupsen/logrus"

	"flag"
)

func main() {
	configPathFlag := flag.String("config", "resource/config.yml", "Config file for start server")
	dbIPFlag := flag.String("dbip", "", "Set IP for database")
	flag.Parse()

	var (
		err error
		app land.Application = land.Application{}
	)
	defer core.ErrorGlobalHandler()

	log.Info("Loading logger..")
	core.NewLogger()

	log.Info("Loading config..")
	if app.Config, err = core.NewConfig(*configPathFlag); err != nil {
		log.WithError(err).Panic("Config is not load")
	}

	if *dbIPFlag != "" {
		app.Config.Database.IP = *dbIPFlag
	}

	log.Info("Loading database..")
	if app.Database, err = core.NewDatabase(&app.Config); err != nil {
		log.WithError(err).Panic("Database is not load")
	}

	log.WithField("address", app.Config.Base.IP+":"+app.Config.Base.Port).Info("Server starting")
	if err := app.Run(); err != nil {
		log.WithError(err).Panic("Server is not started")
	}
}
