package core

import (
	"errors"
	"io"
)

type ArgMap map[string][]string

type AuthLevel uint

const (
	// Guests can only view public data
	AuthGuest AuthLevel = iota

	// Users can view and edit their own data
	AuthUser

	// Admins can edit other user data and most configuration
	AuthAdmin

	// Roots can edit all configuration, do database maintenance
	// and stop the server
	AuthRoot
)

type Command struct {
	Verbs        []string
	Help         string
	MinAuthLevel AuthLevel
	Handler      func(args ArgMap, w io.Writer) error
}

var (
	ErrorCmdNotFound = errors.New("Command not found!")
	ErrorNotAllowed  = errors.New("You are not allowed to do that!")
)

type Core interface {
	Start() error
	Shutdown() error
	Version() string

	RegisterCommand(Command)
	UnregisterCommand(Command)
	CmdLine()
	Cmd(cmd string, args ArgMap, w io.Writer, level AuthLevel) error
}
