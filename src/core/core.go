package core

import (
	"errors"
)

type Command struct {
	Verbs   []string
	Help    string
	Handler func(args []string) error
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
	Cmd(cmd string, args []string) error
}
