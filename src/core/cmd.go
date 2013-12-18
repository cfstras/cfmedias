package core

import (
	"db"
	"errors"
	"fmt"
	"github.com/peterh/liner"
	log "logger"
	"os"
	"strings"
)

type Command struct {
	Verbs   []string
	Help    string
	Handler func(args []string) error
}

// inits the cmd subsystem
func initCmd() {
	CommandMap = make(map[string]Command)
	CommandSet = make(map[string]Command)
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
var CommandMap map[string]Command

// stores the loaded commands, sorted by first verb
var CommandSet map[string]Command

// the REPL state
var repl *liner.State
var replActive bool
var reading bool

// add a command to the list of available commands
func RegisterCommand(command Command) {
	for _, verb := range command.Verbs {
		old, already := CommandMap[verb]
		if already {
			fmt.Println("error registering verb", verb, `for command "`,
				command.Help, `", it already exists with command "`, old.Help, `".`)
			return
		}
	}
	for _, verb := range command.Verbs {
		CommandMap[verb] = command
	}
	CommandSet[command.Verbs[0]] = command
}

// remove a command from the available list
func UnregisterCommand(command Command) {
	for _, verb := range command.Verbs {
		delete(CommandMap, verb)
	}
	delete(CommandSet, command.Verbs[0])
}

func registerBaseCommands() {
	RegisterCommand(Command{
		[]string{"quit", "q", "close", "exit"},
		"Shuts down and exits.",
		func(_ []string) error {
			return Shutdown()
		}})

	RegisterCommand(Command{
		[]string{"help", "h", "?"},
		"Prints help.",
		func(_ []string) error {
			fmt.Println("Available commands:")
			for k, v := range CommandSet {
				fmt.Println(" ", k, "-", v.Help)
			}
			return nil
		}})

	RegisterCommand(Command{
		[]string{"rescan"},
		"Refreshes the database by re-scanning the music folder.",
		func(_ []string) error {
			db.Update()
			return nil
		}})

	RegisterCommand(Command{
		[]string{"stats"},
		"Prints some statistics about the database",
		func(_ []string) error {
			fmt.Printf(" %7s %7s %7s\n", "Titles", "Folders", "Titles/Folder")
			fmt.Printf(" %7d %7d %7f\n", db.TitlesTotal(), db.FoldersTotal(),
				db.AvgFilesPerFolder())
			return nil
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
			err = Cmd(split[0], split[1:])
			if err != nil {
				log.Log.Println(err)
			} else {
				repl.AppendHistory(cmd)
			}
		}
	}
}

func Cmd(cmd string, args []string) error {
	command, ok := CommandMap[cmd]
	if !ok {
		return errors.New("Error: no command for " + cmd)
	}

	return command.Handler(args)
}

func completer(s string) []string {
	out := make([]string, 0)
	// walk cmd map
	for k, _ := range CommandMap {
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
