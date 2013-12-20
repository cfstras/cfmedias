package db

import (
	"config"
	"core"
	"errrs"
	"math"
	"strconv"
)

const (
	ParamString uint = iota
	ParamFloat
	ParamBool
)

func (db *DB) initLogic(c core.Core) {
	c.RegisterCommand(core.Command{
		[]string{"trackplayed"},
		"Inserts a track play into the database",
		core.AuthUser,
		db.TrackPlayed})
}

func (db *DB) TrackPlayed(ctx core.CommandContext) core.Result {
	args := ctx.Args

	tracks, err := db.GetItem(args)

	lengthS, err := getArg(args, "length", true, err)
	length_playedS, err := getArg(args, "length_played", true, err)
	scrobbledS, err := getArg(args, "scrobbled", true, err)

	length, err := castFloat(lengthS, err)
	length_played, err := castFloat(length_playedS, err)
	scrobbled, err := castBool(scrobbledS, err)

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

// Fetches an argument from an ArgMap, used for single args
// Breaks and passes along the error given, if it is not nil.
// If the argument does not exist, the return value is nil.
// Parameter force can be used to return an error if the argument does not exist.
func getArg(args core.ArgMap, arg string, force bool, err error) (*string, error) {
	if err != nil {
		return nil, err
	}
	value, ok := args[arg]
	if !ok || len(value) == 0 {
		if force {
			return nil, errrs.New("argument " + arg + " missing!")
		}
		return nil, nil
	}
	if len(value) > 1 {
		return nil, errrs.New("argument " + arg + " cannot be supplied more than once!")
	}
	return &value[0], nil
}

// Converts a *string to a boolean.
// Passes along errrs, if not nil.
func castBool(arg *string, err error) (*bool, error) {
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, nil
	}
	casted, err := strconv.ParseBool(*arg)
	if err != nil {
		return nil, err
	}
	return &casted, nil
}

// Converts a *string to a float32.
// Passes along errrs, if not nil.
func castFloat(arg *string, err error) (*float32, error) {
	if err != nil {
		return nil, err
	}
	if arg == nil {
		return nil, nil
	}
	casted, err := strconv.ParseFloat(*arg, 32)
	if err != nil {
		return nil, err
	}
	smaller := float32(casted)
	return &smaller, nil
}
