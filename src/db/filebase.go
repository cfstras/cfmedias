package db

import (
	"config"
	"errrs"
	"github.com/coopernurse/gorp"
	log "logger"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"taglib"
)

type updater struct {
	tx *gorp.Transaction
}

var IgnoredTypes = []string{
	"jpg", "jpeg", "png", "gif",
	"nfo", "m3u", "log", "sfv", "txt",
	"cue"}

func (d *DB) Update() {
	// keep file base up to date
	path := config.Current.MediaPath
	if strings.Contains(path, "~") {
		user, err := user.Current()
		if err != nil {
			log.Log.Println("Error getting user home directory:", err)
			return
		}
		path = strings.Replace(path, "~", user.HomeDir, -1)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Log.Println("Error: Music path", path, "does not exist!")
		return
	}
	tx, err := d.dbmap.Begin()
	if err != nil {
		log.Log.Println("Could not start db transaction")
		return
	}
	up := &updater{tx}
	err = filepath.Walk(path, up.step)
	if err != nil {
		log.Log.Println("Updater error:", err)
		if err = up.tx.Commit(); err != nil {
			log.Log.Println("rollback error:", err)
		}
	} else if err = up.tx.Commit(); err != nil {
		log.Log.Println("Updater error:", err)
	}
	log.Log.Println("Filebase updated.")
}

func (up *updater) step(file string, info os.FileInfo, err error) error {
	if info == nil ||
		info.Name() == "." ||
		info.Name() == ".." {
		return nil
	}
	if info.IsDir() {
		//log.Println("in", file)
	} else if linked, err := filepath.EvalSymlinks(file); err != nil || file != linked {
		if err != nil {
			log.Log.Println("Error walking files:", err.Error())
			return nil
		}
		filepath.Walk(linked, up.step)
	} else {
		//log.Log.Println("at", file)
		if err := up.analyze(file, path.Dir(file), info.Name()); err != nil {
			log.Log.Println("analyze error", err)
			return err
		}
	}
	return nil
}

func (up *updater) analyze(path string, parent string, file string) error {
	// check if we already did this one
	itemPath := ItemPathView{}
	if err := up.tx.SelectOne(&itemPath,
		`select item_id Id, filename Filename, path Path
		from `+ItemTable+`
		join `+FolderTable+` on `+FolderTable+`.folder_id = `+ItemTable+`.folder_id
		where filename = ?
		and path = ?`,
		file, parent); err != nil {
		return err
	}
	if itemPath.Id != 0 {
		// this one is already in the db
		//TODO check if the tags have changed anyway
		return nil
	}

	tag, err := taglib.Read(path)
	if err != nil {
		for _, v := range IgnoredTypes {
			if strings.HasSuffix(file, v) {
				return nil //TODO do something with the covers
			}
		}
		log.Log.Println("error reading file", path, "-", err)
		return nil
	}

	title := tag.Title()
	artist := tag.Artist()
	if title == nil || artist == nil {
		return errrs.New("Title and Artist cannot be nil. File " + path)
	}
	item := &Item{
		Title:       *title,
		Artist:      *artist,
		AlbumArtist: nil,
		Album:       tag.Album(),
		Genre:       tag.Genre(),
		TrackNumber: uint32(tag.Track()),
		Folder:      &Folder{Path: parent},
		Filename:    &file,
	}

	//TODO get album, check ID etc

	if err = up.tx.Insert(item); err != nil {
		log.Log.Println("error inserting item", item, err)
		return err
	}
	log.Log.Println("inserted", item)
	return nil
}
