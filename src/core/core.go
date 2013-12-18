package core

import (
	"errors"
	"io"
)

type ArgMap map[string][]string

type Command struct {
	Verbs   []string
	Help    string
	Handler func(args ArgMap, w io.Writer) error
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
	Cmd(cmd string, args ArgMap, w io.Writer) error
}
