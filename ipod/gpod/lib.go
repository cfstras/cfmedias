package gpod

/*
#cgo pkg-config: libgpod-1.0
#include "gpod/itdb.h"
#include "stdlib.h"

const guint32 RatingStep = ITDB_RATING_STEP;
*/
import "C"
import (
	"encoding/json"
	"errors"
)

type GLib interface {
	Tracks() []Track
}

type RatingT int32

var (
	Rating1 RatingT = RatingT(C.RatingStep * 1)
	Rating2         = RatingT(C.RatingStep * 2)
	Rating3         = RatingT(C.RatingStep * 3)
	Rating4         = RatingT(C.RatingStep * 4)
	Rating5         = RatingT(C.RatingStep * 5)
)

type Field string

var (
	Title       Field = "title"
	Album       Field = "album"
	Artist      Field = "artist"
	Genre       Field = "genre"
	Filetype    Field = "filetype"
	Comment     Field = "comment"
	Composer    Field = "composer"
	Description Field = "description"
	Albumartist Field = "albumartist"
	Size        Field = "size"
	Length      Field = "tracklen"
	Rating      Field = "rating"
	Playcount   Field = "playcount"
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
}

type glib struct {
	db *C.Itdb_iTunesDB
}

type track struct {
	t *C.Itdb_Track

	// cached strings

	strCache map[Field]string
}

func New(path string) (GLib, error) {
	var gerr *C.GError
	mntpoint := cstr(path)
	itdb := C.itdb_parse(mntpoint, &gerr)
	free(mntpoint)

	if itdb == nil {
		str := str(gerr.message)
		return nil, errors.New(str)
	}

	return &glib{itdb}, nil
}

func (g *glib) Tracks() []Track {
	ptr := g.db.tracks
	len := C.g_list_length(ptr)
	arr := make([]Track, 0, len)
	for ptr != nil {
		track := &track{(*C.Itdb_Track)(ptr.data), map[Field]string{}}
		arr = append(arr, track)
		ptr = ptr.next
	}
	return arr
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

func (t *track) MarshalJSON() ([]byte, error) {
	m := map[Field]interface{}{
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
	}
	return json.Marshal(m)
}
