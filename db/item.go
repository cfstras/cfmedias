package db

import (
	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/util"
)

func (db *DB) GetItem(args core.ArgMap) ([]Item, error) {
	// check if necessary args are there
	var err error
	title, err := util.GetArg(args, "title", true, err)
	artist, err := util.GetArg(args, "artist", true, err)
	album, err := util.GetArg(args, "album", false, err)
	album_artist, err := util.GetArg(args, "album_artist", false, err)
	//musicbrainz_id, err := util.GetArg(args, "musicbrainz_id", false, err)

	if err != nil {
		return nil, err
	}

	limit := map[string]interface{}{"title": title, "artist": artist}

	if album != nil {
		limit["album"] = album
	}
	if album_artist != nil && false { //TODO album-artists are not implemented yet
		limit["album_artist"] = album_artist
	}
	//TODO mbid is not implemented. If given, only search mbid.
	/*if musicbrainz_id != nil {
		limit["musicbrainz_id"] = musicbrainz_id
	}*/

	tracks := make([]Item, 0)
	err = db.db.Where(limit).Find(tracks).Error
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func InterfaceToItemSlice(slice []interface{}) []Item {
	items := make([]Item, len(slice))
	for i, v := range slice {
		items[i] = v.(Item)
	}
	return items
}

func ItemToInterfaceSlice(slice []Item) []interface{} {
	items := make([]interface{}, len(slice))
	for i, v := range slice {
		items[i] = v
	}
	return items
}
