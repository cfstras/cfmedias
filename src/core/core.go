package core

import (
	"config"
	//	"db"
	"fmt"
)

func Start() error {
	// connect to db
	err := config.Load("config.json")
	if err != nil {
		return err
	}
	return nil
}

// start a REPL shell.
func CmdLine() {
	fmt.Println("cfmedias", currentVersion)
	//TODO
}
