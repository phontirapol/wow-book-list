package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *BookDB {
	driver := "sqlite3"
	dataSource := "./db/book.db"
	initStatement := `
		CREATE TABLE IF NOT EXISTS
		book (
			book_id INTEGER PRIMARY KEY,
			title TEXT,
			author TEXT,
			genre TEXT,
			published_year INTEGER
		)
	`

	db := &BookDB{
		Driver:        driver,
		DataSource:    dataSource,
		InitStatement: initStatement,
	}

	err := db.Init()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

type BookDBInterface interface {
	GetDB() *sql.DB
}

type BookDB struct {
	Driver        string
	DataSource    string
	InitStatement string
	Database      *sql.DB
}

func (db *BookDB) Init() error {
	database, err := sql.Open(db.Driver, db.DataSource)
	if err != nil {
		return err
	}

	db.Database = database
	_, err = db.Database.Exec(db.InitStatement)
	if err != nil {
		return err
	}

	return nil
}

func (db *BookDB) GetDB() *sql.DB {
	return db.Database
}
