package core

import (
	"config"
	"db"
	log "logger"
)

var shutUp bool

func Start() error {

	// load config
	err := config.Load("config.json")
	if err != nil {
		return err
	}
	// connect to db
	if err = db.Open(); err != nil {
		return err
	}
	shutUp = true

	initCmd()

	//TODO call plugin loads

	// update local files
	db.Update()

	return nil
}

func Shutdown() error {
	if !shutUp {
		return nil
	}

	shutUp = false
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

	return err
}
