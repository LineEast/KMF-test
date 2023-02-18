package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"kmfRedirect/internal/configuration"
	"kmfRedirect/internal/database"
	"kmfRedirect/internal/server"
)

func main() {
	configuration, err := configuration.New()
	if err != nil {
		log.Fatal(err)
	}

	database, err := database.Conn(configuration.Database)
	if err != nil {
		log.Fatal(err)
	}

	errs := make(chan error)
	go func() {
		errs <- server.New(database, configuration.Server).Run()
	}()

	log.Println("Starting on " + "'http:\\\\" + configuration.Server.Host + ":" + configuration.Server.Port + "'")

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err = <-errs:
		if err != nil {
			log.Fatal(err)
		}
	case <-signals:
	}
}
