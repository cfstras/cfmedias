package web

import (
	"config"
	"core"
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"util"
)

type NetCmdLine struct {
	core core.Core
	db   *db.DB
}

func (n *NetCmdLine) getCmd(r *http.Request) (err error, cmd string, args core.ArgMap) {
	path := r.URL.Path[1:]
	parts := strings.Split(path, "/")
	// ignore parts[2:]
	if err := r.ParseForm(); err != nil {
		return err, "", nil
	}

	return nil, parts[1], core.ArgMap(r.Form)
}

func (n *NetCmdLine) api(w http.ResponseWriter, r *http.Request) {
	// get command name from url
	err, cmd, args := n.getCmd(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var ctx core.CommandContext
	ctx.Cmd = cmd
	ctx.Args = args

	// check auth token
	token, err := util.GetArg(args, "auth_token", false, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx.AuthLevel = core.AuthGuest
	if token != nil {
		ctx.AuthLevel, ctx.UserId, err = n.db.Authenticate([]byte(*token))
		if err != nil {
			if bytes, err := json.MarshalIndent(err, "", "  "); err != nil {
				http.Error(w, err.Error(), 500)
			} else {
				w.Write(bytes)
				fmt.Fprintln(w)
			}
			return
		}
	}

	// execute command
	result := n.core.Cmd(ctx)
	if result.Error == core.ErrorCmdNotFound {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	if bytes, err := json.MarshalIndent(result, "", "  "); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Write(bytes)
		fmt.Fprintln(w)
	}
}

func (n *NetCmdLine) Start(core core.Core, db *db.DB) {
	n.core = core
	n.db = db
	http.HandleFunc("/api/", n.api)
	url := fmt.Sprint(":", config.Current.WebPort)
	http.ListenAndServe(url, nil)
}
