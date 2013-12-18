package web

import (
	"config"
	"core"
	"fmt"
	"net/http"
	"strings"
)

type NetCmdLine struct {
	core core.Core
}

func (n *NetCmdLine) getCmd(r *http.Request) (cmd string, args []string) {
	path := r.URL.Path[1:]
	parts := strings.Split(path, "/")
	return parts[1], parts[2:]
}

func (n *NetCmdLine) api(w http.ResponseWriter, r *http.Request) {
	cmd, args := n.getCmd(r)

	//fmt.Fprintf(w, "You're at /api/%s!\n", cmd)
	//fmt.Fprintln(w, "args:", args)
	err := n.core.Cmd(cmd, args)
	if err != nil {
		if err == core.ErrorCmdNotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), 500)
		return
	}
}

func (n *NetCmdLine) Start(core core.Core) {
	n.core = core
	http.HandleFunc("/api/", n.api)
	url := fmt.Sprint(":", config.Current.WebPort)
	http.ListenAndServe(url, nil)
}
