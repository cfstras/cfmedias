package core

import (
	"bytes"
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

// inits the cmd subsystem
func initCmd() {
	commandMap = make(map[string]Command)
	registerBaseCommands()
}

// stores the loaded commands, sorted by verb.
// multiple verbs may point to the same command.
var commandMap map[string]Command

// whether the cmd line is currently active
var runCmdLine bool

// add a command to the list of available commands
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

// remove a command from the available list
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

// start a REPL shell.
func CmdLine() {
	log.Println("cfmedias", currentVersion)

	runCmdLine = true
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
