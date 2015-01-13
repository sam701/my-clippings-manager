package main

import (
	"a4world/util/alog"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"time"
)

type bookImportData struct {
	data               *bookData
	clippingsProcessed int
	clippingsImported  int
}

func (s *bookImportData) modified() bool {
	return s.clippingsImported > 0
}

func makeHash(str string) string {
	bb := sha1.Sum([]byte(str))
	return hex.EncodeToString(bb[:])
}

type dbClipping struct {
	*baseClipping
	Note string
}

type importStat struct {
}

type importer struct {
	prev        *rawClipping
	currentTime int64

	importedBooks map[string]*bookImportData
}

func importClippings(r io.Reader, fileName string) *importStat {
	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	// Store the file first
	if !storage.saveUploadFile(bytes.NewReader(buf.Bytes()), fmt.Sprintf("%x", md5.Sum(buf.Bytes()))) {
		log.Println(alog.INFO, "File already exists. Skipping.")
		return nil
	}

	var i importer
	i.currentTime = time.Now().Unix()

	i.importedBooks = make(map[string]*bookImportData)

	p := &parser{i.processRawClipping}
	p.parse(bytes.NewReader(buf.Bytes()))

	// save
	for k, v := range i.importedBooks {
		if v.modified() {
			storage.saveBook(k, v.data)
		}
	}
	i.updateBookIndex()

	return nil
}

func (i *importer) updateBookIndex() {
	ix := storage.readBooksIndex()
	for k, v := range i.importedBooks {
		ib := ix[k]
		if ib == nil {
			ib = &indexBook{
				Title:   v.data.Title,
				Authors: v.data.Authors,
			}
			ix[k] = ib
		}
		ib.ClippingsNo = len(v.data.Clippings)
		ib.ClippingsTimes.First = 1 << 62
		for _, c := range v.data.Clippings {
			if c.CreationTime < ib.ClippingsTimes.First {
				ib.ClippingsTimes.First = c.CreationTime
			}
			if c.CreationTime > ib.ClippingsTimes.Last {
				ib.ClippingsTimes.Last = c.CreationTime
			}
		}
	}
	storage.saveBooksIndex(ix)
}

func (i *importer) processRawClipping(rc *rawClipping) {
	if rc.cType == highlight {
		c := &dbClipping{&rc.baseClipping, ""}
		if i.prev != nil && i.prev.cType == note {
			c.Note = i.prev.Content
		}
		i.importClipping(c)
	}

	i.prev = rc
}

func (i *importer) importClipping(c *dbClipping) {
	bc := &bookClipping{
		Loc:          c.Loc,
		CreationTime: c.CreationTime,
		ImportTime:   time.Now().Unix(),
		Content:      c.Content,
		Note:         c.Note,
	}
	bd := i.getBookData(c.Book)
	cid := bc.id()
	_, exists := bd.data.Clippings[cid]
	bd.clippingsProcessed++
	if !exists {
		bd.clippingsImported++
		bd.data.Clippings[cid] = bc
	}
}

func (i *importer) getBookData(b book) *bookImportData {
	bookId := b.id()
	data := i.importedBooks[bookId]
	if data == nil {
		bd := storage.readBook(bookId)
		if len(bd.Clippings) == 0 {
			bd = &bookData{
				Title:     b.Title,
				Authors:   b.Authors,
				Clippings: make(map[string]*bookClipping),
			}
		}
		data = &bookImportData{bd, 0, 0}
		i.importedBooks[bookId] = data
	}
	return data
}
