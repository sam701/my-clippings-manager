package main

import (
	"database/sql"
	_ "github.com/mxk/go-sqlite/sqlite3"
	"log"
)

var db *sql.DB

func openDb(filePath string) *sql.DB {
	var err error
	if db != nil {
		return db
	}
	db, err = sql.Open("sqlite3", "file:"+filePath)
	if err != nil {
		log.Fatalln("Cannot connect to DB", err)
	}
	if !isDbInitialized(db) {
		initDb(db)
	}
	return db
}

func isDbInitialized(db *sql.DB) bool {
	_, err := db.Exec("select 1 from clipping")
	return err == nil
}

func initDb(db *sql.DB) {
	statements := [...]string{
		`create table book (
			id varchar(64) primary key,
			title varchar(1024),
			authors varchar(1024)
		);`,
		`create table clipping (
			id varchar(64) primary key,
			book varchar(64) not null,
			loc_start int,
			loc_end int,
			creation_time int,
			content text,
			note text,
			import_time int
		);`,
		`create index clipping_book on clipping(book);`,
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalln("Cannot init db, cannot start transaction:", err.Error())
	}
	defer tx.Commit()
	for _, s := range statements {
		_, err = tx.Exec(s)
		if err != nil {
			log.Println("ERROR in SQL:", s)
			log.Fatalln("Cannot init DB:", err.Error())
		}
	}

	log.Println("DB has been initialized.")
}
