package web

import (
	"config"
	"core"
	"db"
	"encoding/json"
	"errrs"
	"fmt"
	"html/template"
	log "logger"
	"mime"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"util"
)

var allowedTemplates = []string{".html"}

type NetCmdLine struct {
	core       core.Core
	db         *db.DB
	renderInfo WebTemplate
	templates  map[string]*template.Template
}

// context information for web templates
type WebTemplate struct {
	ApiPath string
}

func (n *NetCmdLine) getCmd(r *http.Request) (err error, cmd string, args core.ArgMap) {
	path := r.URL.Path[1:]
	parts := strings.Split(path, "/")
	// ignore parts[2:]
	if err := r.ParseForm(); err != nil {
		return errrs.New("bad request: " + err.Error()), "", nil
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
		ctx.AuthLevel, ctx.UserId, err = n.db.Authenticate(*token)
		if err != nil {
			n.printResult(w, core.ResultByError(err))
			return
		}
	}

	// execute command
	result := n.core.Cmd(ctx)
	if result.Error == core.ErrorCmdNotFound {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	n.printResult(w, result)
}

func (n *NetCmdLine) printResult(w http.ResponseWriter, result core.Result) {
	if result.Error != nil {
		log.Log.Println("api error:", result.Status, result.Error)
	}
	if result.IsRaw {
		fmt.Fprint(w, result.Status)
		if result.Error != nil {
			fmt.Fprintln(w, "", result.Error)
		} else {
			fmt.Fprintln(w)
		}
		if result.Result != nil {
			if array, ok := result.Result.([]interface{}); ok {
				for _, v := range array {
					fmt.Fprintln(w, v)
				}
			} else {
				fmt.Fprintln(w, result.Result)
			}
		}
	} else {
		if bytes, err := json.MarshalIndent(result, "", "  "); err != nil {
			http.Error(w, err.Error(), 500)
		} else {
			w.Write(bytes)
			fmt.Fprintln(w)
		}
	}
}

func (n *NetCmdLine) serveAsset(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = strings.TrimLeft(r.URL.Path, "/")
	if r.URL.Path == "" {
		r.URL.Path = "index.html"
	}

	cache := false
	templateAble := false
	for _, suffix := range allowedTemplates {
		if strings.HasSuffix(r.URL.Path, suffix) {
			cache = config.Current.CacheWebTemplates
			templateAble = true
			break
		}
	}
	if templateAble {
		var tmpl *template.Template
		if cache {
			tmpl = n.templates[r.URL.Path]
		}
		if tmpl == nil {
			data, err := Asset(r.URL.Path)
			if len(data) == 0 || err != nil {
				//TODO handle 404's more nicely
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, err)
				return
			}
			tmpl, err = template.New(r.URL.Path).Parse(string(data))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, err)
				log.Log.Println(err)
				return
			}
			n.templates[r.URL.Path] = tmpl
		}
		n.setMime(w, r)
		err := tmpl.Execute(w, n.renderInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			log.Log.Println(err)
		}
	} else {
		data, err := Asset(r.URL.Path)
		if len(data) == 0 || err != nil {
			//TODO handle 404's more nicely
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, err)
			return
		}
		n.setMime(w, r)
		_, err = w.Write(data)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintln(w, err)
		}
	}
}

func (n *NetCmdLine) setMime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", mime.TypeByExtension(r.URL.Path[strings.LastIndex(r.URL.Path, "."):]))
}

func (n *NetCmdLine) Start(core core.Core, db *db.DB) {
	n.core = core
	n.db = db
	n.templates = make(map[string]*template.Template)
	//TODO insert hostname here
	n.renderInfo = WebTemplate{
		ApiPath: fmt.Sprintf("http://localhost:%d/api/", config.Current.WebPort),
	}

	http.HandleFunc("/api/", n.api)
	http.HandleFunc("/", n.serveAsset)
	url := fmt.Sprint(":", config.Current.WebPort)
	http.ListenAndServe(url, nil)
}
