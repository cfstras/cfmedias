package db

import (
	"config"
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/go-contrib/uuid"
	_ "github.com/mattn/go-sqlite3"
	log "logger"
)

var db *sql.DB
var dbmap *gorp.DbMap
var guid string

func Open() error {
	var err error
	file := config.Current.DbFile
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		return err
	}

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.TraceOn("[db]", log.Log)

	if err := checkTables(); err != nil {
		return err
	}

	err = checkSanity()
	if err != nil {
		//TODO what now?
		return err
	}

	return nil
}

func Close() error {
	return dbmap.Db.Close()
}

// checks db schema and tables
func checkTables() error {
	qu := `create table if not exists
	cfmedias
	(guid varchar(34) not null primary key, locked bool)`
	_, err := db.Exec(qu)
	if err != nil {
		log.Log.Println("SQL error", err, "at query:", qu)
		return err
	}

	qu = `select guid, locked from cfmedias`
	res, err := db.Query(qu)
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
		guid = guidRead
		log.Log.Println("Database loaded with GUID", guid)
		//TODO set locked
	} else {
		guid = uuid.NewV4().String()
		locked := true

		qu = `insert into cfmedias values (?, ?)`
		_, err := db.Exec(qu, guid, locked)
		if err != nil {
			log.Log.Println("SQL error", err, "at query:", qu)
			return err
		}
		log.Log.Println("Database created with GUID", guid)
	}
	res.Close()

	dbmap.AddTableWithName(Item{}, "items").SetKeys(true, "Id")
	dbmap.AddTableWithName(Album{}, "albums").SetKeys(true, "Id")
	dbmap.AddTableWithName(Folder{}, "folders").SetKeys(true, "Id")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		log.Log.Println("Could not create database tables!")
		return err
	}

	return nil
}

// performs some integrity tests
func checkSanity() error {
	//TODO checkSanity()
	return nil
}
