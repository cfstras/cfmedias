package core

import (
	"github.com/cfstras/cfmedias/errrs"
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
	Cmd  string // Command given
	Args ArgMap // arguments of the request

	AuthLevel AuthLevel // permission level of current user
	UserId    *int64    // logged in User ID or nil
}

type Command struct {
	Verbs        []string
	Description  string
	ArgsHelp     map[string]string
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

type JobSignal int

const (
	SignalNone JobSignal = iota
	SignalTerminate
	SignalKill
)

type Result struct {
	Status Status      `json:"status"`
	Result interface{} `json:"result,omitempty"`
	Error  error       `json:"error,omitempty"`
	IsRaw  bool        `json:"-"`
}

var (
	ResultOK = Result{Status: StatusOK, Result: nil, Error: nil}

	ErrorCmdNotFound    = errrs.New("Command not found!")
	ErrorNotAllowed     = errrs.New("You are not allowed to do that!")
	ErrorNotLoggedIn    = errrs.New("You are not allowed to do that; you need to be logged in!")
	ErrorNotImplemented = errrs.New("Sorry, this feature is not implemented yet.")
	ErrorUserNotFound   = errrs.New("User not found!")
	ErrorInvalidQuery   = errrs.New("Invalid Query!")

	//ErrorItemNotFound = errrs.New("The requested item was not found.")
)

func ResultByError(err error) Result {
	if err == nil {
		return ResultOK
	} else {
		return Result{Status: StatusError, Error: errrs.New(err.Error())}
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

	// Long-running jobs can use this method to register a shutdown handler.
	// Returned is a channel which should be listened on.
	//
	// When the core is shut down, it will first send a SignalTerminate to each
	// registered job, then wait a few seconds, and then send a SignalKill.
	// After the SignalKill is received, the goroutine will be killed.
	// If necessary, the killing can be delayed by not listening on the channel
	// anymore after receiving the SignalTerminate. Don't wait for too long,
	// though.
	RegisterJob() <-chan JobSignal

	// Use this to unregister your job once it's finished.
	UnregisterJob(job <-chan JobSignal)
}
