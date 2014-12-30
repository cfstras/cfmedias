package sync

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cfstras/cfmedias/config"
	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/db"
	"github.com/cfstras/cfmedias/errrs"
	"github.com/cfstras/cfmedias/logger"
	"github.com/cfstras/cfmedias/util"
)

type Sync struct {
	core   core.Core
	db     *db.DB
	config *Config
}

type Config struct {
	ConvertRules map[string]string
	Formats      map[string]Format
}

type Format struct {
	Name string

	Extension       string // commonly used container extension
	Encoder         string // encoder name in ffmpeg
	Quality         int    // quality or bitrate
	ConstantBitrate bool
}

var (
	FormatFLAC Format = Format{"flac", "flac", "flac", 8, false}
	FormatALAC        = Format{"alac", "m4a", "alac", -1, false}
	FormatWave        = Format{"wav", "wav", "pcm_s16le", -1, true}

	FormatV0       = Format{"mp3v0", "mp3", "mp3", 0, false}
	FormatV1       = Format{"mp3v1", "mp3", "mp3", 1, false}
	FormatV2       = Format{"mp3v2", "mp3", "mp3", 2, false}
	FormatMP3_320K = Format{"mp3_320k", "mp3", "mp3", 320, true}
	FormatMP3_256K = Format{"mp3_256k", "mp3", "mp3", 256, true}
	FormatMP3_192K = Format{"mp3_192k", "mp3", "mp3", 192, true}

	FormatAAC_320K = Format{"aac_320k", "m4a", "aac", 320, true}

	FormatList = []Format{FormatFLAC, FormatALAC,
		FormatV0, FormatV1, FormatV2,
		FormatMP3_192K, FormatMP3_256K, FormatAAC_320K,
		FormatAAC_320K,
	}

	DefaultConfig Config = Config{
		ConvertRules: map[string]string{
			"flac": "mp3v0",
			"alac": "mp3v0",
			"wav":  "mp3v0",
		},
	}
)

const PluginName = "sync"

func init() {
	// map formats
	DefaultConfig.Formats = make(map[string]Format)
	for _, f := range FormatList {
		DefaultConfig.Formats[f.Name] = f
	}

	configCopy := DefaultConfig
	config.RegisterPlugin(PluginName, &configCopy, &Config{})
}

func (s *Sync) Start(c core.Core, db *db.DB) {
	s.core = c
	s.db = db
	s.config = config.Current.Plugins[PluginName].(*Config)

	c.RegisterCommand(core.Command{[]string{"sync", "s"},
		"Syncs media with a device or folder. By default, lossles files are converted to MP3 V0.",
		map[string]string{
			"path":    "Target path",
			"convert": "boolean, default is true"},
		core.AuthAdmin,
		func(ctx core.CommandContext) core.Result {
			args := ctx.Args
			var err error
			pathS, err := util.GetArg(args, "path", true, err)
			convertS, err := util.GetArg(args, "convert", true, err)
			doConvert, err := util.CastBool(convertS, err)
			if err != nil {
				return core.ResultByError(err)
			}
			convert := true
			if doConvert != nil {
				convert = *doConvert
			}
			return core.ResultByError(s.Sync(*pathS, convert))
		}})
}

type file struct {
	// relative to targetPath
	path string
	info os.FileInfo
}

func (s *Sync) Sync(targetPath string, convert bool) error {
	logger.Log.Println("getting data...")
	tracks, err := s.db.ListAll()
	if err != nil {
		return err
	}
	logger.Log.Println(len(tracks), "tracks found.")
	logger.Log.Println("indexing target...")

	cleanPath := path.Clean(targetPath)
	if cleanPath != "" {
		cleanPath += "/"
	}
	var targetFiles []file
	filepath.Walk(targetPath, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		filepath = path.Clean(filepath)
		if !strings.HasPrefix(filepath, cleanPath) {
			return errrs.New("file " + filepath + " is not within " + cleanPath)
		}
		filepath = strings.TrimPrefix(filepath, cleanPath)
		targetFiles = append(targetFiles, file{filepath, info})
		return nil
	})
	logger.Log.Println(len(targetFiles), "files in target")

	return core.ErrorNotImplemented
}
