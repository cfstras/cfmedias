package web

import (
	"bytes"
	"config"
	"core"
	"db"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"mime"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"util"
)

var allowedTemplates = []string{".html"}

type NetCmdLine struct {
	core    core.Core
	db      *db.DB
	martini *martini.ClassicMartini
}

func (n *NetCmdLine) api(r render.Render, ctx *core.CommandContext) (int, string) {
	// execute command
	result := n.core.Cmd(*ctx)
	if result.Error == core.ErrorCmdNotFound {
		return http.StatusNotFound, result.Error.Error()
	}
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error.Error()
	}
	if result.IsRaw {
		str := &bytes.Buffer{}
		if result.Result != nil {
			if array, ok := result.Result.([]interface{}); ok {
				for _, v := range array {
					fmt.Fprintln(str, v)
				}
			} else {
				fmt.Fprintln(str, result.Result)
			}
		}
		return 200, str.String()
	} else {
		r.JSON(200, result.Result)
		return 200, ""
	}
}

func (n *NetCmdLine) Start(coreInstance core.Core, db *db.DB) {
	n.core = coreInstance
	n.db = db

	m := martini.Classic()
	n.martini = m
	m.Use(render.Renderer())

	m.Group("/api", func(r martini.Router) {
		r.Get("/:cmd", n.api)
		r.Post("/:cmd", n.api)
	}, mapArgs, n.authenticate)

	m.Use(func(c martini.Context, r *http.Request) {
		c.Map(r.URL.Path[1:])
	})
	m.Get("/", index_html)
	m.Get("/**", setMime, Asset)

	os.Setenv("PORT", fmt.Sprint(config.Current.WebPort))
	m.Run()
}

func setMime(w http.ResponseWriter, path string) {
	w.Header().Set("Content-Type", mime.TypeByExtension(path[strings.LastIndex(path, "."):]))
}

func mapArgs(c martini.Context, params martini.Params, r *http.Request,
	w http.ResponseWriter) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad request:", err)
		return
	}

	var ctx core.CommandContext
	ctx.Cmd = params["cmd"]
	ctx.Args = core.ArgMap(r.Form)
	c.Map(&ctx)
}

func (n *NetCmdLine) authenticate(c martini.Context, ctx *core.CommandContext,
	w http.ResponseWriter) {
	// check auth token
	token, err := util.GetArg(ctx.Args, "auth_token", false, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad request:", err)
		return
	}
	ctx.AuthLevel = core.AuthGuest
	if token != nil {
		ctx.AuthLevel, ctx.UserId, err = n.db.Authenticate(*token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "bad request:", err)
			return
		}
	}
}
