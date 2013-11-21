package core

import (
	"config"
	//	"db"
	"log"
)

func Start() error {
	// connect to db
	err := config.Load("config.json")
	if err != nil {
		return err
	}
	return nil
}

func Shutdown() error {
	// disconnect from db
	// save config
	err := config.Save("config.json")
	if err != nil {
		//TODO don't catch if this is an init error
		log.Println("Error while saving config:", err.Error())
	}

	return err
}

// start a REPL shell.
func CmdLine() {
	log.Println("cfmedias", currentVersion)
	//TODO
}
