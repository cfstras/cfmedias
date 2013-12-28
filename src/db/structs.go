package db

import (
	"core"
	"errrs"
	"fmt"
	"github.com/coopernurse/gorp"
	"math"
	"time"
)

const (
	ItemTable   = "items"
	FolderTable = "folders"
	UserTable   = "users"
)

// when inserting an item, Folder has to be a pointer with only the Path set.
type Item struct {
	Id            uint64  `db:"item_id"`
	Title         string  `db:"title"`
	Artist        string  `db:"artist"`
	AlbumArtist   *string `db:"album_artist"`
	Album         *string `db:"album"`
	Genre         *string `db:"genre"` //TODO more refined genres
	TrackNumber   uint32  `db:"track_number"`
	Filename      *string `db:"filename"`
	MusicbrainzId *string `db:"musicbrainz_id"`

	Folder   *Folder `db:"-"`
	FolderId *uint64 `db:"folder_id"`

	Added int64 `db:"added"`

	// Total play score of a track.
	// Gets incremented by up to 1 each time a user listens to a track,
	// depending on how far in he listened to the track.
	// The formula for the added value is:
	//  x: amount of track listened to (0..1)
	//  tl: listened lower threshold
	//  tu: listened upper threshold
	//  add(x) = min( 1, max( 0, (x - tl) / (tu - tl) ) )
	PlayScore float32 `db:"play_score"`

	// Total times a track was played (including skips)
	// Track Score := (PlayScore / PlayCount)
	PlayCount uint32 `db:"play_count"`

	// Total times a track was played and scored as positive (x >= tl)
	ScoredCount uint32 `db:"scored_count"`

	// Number of plays that are registered at the local scrobbler.
	//
	// When a track is listened and scored positively, this is incremented
	// and a scrobble is sent to the service.
	ScrobbleCount uint32 `db:"scrobbled_count"`
}

type Folder struct {
	Id   uint64 `db:"folder_id"`
	Path string `db:"path"`

	Added         int64  `db:"added"`
	MusicbrainzId string `db:"musicbrainz_folder_id"`
}

type User struct {
	Id uint64 `db:"user_id" json:"-"`

	// Login name for the user. Has to be unique.
	Name string `db:"name" json:"name"`

	// E-mail address. Has to be unique.
	EMail string `db:"email" json:"email"`

	AuthLevel core.AuthLevel `db:"auth_level" json:"auth_level"`

	// salt + pbkdf2-hashed password
	Password []byte `db:"password" json:"-"`

	// Authentication token, used for sending requests
	AuthToken []byte `db:"auth_token" json:"auth_token"`

	//TODO optionally add a common secret for authenticating messages via HMAC
}

type ItemPathView struct {
	Id       uint64
	Filename string
	Path     string
}

func (i *Item) String() string {
	return fmt.Sprintf("Item[%d]{%s / %s - %s / %d %s, %s]", i.Id, i.Artist,
		str(i.AlbumArtist), str(i.Album), i.TrackNumber, i.Title, str(i.Genre))
}
func str(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

// Returns a rating of the song and an indication of how accurate the score
// might be.
// Both values are in the range of 0 to 1.
// The computation uses a combination of PlayScore, PlayCount and ScrobbledCount.
func (i *Item) Rating() (rating, accuracy float32) {
	// accuracy is based on how much data there is on this track.
	//TODO use an average data count on the whole database instead of hard one
	acc := (i.PlayScore - 5) // every data below 5 listens is basically useless
	score := i.PlayScore / float32(i.PlayCount)

	//incomingScrobbles := i.ScrobbledCount - i.ScoredCount

	//TODO average score over ratio from incoming scrobbles to own data
	// also use the median of all scrobbles to weigh out the score

	return limit(score, 0, 1), limit(acc, 0, 1)
}

func (i *Item) Skipped() uint32 {
	return i.ScoredCount - i.PlayCount
}

func (i *Item) PreInsert(s gorp.SqlExecutor) error {
	i.Added = time.Now().Unix()

	if i.Filename == nil {
		// No filename, no path.
		i.Folder = nil
		i.FolderId = nil
		return nil
	}
	if i.Folder == nil {
		return errrs.New("Item insert with path needs a Folder!")
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
	i.FolderId = &i.Folder.Id
	return nil
}

func (f *Folder) PreInsert(s gorp.SqlExecutor) error {
	f.Added = time.Now().Unix()
	return nil
}

func limit(val, min, max float32) float32 {
	return float32(math.Min(float64(min), math.Max(float64(max), float64(val))))
}
