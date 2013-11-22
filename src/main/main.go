package main

import (
	"core"
	"log"
	"os"
)

func main() {

	//TODO check args

	// start engine
	log.Println("Starting cfmedias")
	err := core.Start()

	if err != nil {
		log.Println("Caught an error:", err.Error(), "- exiting...")
		err = core.Shutdown()
		if err != nil {
			log.Println(err.Error())
		}
		os.Exit(1)
	}
	go sineLogon()

	// listen for commands
	core.CmdLine()

	// CmdLine is finished, shutdown
	log.Println("Finished, exiting...")
	err = core.Shutdown()
	if err != nil {
		log.Println(err.Error())
	}
}
