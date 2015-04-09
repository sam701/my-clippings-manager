package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type bookData struct {
	Title     string
	Authors   string
	Clippings map[string]*bookClipping
}

type bookClipping struct {
	Loc          location
	CreationTime int64
	ImportTime   int64
	Content      string
	Note         string `json:",omitempty"`
}

func (c *bookClipping) id() string {
	return makeHash(c.Content + c.Note)
}

type indexBook struct {
	Title, Authors string
	ClippingsNo    int
	ClippingsTimes struct {
		First, Last int64
	}
}

type bookIndex map[string]*indexBook

type uploadItem struct {
	Id               string
	FileSize         int
	ClippingsTotalNo int
	EmptyClippingsNo int
	BooksNo          int
}

type uploadsIndex []*uploadItem

type clStorage struct {
	dir string
}

func NewStorage() *clStorage {
	dir := defaultStorageDir()
	log.Println("DEBUG: Storage dir:", dir)
	os.MkdirAll(filepath.Join(dir, "uploads"), 0755)
	os.MkdirAll(filepath.Join(dir, "books"), 0755)
	return &clStorage{dir}
}

func defaultStorageDir() string {
	home := os.Getenv("HOME")
	name := "MyClippingsManager"
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("LocalAppData"), name)

	case "darwin":
		return filepath.Join(home, "Library/Application Support", name)

	default:
		if xdgCfg := os.Getenv("XDG_CONFIG_HOME"); xdgCfg != "" {
			return filepath.Join(xdgCfg, name)
		}
		return filepath.Join(home, ".config", name)
	}
}

func (s *clStorage) uploadIndexFileName() string {
	return filepath.Join(s.dir, "uploadsIndex.json")
}

func (s *clStorage) readUploadsIndex() uploadsIndex {
	var ix uploadsIndex
	readJsonFile(s.uploadIndexFileName(), &ix)
	return ix
}

func (s *clStorage) saveUploadsIndex(ix uploadsIndex) {
	writeJsonFile(s.uploadIndexFileName(), ix)
}

// Returns false if such file already exists.
func (s *clStorage) saveUploadFile(r io.Reader, uploadId string) bool {
	fn := s.uploadFileName(uploadId)
	if _, err := os.Stat(fn); !os.IsNotExist(err) {
		return false
	}
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_RDWR, 0644)
	panicOnError(err)
	defer f.Close()

	// gzip first
	gw := gzip.NewWriter(f)
	defer gw.Close()
	io.Copy(gw, r)

	log.Println("DEBUG: Saved uploaded file:", uploadId)
	return true
}

func (s *clStorage) uploadFileName(uploadId string) string {
	return filepath.Join(s.dir, "uploads", uploadId+".txt.gz")
}

func (s *clStorage) bookFileName(bookId string) string {
	return filepath.Join(s.dir, "books", bookId+".json")
}

// Return if read
func readJsonFile(fileName string, out interface{}) bool {
	f, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return false
	}
	panicOnError(err)
	defer f.Close()

	err = json.NewDecoder(f).Decode(out)
	panicOnError(err)
	return true
}

func writeJsonFile(fileName string, data interface{}) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	panicOnError(err)
	defer f.Close()

	panicOnError(json.NewEncoder(f).Encode(data))
}

func (s *clStorage) readBook(bookId string) *bookData {
	var cf bookData
	if !readJsonFile(s.bookFileName(bookId), &cf) {
		cf.Clippings = make(map[string]*bookClipping)
	}
	return &cf
}

func (s *clStorage) saveBook(bookId string, b *bookData) {
	writeJsonFile(s.bookFileName(bookId), b)
}

func (s *clStorage) bookIndexFileName() string {
	return filepath.Join(s.dir, "bookIndex.json")
}

func (s *clStorage) readBooksIndex() bookIndex {
	var ix bookIndex
	if !readJsonFile(s.bookIndexFileName(), &ix) {
		ix = make(map[string]*indexBook)
	}
	return ix
}

func (s *clStorage) saveBooksIndex(ix bookIndex) {
	writeJsonFile(s.bookIndexFileName(), ix)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
