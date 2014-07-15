package db

import (
	"fmt"
	"github.com/cfstras/cfmedias/config"
	"github.com/cfstras/cfmedias/errrs"
	log "github.com/cfstras/cfmedias/logger"
	"github.com/cfstras/go-taglib"
	"github.com/coopernurse/gorp"
	"os"
	osuser "os/user"
	"path"
	"path/filepath"
	"strings"
)

type entry struct {
	folder string
	file   string
}

const bufferSize = 128

type updater struct {
	tx           chan *gorp.Transaction
	allFiles     chan entry // all files seen
	newFiles     chan entry // files not yet in db
	importFiles  chan entry // files to import
	success      chan bool
	stopStepping chan bool

	// the receiving goroutine shall increment these
	numAllFiles      int
	numNewFiles      int
	numImportFiles   int
	numInvalidFiles  int
	numFailedFiles   int
	numImportedFiles int
}

var IgnoredTypes = []string{
	"jpg", "jpeg", "png", "gif", "nfo", "m3u", "log", "sfv", "txt", "cue",
	"itc2", "html", "xml", "ipa", "asd", "plist", "itdb", "itl", "tmp", "ini",
	"sh", "sha1", "blb"}

func (d *DB) Update() {
	// keep file base up to date
	searchPath := config.Current.MediaPath
	if strings.Contains(searchPath, "~") {
		user, err := osuser.Current()
		if err != nil {
			log.Log.Println("Error getting user home directory:", err)
			return
		}
		searchPath = strings.Replace(searchPath, "~", user.HomeDir, -1)
	}
	if _, err := os.Stat(searchPath); os.IsNotExist(err) {
		log.Log.Println("Error: Music path", searchPath, "does not exist!")
		return
	}
	tx, err := d.dbmap.Begin()
	if err != nil {
		log.Log.Println("Could not start db transaction")
		return
	}
	up := &updater{tx: make(chan *gorp.Transaction, 1),
		allFiles:     make(chan entry, bufferSize),
		newFiles:     make(chan entry, bufferSize),
		importFiles:  make(chan entry, bufferSize),
		success:      make(chan bool),
		stopStepping: make(chan bool, 1)}
	up.tx <- tx

	go func() {
		err := filepath.Walk(searchPath, up.step)
		if err != nil {
			log.Log.Println("Updater error:", err)
		}
		close(up.allFiles)
	}()

	go func(input, output chan entry) {
		for entry := range input {
			up.numAllFiles++
			//fmt.Println("suffix filter gets:", entry)
			do := true
			for _, v := range IgnoredTypes {
				if strings.HasSuffix(entry.file, v) {
					//TODO do something with the cover jpgs
					do = false
					break
				}
			}
			if do {
				output <- entry
			}
		}
		close(output)
	}(up.allFiles, up.importFiles)

	go func(input, output chan entry) {
		for entry := range input {
			up.numImportFiles++
			//fmt.Println("seen filter gets:", entry)
			// check if we already did this one
			itemPath := ItemPathView{}
			tx := <-up.tx
			if err := tx.SelectOne(&itemPath,
				`select item_id Id, filename Filename, path Path
				from `+ItemTable+`
				join `+FolderTable+` on `+FolderTable+`.folder_id = `+ItemTable+`.folder_id
				where filename = ?
				and path = ?`,
				entry.file, entry.folder); err != nil {
				fmt.Println("sql error:", err)
				up.success <- false
			}
			up.tx <- tx
			if itemPath.Id != 0 {
				// this one is already in the db
				//TODO check if the tags have changed anyway
				//log.Log.Println("skipping", entry)
			} else {
				output <- entry
			}
		}
		close(output)
	}(up.importFiles, up.newFiles)

	go func(input chan entry, success chan bool) {
		for entry := range input {
			up.numNewFiles++
			err := up.analyze(path.Join(entry.folder, entry.file),
				entry.folder, entry.file)
			if err != nil {
				up.numFailedFiles++
				fmt.Println("import error: ", err)
			}
		}
		up.success <- true
		close(up.success)
	}(up.newFiles, up.success)

	success := true
	for v := range up.success {
		if !v {
			up.stopStepping <- true
			success = false
			break
			if err = tx.Rollback(); err != nil {
				log.Log.Println("rollback error:", err)
			}
		}
	}
	if success {
		tx := <-up.tx
		if err = tx.Commit(); err != nil {
			log.Log.Println("Updater error:", err)
		}
	}
	log.Log.Println("Filebase updated:",
		"\nTotal Files:      ", up.numAllFiles,
		"\nNon-ignored Files:", up.numImportFiles,
		"\nNew Files:        ", up.numNewFiles,
		"\nImported Files:   ", up.numImportedFiles,
		"\nInvalid/Non-media:", up.numInvalidFiles,
		"\nFailed Files:     ", up.numFailedFiles)
}

func (up *updater) step(file string, info os.FileInfo, err error) error {
	if info == nil ||
		info.Name() == "." ||
		info.Name() == ".." {
		return nil
	}
	if info.IsDir() {
		//log.Log.Println("in", file)
	} else if linked, err := filepath.EvalSymlinks(file); err != nil || file != linked {
		if err != nil {
			log.Log.Println("Error walking files:", err.Error())
			return nil
		}
		filepath.Walk(linked, up.step)
	} else {
		select {
		case <-up.stopStepping:
			return errrs.New("aborting")
		case up.allFiles <- entry{path.Dir(file), info.Name()}:
		}
	}
	return nil
}

func (up *updater) analyze(path string, parent string, file string) error {
	//log.Log.Println("doing", path)

	tag, err := taglib.Read(path)
	if err != nil {
		log.Log.Println("error reading file", path, "-", err)
		up.numInvalidFiles++
		return nil
	}

	defer tag.Close()

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

	tx := <-up.tx
	err = tx.Insert(item)
	up.tx <- tx
	if err != nil {
		log.Log.Println("error inserting item", item, err)
	} else {
		up.numImportedFiles++
		log.Log.Println("inserted", item)
	}
	return err
}
