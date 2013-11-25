package core

import (
	"bytes"
	"config"
	"db"
	"fmt"
	"log"
	"os"
	"strings"
)

type Command struct {
	verbs   []string
	help    string
	handler func(args []string)
}

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
	commandMap = make(map[string]Command)
	registerBaseCommands()
	//TODO call plugin loads

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

var commandMap map[string]Command

func RegisterCommand(command Command) {
	for _, verb := range command.verbs {
		old, already := commandMap[verb]
		if already {
			fmt.Println("error registering verb", verb, `for command "`,
				command.help, `", it already exists with command "`, old.help, `".`)
			return
		}
	}
	for _, verb := range command.verbs {
		commandMap[verb] = command
	}
}

func UnregisterCommand(command Command) {
	for _, verb := range command.verbs {
		delete(commandMap, verb)
	}
}

func registerBaseCommands() {
	quit := Command{
		[]string{"quit", "q", "close", "exit"},
		"Shuts down and exits.",
		func(_ []string) {
			Shutdown()
			runCmdLine = false
			os.Stdin.Close()
		}}
	RegisterCommand(quit)

}

var runCmdLine = true

// start a REPL shell.
func CmdLine() {
	log.Println("cfmedias", currentVersion)

	readBuffer := make([]byte, 1)
	var buffer bytes.Buffer

	for runCmdLine {
		fmt.Print("> ")
		buffer.Reset()
		readBuffer[0] = 0
		for readBuffer[0] != '\n' {
			n, err := os.Stdin.Read(readBuffer)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if n != 1 {
				fmt.Println("read error!,", readBuffer[0])
				continue
			}
			if readBuffer[0] != '\n' && readBuffer[0] != '\r' {
				buffer.WriteByte(readBuffer[0])
			}
		}
		cmd := buffer.String()
		// Parse!
		split := strings.Split(cmd, " ")
		if len(split) > 0 {
			command, ok := commandMap[split[0]]
			if !ok {
				fmt.Println("Error: no command for", split[0])
				continue
			}
			command.handler(split[1:])
		}
	}
}

func parseQuoted(s string) string {
	parsed := s[1 : len(s)-1]
	parsed = strings.Replace(parsed, `\"`, `"`, -1)
	return parsed
}
