package main

import (
	"core"
	"fmt"
)

func main() {

	//TODO check args

	// start engine
	fmt.Println("Starting cfmedias")
	err := core.Start()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Exiting")
	}

	// listen for commands
	core.CmdLine()
}
