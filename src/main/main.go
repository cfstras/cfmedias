package main

import (
	"coreimpl"
	"log"
	"os"
)

func main() {

	//TODO check args

	// start engine
	log.Println("Starting cfmedias")

	inst := coreimpl.New()
	err := inst.Start()

	if err != nil {
		log.Println("Caught an error:", err.Error(), "- exiting...")
		err = inst.Shutdown()
		if err != nil {
			log.Println(err.Error())
		}
		os.Exit(1)
	}
	//go sineLogon()

	// listen for commands
	inst.CmdLine()

	// CmdLine is finished, shutdown
	log.Println("Exiting...")

	if err = inst.Shutdown(); err != nil {
		log.Println(err.Error())
	}
}
