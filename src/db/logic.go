package db

import (
	"config"
	"core"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
)

const (
	ParamString uint = iota
	ParamFloat
	ParamBool
)

func (db *DB) TrackPlayed(args core.ArgMap, w io.Writer) error {
	// check if necessary args are there
	var err error
	title, err := getArg(args, "title", true, err)
	artist, err := getArg(args, "artist", true, err)
	album, err := getArg(args, "album", false, err)
	album_artist, err := getArg(args, "album_artist", false, err)
	musicbrainz_id, err := getArg(args, "musicbrainz_id", false, err)

	lengthS, err := getArg(args, "length", true, err)
	length_playedS, err := getArg(args, "length_played", true, err)
	scrobbledS, err := getArg(args, "scrobbled", true, err)

	length, err := castFloat(lengthS, err)
	length_played, err := castFloat(length_playedS, err)
	scrobbled, err := castBool(scrobbledS, err)

	if err != nil {
		return err
	}

	qArgs := []*string{title, artist} // never nil, because force is true for them

	q := `select * from ` + ItemTable + `
		where title = ? and artist = ? `
	if album != nil {
		q += `and album = ? `
		qArgs = append(qArgs, album)
	}
	if album_artist != nil {
		q += `and album_artist = ? `
		qArgs = append(qArgs, album_artist)
	}
	if musicbrainz_id != nil {
		q += `and musicbrainz_id = ? `
		qArgs = append(qArgs, musicbrainz_id)
	}

	// get track info
	//TODO get DB write lock!
	tracks, err := db.dbmap.Select(Item{}, q, qArgs)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		//TODO insert track
		return core.ErrorItemNotFound
	}
	if len(tracks) >= 0 {
		fmt.Fprintln(w, "Multiple tracks found! Please re-try with more "+
			"accurate arguments.")
		for _, t := range tracks {
			fmt.Fprintln(w, t)
		}
		return core.ErrorQueryAmbiguous
	}

	if tracks[0].(type) != *Item {
		return errors.New("a;lsjfa;lskjfa;lskjf;alskjf")
	}
	track := *Item(tracks[0])

	//DEBUG
	fmt.Println("track found:", track)

	// update stats
	x := float64(*length_played / *length)
	tu := float64(config.Current.ListenedUpperThreshold)
	tl := float64(config.Current.ListenedLowerThreshold)

	scoreAdd := float32(math.Min(1, math.Max(0, (x-tl)/(tu-tl))))
	//TODO debug

	track.PlayCount++
	if scoreAdd > 0 {
		track.PlayScore += scoreAdd
		track.ScoredCount++
	}
	if scrobbled {
		track.ScrobbleCount++
	}

	err := db.dbmap.Update(track)
	if err != nil {
		return err
	}

	//TODO
	return errors.New("Not implemented")
}

// for single args
func getArg(args core.ArgMap, arg string, force bool, err error) (*string, error) {
	if err != nil {
		return nil, err
	}
	value, ok := args["title"]
	if !ok || len(value) == 0 {
		if force {
			return nil, errors.New("argument " + arg + " missing!")
		}
		return nil, nil
	}
	if len(value) > 1 {
		return nil, errors.New("argument " + arg + " cannot be supplied more than once!")
	}
	return &value[0], nil
}

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
