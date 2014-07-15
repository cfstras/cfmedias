package db

import (
	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/util"
)

func (db *DB) GetItem(args core.ArgMap) ([]*Item, error) {
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

	qArgs := []interface{}{title, artist} // never nil, because force is true for them

	q := `select * from ` + ItemTable + `
		where title = ? and artist = ? `
	if album != nil {
		q += `and album = ? `
		qArgs = append(qArgs, album)
	}
	if album_artist != nil && false { //TODO album-artists are not implemented yet
		q += `and album_artist = ? `
		qArgs = append(qArgs, album_artist)
	}
	//TODO mbid is not implemented. If given, only search mbid.
	/*if musicbrainz_id != nil {
		q += `and musicbrainz_id = ? `
		qArgs = append(qArgs, musicbrainz_id)
	}*/

	// get track info
	//TODO get DB write lock!
	tracks, err := db.dbmap.Select(Item{}, q, qArgs...)
	if err != nil {
		return nil, err
	}

	return InterfaceToItemSlice(tracks), nil
}

func InterfaceToItemSlice(slice []interface{}) []*Item {
	items := make([]*Item, len(slice))
	for i, v := range slice {
		items[i] = v.(*Item)
	}
	return items
}

func ItemToInterfaceSlice(slice []*Item) []interface{} {
	items := make([]interface{}, len(slice))
	for i, v := range slice {
		items[i] = v
	}
	return items
}
