package db

import (
	"config"
	"database/sql"
	"github.com/go-contrib/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB
var guid string

func Open() error {
	var err error
	file := config.Current.DbFile
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		return err
	}

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

// checks db schema and tables
func checkTables() error {
	qu := `create table if not exists
	cfmedias
	(guid varchar(34) not null primary key, locked bool)`
	_, err := db.Exec(qu)
	if err != nil {
		return e(qu, nil, err)
	}

	qu = `select guid, locked from cfmedias`
	res, err := db.Query(qu)
	if err != nil {
		return q(qu, res, err)
	}

	if res.Next() {
		var guidRead string
		var locked bool
		err = res.Scan(&guidRead, &locked)
		if err != nil {
			return err
		}
		guid = guidRead
		log.Println("Database loaded with GUID", guid)
		//TODO set locked
	} else {
		guid = uuid.NewV4().String()
		locked := true

		qu = `insert into cfmedias values (?, ?)`
		_, err := db.Exec(qu, guid, locked)
		if err != nil {
			return e(qu, nil, err)
		}
		log.Println("Database created with GUID", guid)
	}

	//TODO checkTables()

	return nil
}

func e(query string, res sql.Result, err error) error {
	if err != nil {
		log.Println("SQL error at query:", query)
		return err
	}
	return nil
}

func q(query string, res *sql.Rows, err error) error {
	if err != nil {
		log.Println("SQL error at query:", query)
		return err
	}
	return nil
}

// performs some integrity tests
func checkSanity() error {
	//TODO checkSanity()
	return nil
}
