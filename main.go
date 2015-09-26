package main

import (
	"github.com/Nyarum/noterius/core"
	"github.com/Nyarum/noterius/land"
	log "github.com/Sirupsen/logrus"

	"flag"
)

func main() {
	configPathFlag := flag.String("config", "resource/config.yml", "A config file for start server")
	dbIPFlag := flag.String("dbip", "", "IP for database")
	flag.Parse()

	app := land.Application{}
	defer core.ErrorGlobalHandler()

	log.Info("Loading logger..")
	core.NewLogger()

	log.Info("Loading config..")
	if err := core.NewConfig(&app.Config, *configPathFlag); err != nil {
		log.WithError(err).Panic("Config is not load")
	}

	if *dbIPFlag != "" {
		app.Config.Database.IP = *dbIPFlag
	}

	log.Info("Loading database..")
	if err := core.NewDatabase(&app.Database, &app.Config); err != nil {
		log.WithError(err).Panic("Database is not load")
	}

	log.WithField("address", app.Config.Base.IP+":"+app.Config.Base.Port).Info("Server starting")
	if err := app.Run(); err != nil {
		log.WithError(err).Panic("Server is not started")
	}
}
