package db

import (
	"errors"
	"github.com/coopernurse/gorp"
	"time"
)

const (
	ItemTable   = "items"
	FolderTable = "folders"
)

// when inserting an item, Folder has to be a pointer with only Path set.
type Item struct {
	Id            uint64 `db:"item_id"`
	Title         string `db:"title"`
	Artist        string `db:"artist"`
	AlbumArtist   string `db:"album_artist"`
	Album         string `db:"album"`
	Genre         string `db:"genre"` //TODO more refined genres
	TrackNumber   uint32 `db:"track_number"`
	Filename      string `db:"filename"`
	MusicbrainzId string `db:"musicbrainz_id"`

	Folder   *Folder `db:"-"`
	FolderId uint64  `db:"folder_id"`

	Added          int64   `db:"added_date"`
	Rating         float32 `db:"rating"` //TODO define fancy math around this
	PlayCount      uint32  `db:"play_count"`
	ScrobbledCount uint32  `db:"scrobbled_count"`
	SkipCount      uint32  `db:"skip_count"`
}

type Folder struct {
	Id   uint64 `db:"folder_id"`
	Path string `db:"path"`

	Added         int64  `db:"added_date_folder"`
	MusicbrainzId string `db:"musicbrainz_folder_id"`
}

type ItemPathView struct {
	Id       uint64
	Filename string
	Path     string
}

func (i *Item) PreInsert(s gorp.SqlExecutor) error {
	i.Added = time.Now().Unix()

	if i.Folder == nil {
		return errors.New("Item insert needs a Folder!")
	}

	// set the folder foreign key
	oldFolder := Folder{}
	if err := s.SelectOne(&oldFolder,
		`select * from `+FolderTable+` where path = ?`,
		i.Folder.Path); err != nil {

		return err
	} else if oldFolder.Id == 0 { // not yet there, insert

		if err := s.Insert(i.Folder); err != nil {
			return err
		}
	} else { // all is well, copy folder
		i.Folder = &oldFolder
	}
	i.FolderId = i.Folder.Id
	return nil
}

func (f *Folder) PreInsert(s gorp.SqlExecutor) error {
	f.Added = time.Now().Unix()
	return nil
}
