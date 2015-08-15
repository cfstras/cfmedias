package vlc

/*
#include <stdlib.h>
#include <vlc/vlc.h>
*/
import "C"

type Player struct {
	instance *C.libvlc_media_player_t
}

func (vlc *VLC) NewPlayer() (*Player, error) {
	p := &Player{}
	p.instance = C.libvlc_media_player_new(vlc.instance)
	return p, LastError()
}

func (p *Player) SetMedia(media *Media) {
	C.libvlc_media_player_set_media(p.instance, (*C.struct_libvlc_media_t)(media))
}

func (p *Player) Play() error {
	ret := C.libvlc_media_player_play(p.instance)
	if ret != 0 {
		return LastError()
	}
	return nil
}
