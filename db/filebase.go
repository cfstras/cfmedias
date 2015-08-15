package db

import (
	"fmt"
	"os"
	osuser "os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/cfstras/cfmedias/config"
	"github.com/cfstras/cfmedias/core"
	errors "github.com/cfstras/cfmedias/errrs"
	log "github.com/cfstras/cfmedias/logger"
	"github.com/cfstras/cfmedias/vlc"
	"github.com/jinzhu/gorm"
)

type entry struct {
	Path     string
	Filename string
}

const bufferSize = 16

type updater struct {
	job          <-chan core.JobSignal
	tx           chan *gorm.DB
	allFiles     chan entry // all files seen
	newFiles     chan entry // files not yet in db
	importFiles  chan entry // files to import
	success      chan bool
	stopStepping chan bool
	vlc          *vlc.VLC

	// the receiving goroutine shall increment these
	numAllFiles      int
	numNewFiles      int
	numImportFiles   int
	numInvalidFiles  int
	numFailedFiles   int
	numImportedFiles int
}

var Types = []string{"mp3", "wav", "flac", "m4a", "aac"}

var ErrorTerminate error = errors.New("Terminating")

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
	d.db.LogMode(false)
	tx := d.db.Begin()
	up := &updater{
		job:          d.c.RegisterJob(),
		tx:           make(chan *gorm.DB, 1),
		allFiles:     make(chan entry, bufferSize),
		newFiles:     make(chan entry, bufferSize),
		importFiles:  make(chan entry, bufferSize),
		success:      make(chan bool),
		stopStepping: make(chan bool, 1),
	}
	up.tx <- tx

	var err error
	up.vlc, err = vlc.New()
	if err != nil {
		log.Log.Println("Could not init libVLC", err)
		return
	}

	go func() {
		err := filepath.Walk(searchPath, up.step)
		if err != nil && err != ErrorTerminate {
			log.Log.Println("Updater error:", err)
		}
		close(up.allFiles)
	}()

	// filter out ignored
	go func(input, output chan entry) {
		for entry := range input {
			up.numAllFiles++
			//TODO do something with the cover jpgs
			do := false
			for _, v := range Types {
				if strings.HasSuffix(entry.Filename, v) {
					do = true
					break
				}
			}
			if do {
				output <- entry
			}
		}
		close(output)
	}(up.allFiles, up.importFiles)

	// filter out already inserted files
	go func(input, output chan entry) {
		for entry := range input {
			up.numImportFiles++
			//fmt.Println("seen filter gets:", entry)
			// check if we already did this one
			tx := <-up.tx

			var c int
			err := tx.Table(ItemTable).
				//Joins("JOIN " + FolderTable + " ON " +
				//FolderTable + ".id = " + ItemTable + ".folder_id").
				//Where(entry).
				Where("filename = ? AND folder_id = (SELECT id FROM "+FolderTable+
				" WHERE path = ?)", entry.Filename, entry.Path).
				Count(&c).Error
			if err != nil {
				if strings.Contains(err.Error(), "no rows in result set") {
					log.Log.Println("this should not happen.")
					up.tx <- tx
					continue
				}
				log.Log.Println("sql error:", err)
				up.success <- false
			} else {
				up.tx <- tx
			}
			if c != 0 {
				// this one is already in the db
				//TODO check if the tags have changed anyway
				//log.Log.Println("skipping", entry)
			} else {
				//log.Log.Println("not skipping", entry)
				output <- entry
			}
		}
		close(output)
	}(up.importFiles, up.newFiles)

	// analyze remaining files
	go func(input chan entry, success chan bool) {
		for entry := range input {
			up.numNewFiles++
			err := up.analyze(path.Join(entry.Path, entry.Filename),
				entry.Path, entry.Filename)
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
			log.Log.Println("Rolling back...")
			if err := tx.Rollback().Error; err != nil {
				log.Log.Println("rollback error:", err)
			}
			break
		}
	}
	if success {
		tx := <-up.tx
		log.Log.Println("Committing imported files...")
		if err := tx.Commit().Error; err != nil {
			log.Log.Println("Updater error:", err)
		}
	}
	d.db.LogMode(false)
	log.Log.Println("Filebase updated:",
		"\nTotal Files:      ", up.numAllFiles,
		"\nNon-ignored Files:", up.numImportFiles,
		"\nNew Files:        ", up.numNewFiles,
		"\nImported Files:   ", up.numImportedFiles,
		"\nInvalid/Non-media:", up.numInvalidFiles,
		"\nFailed Files:     ", up.numFailedFiles)

	d.c.UnregisterJob(up.job)
}

func (up *updater) step(file string, info os.FileInfo, err error) error {
	if info == nil ||
		info.Name() == "." ||
		info.Name() == ".." {
		return nil
	}
	if info.IsDir() && info.Name() == ".git" {
		return filepath.SkipDir
	}
	select {
	case sig, ok := <-up.job:
		if !ok || sig >= core.SignalTerminate {
			log.Log.Println("Terminate got, processing remaining files")
			return ErrorTerminate
		}
	default:
	}

	if info.IsDir() {
		//log.Log.Println("in", file)
	} else if linked, err := filepath.EvalSymlinks(file); err != nil || file != linked {
		if err != nil {
			log.Log.Println("Error walking files:", err)
			return nil
		}
		//TODO add loop detection
		err = filepath.Walk(linked, up.step)
		if err != nil {
			return err
		}
	} else if !strings.HasPrefix(info.Name(), ".") {
		select {
		case <-up.stopStepping:
			return errors.New("aborting")
		case up.allFiles <- entry{path.Dir(file), info.Name()}:
		}
	}
	return nil
}

func (up *updater) analyze(path string, parent string, file string) error {
	//log.Log.Println("doing", path)

	tag, err := up.vlc.MediaNewPath(path)
	if err != nil {
		log.Log.Println("error reading file", path, "-", err)
		up.numInvalidFiles++
		return nil
	}
	defer tag.Release()
	tag.Parse()
	err = vlc.LastError()
	if err != nil {
		return err
	}

	title := tag.GetMeta(vlc.MetaTitle)
	artist := tag.GetMeta(vlc.MetaArtist)

	title, artist = TitleMagic(file, title, artist)

	item := Item{
		Title:       title,
		Artist:      artist,
		AlbumArtist: NullStr(nil),
		Album:       SqlStr(tag.GetMeta(vlc.MetaAlbum)),
		Genre:       SqlStr(tag.GetMeta(vlc.MetaGenre)),
		TrackNumber: SqlStr(tag.GetMeta(vlc.MetaTrackNumber)),
		Folder:      Folder{Path: parent},
		Filename:    SqlStr(file),
	}
	//TODO get albumartst, length, check ID etc
	err = vlc.LastError()
	if err != nil {
		return err
	}

	tx := <-up.tx
	err = tx.FirstOrCreate(&item.Folder, item.Folder).Error
	if err != nil {
		log.Log.Println("error inserting folder", item.Folder, err)
		return err
	}
	err = tx.Save(&item).Error
	up.tx <- tx
	if err != nil {
		log.Log.Println("error inserting item", item, err)
	} else {
		up.numImportedFiles++
		//log.Log.Println("inserted", item)
	}
	return err
}
