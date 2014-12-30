package main

import (
	"os"

	"github.com/cfstras/cfmedias/coreimpl"
	log "github.com/cfstras/cfmedias/logger"
)

func main() {

	//TODO check args

	// start engine
	log.Log.Println("Starting cfmedias")

	inst := coreimpl.New()
	err := inst.Start()

	if err != nil {
		log.Log.Println("Caught an error:", err.Error(), "- exiting...")
		err = inst.Shutdown()
		if err != nil {
			log.Log.Println(err.Error())
		}
		os.Exit(1)
	}
	//go sineLogon()

	// listen for commands
	inst.CmdLine()

	// CmdLine is finished, shutdown
	log.Log.Println("Exiting...")

	if err = inst.Shutdown(); err != nil {
		log.Log.Println(err.Error())
	}
}
