package core

import (
	"errors"
)

type ArgMap map[string][]string

type Command struct {
	Verbs   []string
	Help    string
	Handler func(args ArgMap) error
}

var (
	ErrorCmdNotFound = errors.New("Command not found!")
)

type Core interface {
	Start() error
	Shutdown() error
	Version() string

	RegisterCommand(Command)
	UnregisterCommand(Command)
	CmdLine()
	Cmd(cmd string, args ArgMap) error
}
