package ipod

import (
	"github.com/cfstras/cfmedias/config"
	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/db"
	"github.com/cfstras/cfmedias/errrs"
	"github.com/cfstras/cfmedias/ipod/glib"
	"github.com/cfstras/cfmedias/logger"
	"github.com/cfstras/cfmedias/sync"
	"github.com/cfstras/cfmedias/util"
)

/*
#cgo pkg-config: libgpod-1.0
#include "gpod/itdb.h"
#include "stdlib.h"
*/
import "C"

type IPod struct {
	core core.Core
	db   *db.DB
	sync *sync.Sync
}

type Config struct {
	sync.Config
}

var defaultConfig Config = Config{
	sync.Config{
		ConvertRules: map[string]sync.Format{
			"flac": sync.FormatV0,
			"alac": sync.FormatV0,
			"wav":  sync.FormatV0,
		},
	},
}

func init() {
	config.RegisterPlugin("ipod", defaultConfig)
}

func (p *IPod) Start(c core.Core, db *db.DB, s *sync.Sync) {
	p.core = c
	p.db = db
	p.sync = s

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

func (p *IPod) Sync(mountpoint string) error {
	logger.Log.Println("getting data...")
	tracks, err := p.db.ListAll()
	if err != nil {
		return err
	}
	logger.Log.Println(len(tracks), "tracks found.")
	logger.Log.Println("indexing target...")

	var gerr *C.GError
	mntpoint := glib.CStr(mountpoint)
	itdb := C.itdb_parse((*C.gchar)(mntpoint), &gerr)
	glib.Free(mntpoint)
	if itdb == nil {
		str := C.GoString((*C.char)(gerr.message))
		return errrs.New(str)
	}
	ptr := itdb.tracks
	for ptr != nil {
		track := (*C.Itdb_Track)(ptr.data)
		title := C.GoString((*C.char)(track.title))
		logger.Log.Println(title)
		ptr = ptr.next
	}

	return core.ErrorNotImplemented
}
