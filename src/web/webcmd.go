package web

import (
	"bytes"
	"config"
	"core"
	"db"
	"errrs"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"html/template"
	"mime"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"util"
)

var allowedTemplates = []string{".html"}

type NetCmdLine struct {
	core       core.Core
	db         *db.DB
	martini    *martini.ClassicMartini
	renderInfo WebTemplate
	templates  map[string]*template.Template
}

// context information for web templates
type WebTemplate struct {
	ApiPath string
}

func (n *NetCmdLine) api(params martini.Params, r render.Render, req *http.Request) (
	int, string) {
	if err := req.ParseForm(); err != nil {
		return http.StatusBadRequest, "bad request: " + err.Error()
	}

	var ctx core.CommandContext
	ctx.Cmd = params["cmd"]
	ctx.Args = core.ArgMap(req.Form)

	// check auth token
	token, err := util.GetArg(ctx.Args, "auth_token", false, nil)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}
	ctx.AuthLevel = core.AuthGuest
	if token != nil {
		ctx.AuthLevel, ctx.UserId, err = n.db.Authenticate(*token)
		if err != nil {
			return http.StatusUnauthorized, err.Error()
		}
	}

	// execute command
	result := n.core.Cmd(ctx)
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

func (n *NetCmdLine) Start(core core.Core, db *db.DB) {
	n.core = core
	n.db = db
	n.templates = make(map[string]*template.Template)
	//TODO insert hostname here
	n.renderInfo = WebTemplate{
		ApiPath: fmt.Sprintf("http://localhost:%d/api/", config.Current.WebPort),
	}

	m := martini.Classic()
	n.martini = m
	m.Use(render.Renderer())

	m.Group("/api", func(r martini.Router) {
		r.Get("/:cmd", n.api)
		r.Post("/:cmd", n.api)
	})

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
