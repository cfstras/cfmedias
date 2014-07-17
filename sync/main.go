package sync

import (
	"github.com/cfstras/cfmedias/config"
	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/db"
)

type Sync struct {
	core core.Core
	db   *db.DB
}

type Config struct {
	ConvertRules map[string]Format
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

	Formats []Format = []Format{FormatFLAC, FormatALAC,
		FormatV0, FormatV1, FormatV2,
		FormatMP3_192K, FormatMP3_256K, FormatAAC_320K,
		FormatAAC_320K,
	}

	defaultConfig Config = Config{
		ConvertRules: map[string]Format{
			"flac": FormatV0,
			"alac": FormatV0,
			"wav":  FormatV0,
		},
	}
)

func init() {
	config.RegisterPlugin("sync", defaultConfig)
}

func (s *Sync) Start(c core.Core, db *db.DB) {
	s.core = c
	s.db = db
}
