package main

import (
	"github.com/wsxiaoys/terminal/color"
	"log"
)

type oracle struct{}

func (o oracle) listBooks() {
	rows, err := db.Query("select id, title, authors from book order by title")
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var id, title, authors string
		err = rows.Scan(&id, &title, &authors)
		if err != nil {
			log.Fatalln(err)
		}
		title = shortenString(title, 50)
		authors = shortenString(authors, 50)
		color.Printf("@{y!}%s@|\n  @!%s@|\n  %s\n", id[:10], title, authors)
	}
}
