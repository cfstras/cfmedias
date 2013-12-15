package core

import (
	"db"
	"fmt"
	"github.com/peterh/liner"
	log "logger"
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
	commandSet = make(map[string]Command)
	registerBaseCommands()
}

// exits the cmd line
func exitCmd() error {
	replActive = false
	repl.Close()

	if reading {
		fmt.Println("Press enter to continue...")
	}
	err := os.Stdin.Close()
	if err != nil {
		fmt.Println("close error", err)
	}
	return nil
}

// stores the loaded commands, sorted by verb.
// multiple verbs may point to the same command.
var commandMap map[string]Command

// stores the loaded commands, sorted by first verb
var commandSet map[string]Command

// the REPL state
var repl *liner.State
var replActive bool
var reading bool

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
	commandSet[command.verbs[0]] = command
}

// remove a command from the available list
func UnregisterCommand(command Command) {
	for _, verb := range command.verbs {
		delete(commandMap, verb)
	}
	delete(commandSet, command.verbs[0])
}

func registerBaseCommands() {
	RegisterCommand(Command{
		[]string{"quit", "q", "close", "exit"},
		"Shuts down and exits.",
		func(_ []string) {
			Shutdown()
		}})

	RegisterCommand(Command{
		[]string{"help", "h", "?"},
		"Prints help.",
		func(_ []string) {
			fmt.Println("Available commands:")
			for k, v := range commandSet {
				fmt.Println(" ", k, "-", v.help)
			}
		}})

	RegisterCommand(Command{
		[]string{"rescan"},
		"Refreshes the database by re-scanning the music folder.",
		func(_ []string) {
			db.Update()
		}})

	RegisterCommand(Command{
		[]string{"stats"},
		"Prints some statistics about the database",
		func(_ []string) {
			fmt.Printf(" %7s %7s %7s\n", "Titles", "Folders", "Titles/Folder")
			fmt.Printf(" %7d %7d %7f\n", db.TitlesTotal(), db.FoldersTotal(),
				db.AvgFilesPerFolder())
		}})
}

// start a REPL shell.
func CmdLine() {
	log.Log.Println("cfmedias", currentVersion)

	repl = liner.NewLiner()
	repl.SetCompleter(completer)

	for replActive = true; replActive; {
		reading = true
		cmd, err := repl.Prompt("> ")
		reading = false
		if err != nil && replActive {
			fmt.Println(err)
			replActive = false
			break
		}
		if !replActive {
			return
		}
		split := strings.Split(cmd, " ")

		if len(split) > 0 && len(cmd) > 0 {
			command, ok := commandMap[split[0]]
			if !ok {
				fmt.Println("Error: no command for", split[0])
				continue
			}
			repl.AppendHistory(cmd)
			command.handler(split[1:])
		}
	}
}

func completer(s string) []string {
	out := make([]string, 0)
	// walk cmd map
	for k, _ := range commandMap {
		if strings.HasPrefix(k, s) {
			out = append(out, k)
		}
	}
	return out
}

func parseQuoted(s string) string {
	parsed := s[1 : len(s)-1]
	parsed = strings.Replace(parsed, `\"`, `"`, -1)
	return parsed
}
