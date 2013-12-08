package db

import (
	"config"
	log "logger"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"taglib"
)

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
	err := filepath.Walk(path, step)
	if err != nil {
		log.Log.Println(err)
	}
}

func step(file string, info os.FileInfo, err error) error {
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
		filepath.Walk(linked, step)
	} else {
		//TODO do something with it!
		log.Log.Println(file)
		Analyze(file)
	}
	return nil
}

func Analyze(file string) {
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

	err = dbmap.Insert(item)
	if err != nil {
		log.Log.Println("error inserting item", item, err)
		return
	}
	log.Log.Println("inserted", item)
}
