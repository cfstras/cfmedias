package core

import (
	"config"
	"db"
	log "logger"
	"os"
	"os/signal"
)

var shutUp bool
var signals chan os.Signal

func Start() error {

	// load config
	err := config.Load("config.json")
	if err != nil {
		return err
	}

	// set up signal handler
	signals = make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		for _ = range signals {
			// interrupted!
			Shutdown()
		}
	}()

	// connect to db
	if err = db.Open(); err != nil {
		return err
	}
	shutUp = true

	initCmd()

	//TODO call plugin loads

	// update local files
	go db.Update()

	return nil
}

func Shutdown() error {
	if !shutUp {
		return nil
	}

	log.Log.Println("shutting down.")

	// disconnect from db
	if err := db.Close(); err != nil {
		log.Log.Println("Error closing database:", err)
	}

	// save config
	err := config.Save("config.json")
	if err != nil {
		//TODO don't catch if this is an init error
		log.Log.Println("Error while saving config:", err.Error())
	}
	shutUp = false

	if err := exitCmd(); err != nil {
		log.Log.Println("cmd exit error", err)
	}

	return err
}
