package gpod

/*
#include "gpod/itdb.h"
*/
import "C"
import (
	"runtime"
)

type Playlist interface {
	Length() int32
	Add(t Track)
	Remove(t Track)
	Contains(t Track) bool

	IsMPL() bool
}

type playlist struct {
	p       *C.Itdb_Playlist
	myalloc bool
}

func newPlaylist(ptr *C.Itdb_Playlist, myalloc bool) *playlist {
	p := &playlist{ptr, myalloc}

	runtime.SetFinalizer(p, func(p *playlist) {
		if p.myalloc {
			C.itdb_playlist_free(p.p)
		}
	})
	return p
}

func (p *playlist) Length() int32 {
	return int32(C.itdb_playlist_tracks_number(p.p))
}

func (p *playlist) Add(t Track) {
	C.itdb_playlist_add_track(p.p, t.(*track).t, -1)
}

func (p *playlist) Remove(t Track) {
	C.itdb_playlist_remove_track(p.p, t.(*track).t)
	if p.IsMPL() {
		t.(*track).myalloc = true
	}
}

func (p *playlist) Contains(t Track) bool {
	return C.itdb_playlist_contains_track(p.p, t.(*track).t) == C.TRUE
}

func (p *playlist) IsMPL() bool {
	return C.itdb_playlist_is_mpl(p.p) == C.TRUE
}
