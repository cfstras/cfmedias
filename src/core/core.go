package core

import (
	"errors"
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
	Handler      func(args ArgMap) Result
}

type Status string

const (
	StatusOK             Status = "OK"
	StatusItemNotFound          = "ItemNotFound"
	StatusError                 = "Error"
	StatusQueryNotUnique        = "QueryNotUnique"
)

type Result struct {
	Status  Status
	Results []interface{}
	Error   error
}

var (
	ResultOK = Result{Status: StatusOK, Results: nil, Error: nil}

	ErrorCmdNotFound = errors.New("Command not found!")
	ErrorNotAllowed  = errors.New("You are not allowed to do that!")

	//ErrorItemNotFound = errors.New("The requested item was not found.")
)

func ResultByError(err error) Result {
	if err != nil {
		return ResultOK
	} else {
		return Result{Status: StatusError, Error: err}
	}
}

type Core interface {
	Start() error
	Shutdown() error
	Version() string

	RegisterCommand(Command)
	UnregisterCommand(Command)
	CmdLine()
	Cmd(cmd string, args ArgMap, level AuthLevel) Result
}
