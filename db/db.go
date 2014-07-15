package db

import (
	"database/sql"
	"github.com/cfstras/cfmedias/config"
	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/errrs"
	log "github.com/cfstras/cfmedias/logger"
	"github.com/coopernurse/gorp"
	"github.com/go-contrib/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type DB dbstruct

type dbstruct struct {
	dbmap *gorp.DbMap
	guid  string
}

func (d *DB) Open(c core.Core) error {
	if d.dbmap != nil {
		return errrs.New("DB is already opened!")
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
	if d.dbmap == nil {
		return errrs.New("DB is not open!")
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
	d.dbmap.AddTableWithName(User{}, UserTable).SetKeys(true, "Id")

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
