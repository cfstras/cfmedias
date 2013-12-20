package coreimpl

import (
	"core"
	"fmt"
	"github.com/peterh/liner"
	log "logger"
	"os"
	"strings"
)

// inits the cmd subsystem
func (c *impl) initCmd() {
	c.commandMap = make(map[string]core.Command)
	c.commandSet = make(map[string]core.Command)
	c.registerBaseCommands()
}

// exits the cmd line
func (c *impl) exitCmd() error {
	c.replActive = false
	c.repl.Close()

	if c.reading {
		fmt.Println("Press enter to continue...")
	}
	err := os.Stdin.Close()
	if err != nil {
		fmt.Println("close error", err)
	}
	return nil
}

// add a command to the list of available commands
func (c *impl) RegisterCommand(command core.Command) {
	for _, verb := range command.Verbs {
		old, already := c.commandMap[verb]
		if already {
			fmt.Println("error registering verb", verb, `for command "`,
				command.Help, `", it already exists with command "`, old.Help, `".`)
			return
		}
	}
	for _, verb := range command.Verbs {
		c.commandMap[verb] = command
	}
	c.commandSet[command.Verbs[0]] = command
}

// remove a command from the available list
func (c *impl) UnregisterCommand(command core.Command) {
	for _, verb := range command.Verbs {
		delete(c.commandMap, verb)
	}
	delete(c.commandSet, command.Verbs[0])
}

func (c *impl) registerBaseCommands() {
	c.RegisterCommand(core.Command{
		[]string{"quit", "q", "close", "exit"},
		"Shuts down and exits.",
		core.AuthRoot,
		func(_ core.ArgMap) core.Result {
			return core.ResultByError(c.Shutdown())
		}})

	c.RegisterCommand(core.Command{
		[]string{"help", "h", "?"},
		"Prints help.",
		core.AuthGuest,
		func(_ core.ArgMap) core.Result {
			res := make(map[string]interface{})
			for k, v := range c.commandSet {
				res[k] = v.Help
			}
			return core.Result{Status: core.StatusOK, Results: []interface{}{res}}
		}})

}

const maxUnicodeString = "\U0010FFFF"

// start a REPL shell.
func (c *impl) CmdLine() {
	log.Log.Println("cfmedias", c.currentVersion)

	c.repl = liner.NewLiner()
	c.repl.SetCompleter(c.completer)

	for c.replActive = true; c.replActive; {
		c.reading = true
		cmd, err := c.repl.Prompt("> ")
		//c.repl.Close()
		c.reading = false
		if err != nil && c.replActive {
			fmt.Println(err)
			c.replActive = false
			break
		}
		if !c.replActive {
			return
		}

		cmd = strings.Replace(cmd, `\ `, maxUnicodeString, -1)
		cmd = strings.Replace(cmd, `\\`, `\`, -1)
		split := strings.Split(cmd, " ")

		if len(split) > 0 && len(cmd) > 0 {
			// convert arg list to map, using format
			// name=max fruits=apple fruits=orange
			// ==> map[name: [max], fruits: [apple, orange]]
			args := make(core.ArgMap)
			for _, e := range split[1:] {
				e = strings.Replace(e, maxUnicodeString, ` `, -1)
				tuple := strings.Split(e, "=")
				if _, ok := args[tuple[0]]; !ok {
					args[tuple[0]] = make([]string, 0)
				}
				if len(tuple) > 0 {
					args[tuple[0]] = append(args[tuple[0]],
						strings.Join(tuple[1:], "="))
				}
			}

			result := c.Cmd(split[0], args, core.AuthRoot)
			if result.Error != core.ErrorCmdNotFound {
				c.repl.AppendHistory(cmd)
			}
			log.Log.Println(result)
		}
	}
}

func (c *impl) Cmd(cmd string, args core.ArgMap, level core.AuthLevel) core.Result {
	command, ok := c.commandMap[cmd]
	if !ok {
		return core.Result{Error: core.ErrorCmdNotFound}
	}
	if level < command.MinAuthLevel {
		return core.Result{Error: core.ErrorNotAllowed}
	}

	return command.Handler(args)
}

func (c *impl) completer(s string) []string {
	out := make([]string, 0)
	// walk cmd map
	for k, _ := range c.commandMap {
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
