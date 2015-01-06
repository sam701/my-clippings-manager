package main

import (
	"a4world/util/asignal"
	"flag"
)

func main() {
	fileToImport := flag.String("import", "", "File path to the 'My Clipping' file")
	dbFile := flag.String("db", "", "File path to the clipping DB (will be created if not exists)")
	flagListBooks := flag.Bool("list-books", false, "List books in the database.")
	bookHtmlFile := flag.String("book-html", "", "Output file name for book html.")
	bookId := flag.String("book-id", "", "Book ID as listen with \"list-books\"")
	flag.Parse()

	if *dbFile == "" {
		flag.PrintDefaults()
		return
	}
	db := openDb(*dbFile)
	defer db.Close()

	StartHttpServer()
	asignal.WaitForShutdown(nil)

	if *flagListBooks {
		oracle{}.listBooks()
	} else if *fileToImport != "" {
		newImporter().importClippings(*fileToImport)
	} else if *bookHtmlFile != "" && *bookId != "" {
		createBookHtml(*bookId, *bookHtmlFile)
	} else {
		flag.PrintDefaults()
	}

}
