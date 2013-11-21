package main

import (
	"core"
	"log"
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
	}

	// listen for commands
	core.CmdLine()
}
