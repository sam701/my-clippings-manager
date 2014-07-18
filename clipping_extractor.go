package main

import (
	"flag"
)

func main() {
	fileToImport := flag.String("import", "", "File path to the 'My Clipping' file")
	dbFile := flag.String("db", "", "File path to the clipping DB (will be created if not exists)")
	flagListBooks := flag.Bool("list-books", false, "List all books in the database.")
	flag.Parse()

	if *dbFile == "" {
		flag.PrintDefaults()
		return
	}
	db := openDb(*dbFile)
	defer db.Close()

	if *flagListBooks {
		listBooks()
	} else if *fileToImport != "" {
		importFunc(*fileToImport)
	} else {
		flag.PrintDefaults()
	}

}
