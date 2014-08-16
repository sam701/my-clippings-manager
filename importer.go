package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
)

func (b book) getId() string {
	return getHash(b.Title + ":" + b.Authors)
}

func getHash(str string) string {
	bb := sha1.Sum([]byte(str))
	return hex.EncodeToString(bb[:])
}

type dbClipping struct {
	*baseClipping
	Note string
}

func (c *dbClipping) getId() string {
	return getHash(c.Book.getId() + ":" + c.Content + ":" + c.Note)
}

func newImporter() *importer {
	im := new(importer)
	im.booksProcessed = make(map[string]*book)
	im.booksImported = make(map[string]*book)
	im.bookClippings = make(map[string]int)
	return im
}

type importer struct {
	clippingsProcessed int
	clippingsImported  int
	emptyClippings     int
	booksProcessed     map[string]*book
	booksImported      map[string]*book
	bookClippings      map[string]int
	prev               *rawClipping
}

func (i *importer) importClippings(clippingFile string) {
	p := &parser{i.processRawClipping}
	p.parseClippingFile(clippingFile)
	i.printStat()
}

func shortenString(str string, n int) string {
	if len(str) > n {
		str = str[:n] + "..."
	}
	return str
}

func (i *importer) printStat() {
	fmt.Println("Clippings:")
	fmt.Println("  processed:\t", i.clippingsProcessed)
	fmt.Println("  empty:\t", i.emptyClippings)
	fmt.Println("  imported:\t", i.clippingsImported)
	fmt.Println("Books:")
	for k, v := range i.booksProcessed {
		prefix := "     "
		if i.booksImported[k] != nil {
			prefix = "(new)"
		}
		title := shortenString(v.Title, 50)
		fmt.Printf("  %s %5d %s\n", prefix, i.bookClippings[k], title)
	}
}

func (i *importer) processRawClipping(rc *rawClipping) {
	i.clippingsProcessed++

	if rc.Content == "" {
		i.emptyClippings++
	} else if rc.cType == highlight {
		c := &dbClipping{&rc.baseClipping, ""}
		if i.prev != nil && i.prev.cType == note {
			c.Note = i.prev.Content
		}
		i.importClipping(c)
	}

	i.prev = rc
}

func (i *importer) importClipping(c *dbClipping) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalln("Cannot start transaction:", err)
	}
	defer tx.Commit()

	bookId := c.Book.getId()

	_, err = tx.Exec(`insert into clipping (id, book, loc_start, loc_end, creation_time, content, note)
		values($1, $2, $3, $4, $5, $6, $7)`,
		c.getId(), bookId, c.Loc.Start, c.Loc.End, c.CreationTime, c.Content, c.Note)
	i.booksProcessed[bookId] = &c.Book
	i.bookClippings[bookId]++
	if err != nil {
		if isUniqueError(err) {
			return
		}
		log.Fatalln("Cannot insert clipping", err)
	}
	i.clippingsImported++

	_, err = tx.Exec("insert into book (id, title, authors) values($1, $2, $3)",
		bookId, c.Book.Title, c.Book.Authors)
	if err != nil {
		if isUniqueError(err) {
			return
		}
		log.Fatalln("Cannot insert book", err)
	}
	i.booksImported[bookId] = &c.Book

}

func isUniqueError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE")
}
