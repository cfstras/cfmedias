package gpod

/*
#cgo pkg-config: libgpod-1.0
#include "gpod/itdb.h"
#include "stdlib.h"
*/
import "C"
import (
	"runtime"
)

type GLib interface {
	Tracks() []Track
	MPL() Playlist
	Add(t Track)
	Remove(t Track)
	Copy(t Track, path string) error
	Save() error
}

type glib struct {
	db  *C.Itdb_iTunesDB
	mpl *playlist
}

func New(path string) (GLib, error) {
	var gerr *C.GError
	mntpoint := cstr(path)
	itdb := C.itdb_parse(mntpoint, &gerr)
	free(mntpoint)
	if itdb == nil {
		return nil, err(gerr)
	}

	lib := &glib{itdb,
		newPlaylist(C.itdb_playlist_mpl(itdb), false)}

	runtime.SetFinalizer(itdb, func(l *glib) {
		C.itdb_free(l.db)
	})

	return lib, nil
}

func (g *glib) Tracks() []Track {
	ptr := g.db.tracks
	len := C.g_list_length(ptr)
	arr := make([]Track, 0, len)
	for ptr != nil {
		track := newTrack((*C.Itdb_Track)(ptr.data), false)
		arr = append(arr, track)
		ptr = ptr.next
	}
	return arr
}

func (g *glib) MPL() Playlist {
	return g.mpl
}

func (g *glib) Add(t Track) {
	C.itdb_track_add(g.db, t.(*track).t, -1)
	t.(*track).myalloc = false
}
func (g *glib) Remove(t Track) {
	C.itdb_track_remove(t.(*track).t)
	t.(*track).myalloc = true
}

func (g *glib) Copy(t Track, path string) error {
	pathC := cstr(path)
	defer free(pathC)
	var gerr *C.GError
	success := C.itdb_cp_track_to_ipod(t.(*track).t, pathC, &gerr)
	if success != C.TRUE {
		return err(gerr)
	}
	return nil
}

func (g *glib) Save() error {
	//TODO itdb_spl_update_all

	var gerr *C.GError
	success := C.itdb_write(g.db, &gerr)
	if success != C.TRUE {
		return err(gerr)
	}
	return nil
}
