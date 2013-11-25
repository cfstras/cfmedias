package db

import (
	"config"
	"fmt"
	"log"
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
			log.Println("Error getting user home directory:", err)
			return
		}
		path = strings.Replace(path, "~", user.HomeDir, -1)
	}
	err := filepath.Walk(path, step)
	fmt.Println(err)
}

func step(file string, info os.FileInfo, err error) error {
	if info.Name() == "." || info.Name() == ".." {
		return nil
	}
	if info.IsDir() {
		//log.Println("in", file)
	} else if linked, err := filepath.EvalSymlinks(file); err != nil || file != linked {
		if err != nil {
			log.Println("Error walking files:", err.Error())
			return nil
		}
		filepath.Walk(linked, step)
	} else {
		//TODO do something with it!
		log.Println(file)
		Analyze(file)
	}
	return nil
}

func Analyze(file string) {
	tag, err := taglib.Read(file)
	if err != nil {
		log.Println("error reading file", file, "-", err)
		return
	}

	title := tag.Title()
	artist := tag.Artist()
	album := tag.Album()

	fmt.Println(artist, "/", album, "-", title)
}
