package db

type Item struct {
	Id             int64  `db:"item_id"`
	Title          string `db:"title"`
	Artist         string `db:"artist"`
	Genre          string `db:"genre"` //TODO more refined genres
	TrackNumber    uint32 `db:"track_number"`
	musicbrainz_id string `db:"musicbrainz_id"`

	Folder   *Folder `db:"-"`
	FolderId Folder  `db:"folder_id"`

	Album   *Album `db:"-"`
	AlbumId int64  `db:"album_id"`

	Added          int64     `db:"added_date"`
	Rating         complex64 `db:"rating_complex"` //TODO define fancy math around this
	PlayCount      uint32    `db:"play_count"`
	ScrobbledCount uint32    `db:"scrobbled_count"`
	SkipCount      uint32    `db:"skip_count"`
}

type Album struct {
	AlbumId        int64  `db:"album_id"`
	Artist         string `db:"album_artist"`
	Title          string `db:"title"`
	Date           int64  `db:"date"`
	musicbrainz_id string `db:"musicbrainz_id"`
}

type Folder struct {
	Id      int64  `db:"folder_id"`
	AlbumId int64  `db:"album_id"`
	Album   *Album `db:"-"`
	Title   string `db:"title"`
}
