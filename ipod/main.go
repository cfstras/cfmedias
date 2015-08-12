package ipod

import (
	"encoding/json"
	"fmt"
	"github.com/cfstras/cfmedias/config"
	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/db"
	"github.com/cfstras/cfmedias/errrs"
	"github.com/cfstras/cfmedias/logger"
	"github.com/cfstras/cfmedias/sync"
	"github.com/cfstras/cfmedias/util"
	"io/ioutil"
)

type IPod struct {
	core core.Core
	db   *db.DB
	sync *sync.Sync

	config *Config
}

type Config struct {
	ConvertRules map[string]string
}

var defaultConfig Config = Config{
	ConvertRules: sync.DefaultConfig.ConvertRules,
}

const (
	Seperator  = "\x00"
	PluginName = "ipod"
)

func init() {
	configCopy := defaultConfig
	config.RegisterPlugin(PluginName, &configCopy, &Config{})
}

func (p *IPod) Start(c core.Core, db *db.DB, s *sync.Sync) {
	p.core = c
	p.db = db
	p.sync = s
	p.config = config.Current.Plugins[PluginName].(*Config)

	c.RegisterCommand(core.Command{[]string{"ipod"},
		"Syncs media with an iPod device. By default, Lossles files are converted to MP3 V0.",
		map[string]string{"mountpoint": "mountpoint of the iPod"},
		core.AuthAdmin,
		func(ctx core.CommandContext) core.Result {
			path, err := util.GetArg(ctx.Args, "mountpoint", true, nil)
			if err != nil {
				return core.ResultByError(err)
			}
			return core.ResultByError(p.Sync(*path))
		}})
}


