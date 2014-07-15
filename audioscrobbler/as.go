package audioscrobbler

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/db"
	"github.com/cfstras/cfmedias/errrs"
	"github.com/cfstras/cfmedias/logger"
	"github.com/cfstras/cfmedias/util"
)

const (
	StatusOK      core.Status = "OK"
	StatusBanned              = "BANNED"
	StatusBadAuth             = "BADAUTH"
	StatusBadTime             = "BADTIME"
	StatusFailed              = "FAILED"
)

type AS struct {
	db            *db.DB
	nowPlayingURL string
	submissionURL string
}

func (a *AS) Start(c core.Core, d *db.DB) {
	a.db = d

	//TODO get base url
	baseURL := "http://localhost/api/"
	a.nowPlayingURL = baseURL + "as_nowplaying"
	a.submissionURL = baseURL + "as_scrobble"

	c.RegisterCommand(core.Command{
		[]string{"audioscrobbler", "as"},
		`Interface for clients using the Last.fm Submissions Protocol v1.2.1.
In your favourite scrobbler, use this as handshake URL:
http://cfmedias-server-ip/api/audioscrobbler/?hs=true
See also http://www.last.fm/api/submissions.`,
		map[string]string{
			"hs": "hs=true marks a handshake request. requires p, c, v, u, t, a.",
			"p":  "must be 1.2.1",
			"c":  "identifier for the client software",
			"v":  "version of the client software",
			"u":  "username",
			"t":  "timestamp of the request",
			"a":  "authentication message: md5(md5(auth_token) + timestamp)",
			"__response": `The server will respond with:
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
		core.AuthGuest,
		a.LoginAS})
}

func (a *AS) LoginAS(ctx core.CommandContext) core.Result {
	var err error
	args := ctx.Args
	handshakeS, err := util.GetArg(args, "hs", true, err)
	handshake, err := util.CastBool(handshakeS, err)
	pVer, err := util.GetArg(args, "p", true, err)
	_, err = util.GetArg(args, "c", true, err)
	_, err = util.GetArg(args, "v", true, err)
	userS, err := util.GetArg(args, "u", true, err)
	timestampS, err := util.GetArg(args, "t", true, err)
	timestamp, err := util.CastInt64(timestampS, err)
	authMsg, err := util.GetArg(args, "a", true, err)

	if err != nil {
		return errorAS(StatusFailed, err)
	}
	if !*handshake {
		return errorAS(StatusFailed, errrs.New("this URL is only for handshake"+
			" requests, hs must be true"))
	}
	if *pVer != "1.2.1" {
		return errorAS(StatusFailed, errrs.New("Protocol version must be 1.2.1"))
	}

	//TODO maybe check audioscrobbler client id and version

	// check timestamp
	timestampReal := time.Now().Unix()
	if util.Abs(timestampReal-*timestamp) > 120 {
		return errorAS(StatusBadTime, nil)
	}

	// get user
	user, err := a.db.GetUserByName(*userS)
	if user == nil {
		logger.Log.Println("incorrect auth from unknown user ", *userS)
		return errorAS(StatusBadAuth, nil)
	}
	// check auth
	// md5(md5(auth_token) + timestamp)
	md5Token := fmt.Sprintf("%x", md5.Sum([]byte(user.AuthToken)))
	correctStr := fmt.Sprintf("%s%d", md5Token, *timestamp)
	correctAuthMsg := fmt.Sprintf("%x", md5.Sum([]byte(correctStr)))
	if correctAuthMsg != *authMsg {
		logger.Log.Println("incorrect auth from", user.Name, "with", *authMsg,
			"instead of", correctAuthMsg)
		return errorAS(StatusBadAuth, nil)
	}

	return core.Result{StatusOK, []interface{}{user.AuthToken, a.nowPlayingURL,
		a.submissionURL}, nil, true}
}

func (a *AS) scrobbleAS(ctx core.CommandContext) core.Result {
	//TODO implement

	return core.ResultByError(core.ErrorNotImplemented)
}

func errorAS(status core.Status, err error) core.Result {
	return core.Result{status, nil, err, true}
}
