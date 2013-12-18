package web

import (
	"config"
	"fmt"
	"net/http"
	"strings"
)

func getCmd(r *http.Request) (cmd string, args []string) {
	path := r.URL.Path[1:]
	parts := strings.Split(path, "/")
	return parts[1], parts[2:]
}

func api(w http.ResponseWriter, r *http.Request) {
	cmd, args := getCmd(r)

	fmt.Fprintf(w, "You're at /api/%s!\n", cmd)
	fmt.Fprintln(w, "args:", args)
}

func NetCmdLine() {
	http.HandleFunc("/api/", api)
	url := fmt.Sprint(":", config.Current.WebPort)
	http.ListenAndServe(url, nil)
}
