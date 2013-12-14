package db

import (
	"config"
	"github.com/coopernurse/gorp"
	log "logger"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"taglib"
)

type updater struct {
	tx *gorp.Transaction
}

func Update() {
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
	tx, err := dbmap.Begin()
	if err != nil {
		log.Log.Println("Could not start db transaction")
		return
	}
	up := &updater{tx}
	err = filepath.Walk(path, up.step)
	if err != nil {
		log.Log.Println(err)
	}

	if err = up.tx.Commit(); err != nil {
		log.Log.Println("Updater error:", err)
	}
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
		log.Log.Println(file)
		up.analyze(file)
	}
	return nil
}

func (up *updater) analyze(file string) {
	tag, err := taglib.Read(file)
	if err != nil {
		log.Log.Println("error reading file", file, "-", err)
		return
	}

	item := &Item{
		Title:       tag.Title(),
		Artist:      tag.Artist(),
		Genre:       tag.Genre(),
		TrackNumber: uint32(tag.Track())}

	//TODO get album, folder, added, check ID etc

	if err = up.tx.Insert(item); err != nil {
		log.Log.Println("error inserting item", item, err)
		return
	}
	log.Log.Println("inserted", item)
}
