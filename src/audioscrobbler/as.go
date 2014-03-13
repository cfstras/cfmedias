package audioscrobbler

import (
	"core"
	"db"
)

type AS struct {
	db *db.DB
}

func (a *AS) Start(c core.Core, d *db.DB) {
	a.db = d
	c.RegisterCommand(core.Command{
		[]string{"audioscrobbler", "as"},
		`Interface for clients using the Last.fm Submissions Protocol v1.2.1.
		See also http://www.last.fm/api/submissions.`,
		map[string]string{
			"hs": "hs=true marks a handshake request. requires p, c, v, u, t, a.",
			"p":  "must be 1.2.1",
			"c":  "identifier for the client software",
			"v":  "version of the client software",
			"u":  "username",
			"t":  "timestamp of the request",
			"a":  "authentication token: md5(md5(password) + timestamp)",
			"response": `The server will respond with:
			OK
			<session ID>
			<url to be used for a now-playing request>
			<url to be used for a submission request>

			or one of these:
			BANNED: Your client is doing bad things and is therefore banned.
			BADAUTH: Authentication details are incorrect.
			BADTIME: Time difference between server and client time too high.
				Correct your system clock.
			FAILED <reason>: Your request failed for the specified reason.`},
		core.AuthUser,
		a.LoginAS})
}

func (a *AS) LoginAS(core.CommandContext) core.Result {
	//TODO implement

	//success, authToken, err := a.db.Login(name, password)

	return core.ResultByError(core.ErrorNotImplemented)
}
