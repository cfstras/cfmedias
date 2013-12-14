package db

type Item struct {
	Id            uint64 `db:"item_id"`
	Title         string `db:"title"`
	Artist        string `db:"artist"`
	Genre         string `db:"genre"` //TODO more refined genres
	TrackNumber   uint32 `db:"track_number"`
	MusicbrainzId string `db:"musicbrainz_id"`

	Folder   *Folder `db:"-"`
	FolderId uint64  `db:"folder_id"`

	Album   *Album `db:"-"`
	AlbumId uint64 `db:"album_id"`

	Added          uint64  `db:"added_date"`
	Rating         float32 `db:"rating"` //TODO define fancy math around this
	PlayCount      uint32  `db:"play_count"`
	ScrobbledCount uint32  `db:"scrobbled_count"`
	SkipCount      uint32  `db:"skip_count"`
}

type Album struct {
	Id            uint64 `db:"album_id"`
	Artist        string `db:"album_artist"`
	Title         string `db:"title"`
	Date          uint64 `db:"date"`
	MusicbrainzId string `db:"musicbrainz_id"`
}

type Folder struct {
	Id      uint64 `db:"folder_id"`
	AlbumId uint64 `db:"album_id"`
	Album   *Album `db:"-"`
	Title   string `db:"title"`
}
