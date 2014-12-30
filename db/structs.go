package db

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/cfstras/cfmedias/core"

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
	Id            int64
	Title         string         `sql:"size:255"`
	Artist        string         `sql:"size:255"`
	AlbumArtist   sql.NullString `sql:"size:255"`
	Album         sql.NullString `sql:"size:255" json:",string"`
	Genre         sql.NullString `sql:"size:255"` //TODO more refined genres
	TrackNumber   uint32
	Filename      sql.NullString `sql:"size:255"`
	MusicbrainzId sql.NullString `sql:"size:36"`

	Folder   Folder
	FolderId sql.NullInt64

	CreatedAt time.Time

	// Total play score of a track.
	// Gets incremented by up to 1 each time a user listens to a track,
	// depending on how far in he listened to the track.
	// The formula for the added value is:
	//  x: amount of track listened to (0..1)
	//  tl: listened lower threshold
	//  tu: listened upper threshold
	//  add(x) = min( 1, max( 0, (x - tl) / (tu - tl) ) )
	PlayScore float32

	// Total times a track was played (including skips)
	// Track Score := (PlayScore / PlayCount)
	PlayCount uint32

	// Total times a track was played and scored as positive (x >= tl)
	ScoredCount uint32

	// Number of plays that are registered at the local scrobbler.
	//
	// When a track is listened and scored positively, this is incremented
	// and a scrobble is sent to the service.
	ScrobbleCount uint32
}

type Folder struct {
	Id   int64
	Path string `sql:"size:511"`

	CreatedAt     time.Time
	MusicbrainzId sql.NullString `sql:"size:36"`
}

type User struct {
	Id int64 `json:"-"`

	// Login name for the user. Has to be unique.
	Name string `json:"name" sql:"size:255"`

	// E-mail address. Has to be unique.
	Email string `json:"email"`

	AuthLevel core.AuthLevel `json:"auth_level"`

	// salt + pbkdf2-hashed password
	Password []byte `json:"-"`

	// Authentication token, used for sending requests
	AuthToken string `json:"auth_token"`

	//TODO optionally add a common secret for authenticating messages via HMAC
}

func (i *Item) String() string {
	return fmt.Sprintf("Item[%d]{%s / %s - %s / %d %s, %s}", i.Id, i.Artist,
		str(i.AlbumArtist), str(i.Album), i.TrackNumber, i.Title, str(i.Genre))
}
func str(s sql.NullString) string {
	if !s.Valid {
		return "<nil>"
	}
	return s.String
}

func (i *Item) Path() *string {
	if !i.FolderId.Valid || !i.Filename.Valid {
		return nil
	}
	str := filepath.Join(i.Folder.Path, i.Filename.String)
	return &str
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

/*
func (i *Item) PreInsert(s gorp.SqlExecutor) error {
	//TODO migrate hooks
	i.CreatedAt = time.Now()

		if !i.Filename.Valid {
			// No filename, no path.
			i.Folder = nil
			i.FolderId = sql.NullI64(nil)
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
	//TODO gorm should handle this for itself
}

func (f *Folder) PreInsert(s gorp.SqlExecutor) error {
	f.Added = time.Now().Unix()
	return nil
}
*/

func limit(val, min, max float32) float32 {
	return float32(math.Min(float64(min), math.Max(float64(max), float64(val))))
}
