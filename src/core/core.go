package core

import (
	"config"
	"db"
	"fmt"
	"log"
)

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

	// update local files
	db.Update()
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

	for run := true; run; {
		fmt.Print("> ")
		var str string
		_, err := fmt.Scanln(&str)
		if err != nil && err.Error() == "unexpected newline" {

		} else if err != nil {
			fmt.Println("error:", err.Error())
		} else {
			//TODO parse
			fmt.Println("read", str)
		}
	}

	//TODO
}
