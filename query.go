package main

import (
	"fmt"
	"log"
)

func listBooks() {
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
		fmt.Printf("%s\n  %s\n  %s\n", id[:10], title, authors)
	}
}
