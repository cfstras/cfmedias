package db

import (
	"config"
	"core"
	"errrs"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type DB dbstruct

type dbstruct struct {
	db   gorm.DB
	open bool
}

func (d *DB) Open(c core.Core) error {
	if d.open {
		return errrs.New("DB is already opened!")
	}

	var err error
	file := config.Current.DbFile

	d.db, err = gorm.Open("sqlite3", file)
	if err != nil {
		return err
	}
	d.open = true

	if err := d.checkTables(); err != nil {
		return err
	}

	err = d.checkSanity()
	if err != nil {
		//TODO what now?
		return err
	}

	d.initStats(c)
	d.initLogic(c) // hear the difference?
	d.initLogin(c) // it's subtle but it could save your life
	d.initList(c)

	c.RegisterCommand(core.Command{
		[]string{"rescan"},
		"Refreshes the database by re-scanning the media folder.",
		map[string]string{},
		core.AuthAdmin,
		func(_ core.CommandContext) core.Result {
			go d.Update()
			return core.ResultOK
		}})

	return nil
}

func (d *DB) Close() error {
	if !d.open {
		return errrs.New("DB is not open!")
	}
	return d.db.Close()
}

// checks db schema and tables
func (d *DB) checkTables() error {
	d.db.AutoMigrate(Item{})
	//db.AutoMigrate(Album{})
	d.db.AutoMigrate(Folder{})
	d.db.AutoMigrate(User{})

	return nil
}

// performs some integrity tests
func (d *DB) checkSanity() error {
	//TODO checkSanity()
	return nil
}
