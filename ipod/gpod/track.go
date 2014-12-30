package gpod

/*
#include "gpod/itdb.h"
const guint32 RatingStep = ITDB_RATING_STEP;
*/
import "C"
import (
	"encoding/json"
	"runtime"
	"time"
)

type RatingT int32
type TrackField string

var (
	Rating1 RatingT = RatingT(C.RatingStep * 1)
	Rating2         = RatingT(C.RatingStep * 2)
	Rating3         = RatingT(C.RatingStep * 3)
	Rating4         = RatingT(C.RatingStep * 4)
	Rating5         = RatingT(C.RatingStep * 5)
)

var (
	Title       TrackField = "title"
	Album       TrackField = "album"
	Artist      TrackField = "artist"
	Genre       TrackField = "genre"
	Filetype    TrackField = "filetype"
	Comment     TrackField = "comment"
	Composer    TrackField = "composer"
	Description TrackField = "description"
	Albumartist TrackField = "albumartist"
	Size        TrackField = "size"
	Length      TrackField = "tracklen"
	Rating      TrackField = "rating"
	Playcount   TrackField = "playcount"
	TimeAdded   TrackField = "time_added"
)

type Track interface {
	Title() string
	Album() string
	Artist() string
	Genre() string
	Filetype() string
	Comment() string
	Composer() string
	Description() string
	Albumartist() string
	Size() int32
	Length() int32
	Rating() RatingT
	Playcount() int32
	TimeAdded() time.Time

	SetTitle(val string)
	SetAlbum(val string)
	SetArtist(val string)
	SetGenre(val string)
	SetFiletype(val string)
	SetComment(val string)
	SetComposer(val string)
	SetDescription(val string)
	SetAlbumartist(val string)
	SetSize(val int32)
	SetLength(val int32)
	SetRating(val RatingT)
	SetPlaycount(val int32)
	SetTimeAdded(val time.Time)
}

type track struct {
	t *C.Itdb_Track

	// cached strings
	strCache map[TrackField]string
	myalloc  bool
}

func newTrack(ptr *C.Itdb_Track, myalloc bool) *track {
	t := &track{ptr, map[TrackField]string{}, myalloc}

	runtime.SetFinalizer(t, func(t *track) {
		if myalloc {
			C.itdb_track_free(t.t)
		}
	})
	return t
}

func NewTrack() Track {
	ptr := C.itdb_track_new()
	return newTrack(ptr, true)
}

func (t *track) Title() string {
	if t.strCache[Title] == "" {
		t.strCache[Title] = str(t.t.title)
	}
	return t.strCache[Title]
}
func (t *track) Album() string {
	if t.strCache[Album] == "" {
		t.strCache[Album] = str(t.t.album)
	}
	return t.strCache[Album]
}
func (t *track) Artist() string {
	if t.strCache[Artist] == "" {
		t.strCache[Artist] = str(t.t.artist)
	}
	return t.strCache[Artist]
}
func (t *track) Genre() string {
	if t.strCache[Genre] == "" {
		t.strCache[Genre] = str(t.t.genre)
	}
	return t.strCache[Genre]
}
func (t *track) Filetype() string {
	if t.strCache[Filetype] == "" {
		t.strCache[Filetype] = str(t.t.filetype)
	}
	return t.strCache[Filetype]
}
func (t *track) Comment() string {
	if t.strCache[Comment] == "" {
		t.strCache[Comment] = str(t.t.comment)
	}
	return t.strCache[Comment]
}
func (t *track) Composer() string {
	if t.strCache[Composer] == "" {
		t.strCache[Composer] = str(t.t.composer)
	}
	return t.strCache[Composer]
}
func (t *track) Description() string {
	if t.strCache[Description] == "" {
		t.strCache[Description] = str(t.t.description)
	}
	return t.strCache[Description]
}
func (t *track) Albumartist() string {
	if t.strCache[Albumartist] == "" {
		t.strCache[Albumartist] = str(t.t.albumartist)
	}
	return t.strCache[Albumartist]
}
func (t *track) Size() int32 {
	return int32(t.t.size)
}
func (t *track) Length() int32 {
	return int32(t.t.tracklen)
}
func (t *track) Rating() RatingT {
	return RatingT(t.t.rating)
}
func (t *track) Playcount() int32 {
	return int32(t.t.playcount)
}
func (t *track) TimeAdded() time.Time {
	return time.Unix(int64(t.t.time_added), 0)
}

func (t *track) SetTitle(val string) {
	t.strCache[Title] = val
	free(t.t.title)
	t.t.title = cstr(val)
}
func (t *track) SetAlbum(val string) {
	t.strCache[Album] = val
	free(t.t.album)
	t.t.album = cstr(val)
}
func (t *track) SetArtist(val string) {
	t.strCache[Artist] = val
	free(t.t.artist)
	t.t.artist = cstr(val)
}
func (t *track) SetGenre(val string) {
	t.strCache[Genre] = val
	free(t.t.genre)
	t.t.genre = cstr(val)
}
func (t *track) SetFiletype(val string) {
	t.strCache[Filetype] = val
	free(t.t.filetype)
	t.t.filetype = cstr(val)
}
func (t *track) SetComment(val string) {
	t.strCache[Comment] = val
	free(t.t.comment)
	t.t.comment = cstr(val)
}
func (t *track) SetComposer(val string) {
	t.strCache[Composer] = val
	free(t.t.composer)
	t.t.composer = cstr(val)
}
func (t *track) SetDescription(val string) {
	t.strCache[Description] = val
	free(t.t.description)
	t.t.description = cstr(val)
}
func (t *track) SetAlbumartist(val string) {
	t.strCache[Albumartist] = val
	free(t.t.albumartist)
	t.t.albumartist = cstr(val)
}
func (t *track) SetSize(val int32) {
	t.t.size = C.guint32(val)
}
func (t *track) SetLength(val int32) {
	t.t.tracklen = C.gint32(val)
}
func (t *track) SetRating(val RatingT) {
	t.t.rating = C.guint32(val)
}
func (t *track) SetPlaycount(val int32) {
	t.t.playcount = C.guint32(val)
}
func (t *track) SetTimeAdded(val time.Time) {
	t.t.time_added = C.time_t(val.Unix())
}

func (t *track) MarshalJSON() ([]byte, error) {
	m := map[TrackField]interface{}{
		Title:       t.Title(),
		Album:       t.Album(),
		Artist:      t.Artist(),
		Genre:       t.Genre(),
		Filetype:    t.Filetype(),
		Comment:     t.Comment(),
		Composer:    t.Composer(),
		Description: t.Description(),
		Albumartist: t.Albumartist(),
		Size:        t.Size(),
		Length:      t.Length(),
		Rating:      t.Rating(),
		Playcount:   t.Playcount(),
		TimeAdded:   t.TimeAdded(),
	}
	return json.Marshal(m)
}
