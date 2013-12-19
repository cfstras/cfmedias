package db

import (
	"config"
	"core"
	"database/sql"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/go-contrib/uuid"
	_ "github.com/mattn/go-sqlite3"
	"io"
	log "logger"
)

type DB struct {
	dbmap *gorp.DbMap
	guid  string
}

func (d *DB) Open(c core.Core) error {
	if d.dbmap != nil {
		return errors.New("DB is already opened!")
	}

	var err error
	file := config.Current.DbFile
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return err
	}

	d.dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	//dbmap.TraceOn("[db]", log.Log)

	if err := d.checkTables(); err != nil {
		return err
	}

	err = d.checkSanity()
	if err != nil {
		//TODO what now?
		return err
	}

	d.initStats(c)

	c.RegisterCommand(core.Command{
		[]string{"rescan"},
		"Refreshes the database by re-scanning the media folder.",
		core.AuthAdmin,
		func(_ core.ArgMap, w io.Writer) error {
			fmt.Fprintln(w, "Rescanning media folder...")
			go d.Update()
			return nil
		}})

	return nil
}

func (d *DB) Close() error {
	if d.dbmap == nil {
		return errors.New("DB is not open!")
	}
	return d.dbmap.Db.Close()
}

// checks db schema and tables
func (d *DB) checkTables() error {
	//TODO replace this with gorp
	qu := `create table if not exists
	cfmedias
	(guid varchar(34) not null primary key, locked bool)`
	_, err := d.dbmap.Db.Exec(qu)
	if err != nil {
		log.Log.Println("SQL error", err, "at query:", qu)
		return err
	}

	qu = `select guid, locked from cfmedias`
	res, err := d.dbmap.Db.Query(qu)
	if err != nil {
		log.Log.Println("SQL error", err, "at query:", qu)
		return err
	}

	if res.Next() {
		var guidRead string
		var locked bool
		err = res.Scan(&guidRead, &locked)
		if err != nil {
			return err
		}
		d.guid = guidRead
		log.Log.Println("Database loaded with GUID", d.guid)
		//TODO set locked
	} else {
		d.guid = uuid.NewV4().String()
		locked := true

		qu = `insert into cfmedias values (?, ?)`
		_, err := d.dbmap.Db.Exec(qu, d.guid, locked)
		if err != nil {
			log.Log.Println("SQL error", err, "at query:", qu)
			return err
		}
		log.Log.Println("Database created with GUID", d.guid)
	}
	res.Close()

	d.dbmap.AddTableWithName(Item{}, ItemTable).SetKeys(true, "Id")
	//dbmap.AddTableWithName(Album{}, "albums").SetKeys(true, "Id")
	d.dbmap.AddTableWithName(Folder{}, FolderTable).SetKeys(true, "Id")

	err = d.dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Log.Println("Could not create database tables!")
		return err
	}

	return nil
}

// performs some integrity tests
func (d *DB) checkSanity() error {
	//TODO checkSanity()
	return nil
}
