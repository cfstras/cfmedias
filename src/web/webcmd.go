package web

import (
	"config"
	"core"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type NetCmdLine struct {
	core core.Core
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

	if result := n.core.Cmd(cmd, args, core.AuthGuest); err != nil {
		if err == core.ErrorCmdNotFound {
			http.Error(w, "Command not found", 404)
			return
		}
		http.Error(w, err.Error(), 500)
	} else {
		if str, err := json.MarshalIndent(result, "", "  "); err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			fmt.Fprint(w, str)
		}
	}
}

func (n *NetCmdLine) Start(core core.Core) {
	n.core = core
	http.HandleFunc("/api/", n.api)
	url := fmt.Sprint(":", config.Current.WebPort)
	http.ListenAndServe(url, nil)
}
