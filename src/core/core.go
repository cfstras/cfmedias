package core

import (
	"errrs"
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

type CommandContext struct {
	Cmd       string
	Args      ArgMap
	AuthLevel AuthLevel
}

type Command struct {
	Verbs        []string
	Help         string
	MinAuthLevel AuthLevel
	Handler      func(ctx CommandContext) Result
}

type Status string

const (
	StatusOK             Status = "OK"
	StatusItemNotFound          = "ItemNotFound"
	StatusError                 = "Error"
	StatusQueryNotUnique        = "QueryNotUnique"
)

type Result struct {
	Status  Status        `json:"status"`
	Results []interface{} `json:"results,omitempty"`
	Error   error         `json:"error,omitempty"`
}

var (
	ResultOK = Result{Status: StatusOK, Results: nil, Error: nil}

	ErrorCmdNotFound = errrs.New("Command not found!")
	ErrorNotAllowed  = errrs.New("You are not allowed to do that!")
	ErrorNotLoggedIn = errrs.New("You are not allowed to do that; you need to be logged in!")

	//ErrorItemNotFound = errrs.New("The requested item was not found.")
)

func ResultByError(err error) Result {
	if err == nil {
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
	Cmd(ctx CommandContext) Result
	IsCmdAllowed(level AuthLevel, cmd string) (bool, error)
}
