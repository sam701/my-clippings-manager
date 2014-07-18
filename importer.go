package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func importFunc(clippingFile string) {
	f, err := os.Open(clippingFile)
	if err != nil {
		log.Fatalln("Cannot open file", err.Error())
	}
	defer f.Close()

	im := newImporter()
	im.read(f)
	im.printStat()
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
		title := shortenString(v.title, 50)
		fmt.Printf("  %s %5d %s\n", prefix, i.bookClippings[k], title)
	}
}

func shortenString(str string, n int) string {
	if len(str) > n {
		str = str[:n] + "..."
	}
	return str
}

func (i *importer) read(r io.Reader) {
	c := new(clipping)

	lineNo := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		lineNo++

		if line == "==========" {
			lineNo = 0
			if c != nil {
				i.processClipping(c)
			}
			c = new(clipping)
			continue
		}

		switch lineNo {
		case 1:
			i.extractBook(line, c)
		case 2:
			i.extractLocationAndDate(line, c)
		case 3:
			continue
		default:
			if c.content != "" {
				c.content += "\n"
			}
			c.content += line
		}

	}
}

func (i *importer) extractBook(str string, c *clipping) {
	ix := strings.LastIndex(str, " (")
	if ix < 0 {
		c.book.title = str
	} else {
		c.book.title = str[:ix]
		c.book.authors = str[ix+2 : len(str)-1]
	}
}

func (i *importer) extractLocationAndDate(str string, c *clipping) {
	ix := strings.LastIndex(str, " | ")

	i.extractLocation(str[:ix], c)
	i.extractAddDate(str[ix+3:], c)
}

func (i *importer) extractAddDate(str string, c *clipping) {
	ix := strings.Index(str, ",")
	dateStr := str[ix+2:]
	t, err := time.Parse("January 2, 2006 3:04:05 PM", dateStr)
	if err != nil {
		log.Fatalln("Cannot parse date:", dateStr)
	}
	c.creationTime = t.Unix()
}

func (i *importer) extractLocation(str string, c *clipping) {
	ix := strings.LastIndex(str, " ")
	pageStr := strings.Split(str[ix+1:], "-")
	ii, err := strconv.Atoi(pageStr[0])
	if err != nil {
		log.Fatalln("Cannot parse start page", pageStr[0], "in", str)
	}
	c.loc.start = ii

	if len(pageStr) == 2 {
		ii, err = strconv.Atoi(pageStr[1])
		if err != nil {
			log.Fatalln("Cannot pares end page", pageStr[1], "in", str)
		}
		c.loc.end = ii
	} else {
		c.loc.end = ii
	}
}

func (i *importer) processClipping(c *clipping) {
	i.clippingsProcessed++

	if c.content == "" {
		i.emptyClippings++
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalln("Cannot start transaction:", err)
	}
	defer tx.Commit()

	bookId := c.book.getId()

	_, err = tx.Exec(`insert into clipping (id, book, loc_start, loc_end, creation_time, content)
		values($1, $2, $3, $4, $5, $6)`,
		c.getId(), bookId, c.loc.start, c.loc.end, c.creationTime, c.content)
	i.booksProcessed[bookId] = &c.book
	i.bookClippings[bookId]++
	if err != nil {
		if isUniqueError(err) {
			return
		}
		log.Fatalln("Cannot insert clipping", err)
	}
	i.clippingsImported++

	_, err = tx.Exec("insert into book (id, title, authors) values($1, $2, $3)",
		bookId, c.book.title, c.book.authors)
	if err != nil {
		if isUniqueError(err) {
			return
		}
		log.Fatalln("Cannot insert book", err)
	}
	i.booksImported[bookId] = &c.book

}

func isUniqueError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE")
}
