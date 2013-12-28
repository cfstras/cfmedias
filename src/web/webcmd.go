package web

import (
	"config"
	"core"
	"db"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	err, cmd, args := n.getCmd(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Auth
	//TODO get auth token from request
	authToken := []byte{0, 0, 0, 0}
	authLevel, err := n.db.Authenticate(authToken)
	if err != nil {
		if bytes, err := json.MarshalIndent(err, "", "  "); err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Write(bytes)
			fmt.Fprintln(w)
		}
		return
	}

	result := n.core.Cmd(core.CommandContext{cmd, args, authLevel})
	if err == core.ErrorCmdNotFound {
		http.Error(w, "Command not found", 404)
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
