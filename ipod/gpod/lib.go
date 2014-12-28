package gpod

/*
#cgo pkg-config: libgpod-1.0
#include "gpod/itdb.h"
#include "stdlib.h"

const guint32 RatingStep = ITDB_RATING_STEP;
*/
import "C"
import (
	"errors"
)

type GLib interface {
	Tracks() []Track
}

type Rating uint32

var (
	Rating1 Rating = Rating(C.RatingStep * 1)
	Rating2        = Rating(C.RatingStep * 2)
	Rating3        = Rating(C.RatingStep * 3)
	Rating4        = Rating(C.RatingStep * 4)
	Rating5        = Rating(C.RatingStep * 5)
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
	Size() uint32
	Length() int32
	Rating() Rating
	Playcount() uint32

	SetTitle(val string)
	SetAlbum(val string)
	SetArtist(val string)
	SetGenre(val string)
	SetFiletype(val string)
	SetComment(val string)
	SetComposer(val string)
	SetDescription(val string)
	SetAlbumartist(val string)
	SetSize(val uint32)
	SetLength(val int32)
	SetRating(val Rating)
	SetPlaycount(val uint32)
}

type glib struct {
	db *C.Itdb_iTunesDB
}

type track struct {
	t *C.Itdb_Track

	_Title       string
	_Album       string
	_Artist      string
	_Genre       string
	_Filetype    string
	_Comment     string
	_Composer    string
	_Description string
	_Albumartist string
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
		track := &track{t: (*C.Itdb_Track)(ptr.data)}
		arr = append(arr, track)
		ptr = ptr.next
	}
	return arr
}

func (t *track) Title() string {
	if t._Title == "" {
		t._Title = str(t.t.title)
	}
	return t._Title
}
func (t *track) Album() string {
	if t._Album == "" {
		t._Album = str(t.t.album)
	}
	return t._Album
}
func (t *track) Artist() string {
	if t._Artist == "" {
		t._Artist = str(t.t.artist)
	}
	return t._Artist
}
func (t *track) Genre() string {
	if t._Genre == "" {
		t._Genre = str(t.t.genre)
	}
	return t._Genre
}
func (t *track) Filetype() string {
	if t._Filetype == "" {
		t._Filetype = str(t.t.filetype)
	}
	return t._Filetype
}
func (t *track) Comment() string {
	if t._Comment == "" {
		t._Comment = str(t.t.comment)
	}
	return t._Comment
}
func (t *track) Composer() string {
	if t._Composer == "" {
		t._Composer = str(t.t.composer)
	}
	return t._Composer
}
func (t *track) Description() string {
	if t._Description == "" {
		t._Description = str(t.t.description)
	}
	return t._Description
}
func (t *track) Albumartist() string {
	if t._Albumartist == "" {
		t._Albumartist = str(t.t.albumartist)
	}
	return t._Albumartist
}
func (t *track) Size() uint32 {
	return uint32(t.t.size)
}
func (t *track) Length() int32 {
	return int32(t.t.tracklen)
}
func (t *track) Rating() Rating {
	return Rating(t.t.rating)
}
func (t *track) Playcount() uint32 {
	return uint32(t.t.playcount)
}

func (t *track) SetTitle(val string) {
	t._Title = val
	free(t.t.title)
	t.t.title = cstr(val)
}
func (t *track) SetAlbum(val string) {
	t._Album = val
	free(t.t.album)
	t.t.album = cstr(val)
}
func (t *track) SetArtist(val string) {
	t._Artist = val
	free(t.t.artist)
	t.t.artist = cstr(val)
}
func (t *track) SetGenre(val string) {
	t._Genre = val
	free(t.t.genre)
	t.t.genre = cstr(val)
}
func (t *track) SetFiletype(val string) {
	t._Filetype = val
	free(t.t.filetype)
	t.t.filetype = cstr(val)
}
func (t *track) SetComment(val string) {
	t._Comment = val
	free(t.t.comment)
	t.t.comment = cstr(val)
}
func (t *track) SetComposer(val string) {
	t._Composer = val
	free(t.t.composer)
	t.t.composer = cstr(val)
}
func (t *track) SetDescription(val string) {
	t._Description = val
	free(t.t.description)
	t.t.description = cstr(val)
}
func (t *track) SetAlbumartist(val string) {
	t._Albumartist = val
	free(t.t.albumartist)
	t.t.albumartist = cstr(val)
}
func (t *track) SetSize(val uint32) {
	t.t.size = C.guint32(val)
}
func (t *track) SetLength(val int32) {
	t.t.tracklen = C.gint32(val)
}
func (t *track) SetRating(val Rating) {
	t.t.rating = C.guint32(val)
}
func (t *track) SetPlaycount(val uint32) {
	t.t.playcount = C.guint32(val)
}
