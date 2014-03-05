package db

import (
	"config"
	"core"
	"math"
	"util"
)

const (
	ParamString uint = iota
	ParamFloat
	ParamBool
)

func (db *DB) initLogic(c core.Core) {
	c.RegisterCommand(core.Command{
		[]string{"trackplayed"},
		"Inserts a track playback into the database statistics",
		map[string]string{
			"title":         "Title",
			"artist":        "Artist",
			"album":         "Album",
			"album_artist":  "Album Artist",
			"length":        "Track length in ms",
			"date":          "Time the listening occurred, as Unix timestamp",
			"length_played": "The time the track was listened to (when fully played: length) in ms",
			"scrobbled":     "Whether the track was scrobbled to last.fm/libre.fm"},
		core.AuthUser,
		db.TrackPlayed})
}

func (db *DB) TrackPlayed(ctx core.CommandContext) core.Result {
	args := ctx.Args

	tracks, err := db.GetItem(args)

	lengthS, err := util.GetArg(args, "length", true, err)
	length_playedS, err := util.GetArg(args, "length_played", true, err)
	scrobbledS, err := util.GetArg(args, "scrobbled", true, err)

	length, err := util.CastFloat(lengthS, err)
	length_played, err := util.CastFloat(length_playedS, err)
	scrobbled, err := util.CastBool(scrobbledS, err)

	if err != nil {
		return core.Result{Status: core.StatusError, Error: err}
	}

	if len(tracks) == 0 {
		//TODO insert track
		return core.Result{Status: core.StatusItemNotFound}
	}
	if len(tracks) > 1 {
		return core.Result{Status: core.StatusQueryNotUnique, Results: ItemToInterfaceSlice(tracks)}
	}

	track := tracks[0]

	// update stats
	x := float64(*length_played / *length)
	tu := float64(config.Current.ListenedUpperThreshold)
	tl := float64(config.Current.ListenedLowerThreshold)

	scoreAdd := float32(math.Min(1, math.Max(0, (x-tl)/(tu-tl))))
	//TODO test this

	track.PlayCount++
	if scoreAdd > 0 {
		track.PlayScore += scoreAdd
		track.ScoredCount++
	}
	if *scrobbled {
		track.ScrobbleCount++
	}

	rows, err := db.dbmap.Update(track)
	if err != nil {
		return core.Result{Status: core.StatusError, Error: err}
	}
	if rows == 0 {
		return core.Result{Status: core.StatusItemNotFound}
	}

	return core.Result{Status: core.StatusOK, Results: ItemToInterfaceSlice(tracks)}
}
