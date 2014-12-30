package ipod

import (
	"encoding/json"
	"fmt"
	"github.com/cfstras/cfmedias/config"
	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/db"
	"github.com/cfstras/cfmedias/errrs"
	"github.com/cfstras/cfmedias/ipod/gpod"
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

func (p *IPod) Sync(mountpoint string) error {
	logger.Log.Println("getting data...")
	tracks, err := p.db.ListAll()
	if err != nil {
		return err
	}
	logger.Log.Println(len(tracks), "tracks found.")
	logger.Log.Println("indexing target...")

	ipodDb, err := gpod.New(mountpoint)
	if err != nil {
		return err
	}
	ipodTracks := ipodDb.Tracks()
	logger.Log.Println(len(ipodTracks), "tracks on iPod.")

	idxFuncGpod := func(t gpod.Track) string {
		return t.Title() + Seperator + t.Album() + Seperator + t.Artist()
	}
	idxFunc := func(t db.Item) string {
		str := t.Title + Seperator
		if t.Album.Valid {
			str += t.Album.String
		}
		return str + Seperator + t.Artist
	}
	idx := make(map[string]gpod.Track)
	for _, t := range ipodTracks {
		idx[idxFuncGpod(t)] = t
	}
	var tracksMissing []*db.Item
	// cfmedias db id -> gpod track
	tracksFound := make(map[int64]gpod.Track)
	var tracksUnmatched []gpod.Track
	matched := make(map[string]bool)
	for i := range tracks {
		t := &tracks[i]
		match, ok := idx[idxFunc(*t)]
		if !ok {
			tracksMissing = append(tracksMissing, t)
			continue
		}
		matched[idxFunc(*t)] = true
		tracksFound[t.Id] = match
	}
	for i := range ipodTracks {
		t := &ipodTracks[i]
		if !matched[idxFuncGpod(*t)] {
			tracksUnmatched = append(tracksUnmatched, *t)
		}
	}
	logger.Log.Println(len(tracksFound), "tracks found,",
		len(tracksUnmatched), "unknown tracks on iPod,",
		len(tracksMissing), "tracks to copy")
	m := map[string]interface{}{
		"unmatched": tracksUnmatched,
		"missing":   tracksMissing,
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("sync-info.json", b, 0644)
	if err != nil {
		return err
	}
	logger.Log.Println("Wrote sync-info.json.")

	//TODO update tags

	// delete unmatched
	deleteErrors := make(map[gpod.Track]error)
	for _, t := range tracksUnmatched {
		logger.Log.Println("Deleting", t)
		err := t.Delete()
		if err != nil {
			logger.Log.Println("Error deleting", t, err)
			deleteErrors[t] = err
			continue
		}
	}

	// add missing
	copyErrors := make(map[*db.Item]error)
	for _, t := range tracksMissing {
		logger.Log.Println("Adding", t)
		ipodT := gpod.NewTrack(ipodDb)
		ipodT.SetTitle(t.Title)
		ipodT.SetAlbum(db.StrStr(t.Album))
		ipodT.SetArtist(t.Artist)
		ipodT.SetGenre(db.StrStr(t.Genre))
		//ipodT.SetFiletype(t.Filename)
		//ipodT.SetComment(t.Comment)
		//ipodT.SetComposer(t.Composer)
		//ipodT.SetDescription(t.Description)
		ipodT.SetAlbumartist(db.StrStr(t.AlbumArtist))
		//ipodT.SetSize(t.Size)
		//ipodT.SetLength(t.Length)

		//TODO add missing fields
		// don't make sense yet
		//ipodT.SetRating(t.Rating())
		//ipodT.SetPlaycount(t.Playcount)
		ipodT.SetTimeAdded(t.CreatedAt)

		//TODO playlist hierarchy

		path := t.Path()
		if path == nil {
			err := errrs.New("No file to copy")
			logger.Log.Println("Error finding", t, err)
			copyErrors[t] = err
			continue
		}
		ipodDb.Add(ipodT)
		ipodDb.MPL().Add(ipodT)
		err := ipodDb.Copy(ipodT, *path)
		if err != nil {
			logger.Log.Println("Error copying", t, err)
			copyErrors[t] = err
			ipodDb.MPL().Remove(ipodT)
			ipodDb.Remove(ipodT)
		}
	}
	err = ipodDb.Save()
	if err != nil {
		return err
	}
	// give back other errors, if any
	if len(copyErrors) > 0 {
		serr := &syncError{deleteErrors, copyErrors}
		return serr
	}
	return nil
}

type syncError struct {
	DeleteErrors map[gpod.Track]error
	CopyErrors   map[*db.Item]error
}

func (e *syncError) Error() string {
	return fmt.Sprint(len(e.DeleteErrors), "delete errors,",
		len(e.CopyErrors), "copy errors")
}
