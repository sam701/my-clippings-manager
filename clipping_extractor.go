package main

import (
	"flag"
)

func main() {
	fileToImport := flag.String("import", "", "File path to the 'My Clipping' file")
	dbFile := flag.String("db", "", "File path to the clipping DB (will be created if not exists)")
	flag.Parse()

	if *fileToImport == "" || *dbFile == "" {
		flag.PrintDefaults()
		return
	}

	db := openDb(*dbFile)
	defer db.Close()

	importFunc(*fileToImport)
}
