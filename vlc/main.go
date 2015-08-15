/*
The packaged version of VLC is subject to the GNU LGPL, v2.1 and
Copyright (C) 1998-2009 VLC authors and VideoLAN. For details, see the files in
the vlc/ folder.
*/
package vlc

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lvlc

#include <stdlib.h>
#include <vlc/vlc.h>
*/
import "C"
import (
	"errors"

	"unsafe"
)

type VLC struct {
	instance *C.libvlc_instance_t
}

type Media C.struct_libvlc_media_t

type Meta C.libvlc_meta_t

var (
	MetaTitle       Meta = C.libvlc_meta_Title
	MetaArtist      Meta = C.libvlc_meta_Artist
	MetaGenre       Meta = C.libvlc_meta_Genre
	MetaCopyright   Meta = C.libvlc_meta_Copyright
	MetaAlbum       Meta = C.libvlc_meta_Album
	MetaTrackNumber Meta = C.libvlc_meta_TrackNumber
	MetaDescription Meta = C.libvlc_meta_Description
	MetaRating      Meta = C.libvlc_meta_Rating
	MetaDate        Meta = C.libvlc_meta_Date
	MetaSetting     Meta = C.libvlc_meta_Setting
	MetaURL         Meta = C.libvlc_meta_URL
	MetaLanguage    Meta = C.libvlc_meta_Language
	MetaNowPlaying  Meta = C.libvlc_meta_NowPlaying
	MetaPublisher   Meta = C.libvlc_meta_Publisher
	MetaEncodedBy   Meta = C.libvlc_meta_EncodedBy
	MetaArtworkURL  Meta = C.libvlc_meta_ArtworkURL
	MetaTrackID     Meta = C.libvlc_meta_TrackID
	/*MetaTrackTotal     Meta = C.libvlc_meta_TrackTotal
	MetaDirector        Meta = C.libvlc_meta_Director
	MetaSeason          Meta = C.libvlc_meta_Season
	MetaEpisode         Meta = C.libvlc_meta_Episode
	MetaShowName        Meta = C.libvlc_meta_ShowName
	MetaActors          Meta = C.libvlc_meta_Actors
	MetaAlbumArtist     Meta = C.libvlc_meta_AlbumArtist
	MetaDiscNumber      Meta = C.libvlc_meta_DiscNumber
	*/

	MetaTags map[string]Meta = map[string]Meta{
		"MetaTitle":       MetaTitle,
		"MetaArtist":      MetaArtist,
		"MetaGenre":       MetaGenre,
		"MetaCopyright":   MetaCopyright,
		"MetaAlbum":       MetaAlbum,
		"MetaTrackNumber": MetaTrackNumber,
		"MetaDescription": MetaDescription,
		"MetaRating":      MetaRating,
		"MetaDate":        MetaDate,
		"MetaSetting":     MetaSetting,
		"MetaURL":         MetaURL,
		"MetaLanguage":    MetaLanguage,
		"MetaNowPlaying":  MetaNowPlaying,
		"MetaPublisher":   MetaPublisher,
		"MetaEncodedBy":   MetaEncodedBy,
		"MetaArtworkURL":  MetaArtworkURL,
		"MetaTrackID":     MetaTrackID,
	}
)

// New creates a new instance of libVLC.
func New() (*VLC, error) {
	vlc := &VLC{}
	vlc.instance = C.libvlc_new(0, nil)
	return vlc, LastError()
}

func LastError() error {
	str := C.libvlc_errmsg()
	if str == nil {
		return nil
	}

	defer C.free(unsafe.Pointer(str))
	defer C.libvlc_clearerr()

	return errors.New(C.GoString(str))
}

// MediaNewPath creates a media for a certain file path.
func (vlc *VLC) MediaNewPath(path string) (*Media, error) {
	str := C.CString(path)
	defer C.free(unsafe.Pointer(str))
	media := C.libvlc_media_new_path(vlc.instance, str)
	return (*Media)(media), LastError()
}

// Parse parses a media.
// This fetches (local) art, meta data and tracks information.
// The method is synchronous.
func (m *Media) Parse() {
	C.libvlc_media_parse((*C.struct_libvlc_media_t)(m))
}

// Releases resources associated to a media.
func (m *Media) Release() {
	C.libvlc_media_release((*C.struct_libvlc_media_t)(m))
}

func (m *Media) GetMeta(meta Meta) string {
	str := C.libvlc_media_get_meta((*C.struct_libvlc_media_t)(m), C.libvlc_meta_t(meta))
	defer C.free(unsafe.Pointer(str))
	return C.GoString(str)
}
