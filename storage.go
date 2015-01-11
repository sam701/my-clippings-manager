package main

import (
	"a4world/util/alog"
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
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
	UploadTime       int64
	FileName         string
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
	os.MkdirAll(dir, 0755)
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

func (s *clStorage) uploadArchiveTarFileName() string {
	return filepath.Join(s.dir, "uploadsArchive.tar")
}

func (s *clStorage) uploadIndexFileName() string {
	return filepath.Join(s.dir, "uploadsIndex.json")
}

func (s *clStorage) readUploadsIndex() uploadsIndex {
	f, err := os.Open(s.uploadIndexFileName())
	panicOnError(err)
	defer f.Close()
	var ix uploadsIndex
	err = json.NewDecoder(f).Decode(&ix)
	panicOnError(err)
	return ix
}

func (s *clStorage) saveUploadItem(r io.Reader, item *uploadItem) {
	ix := s.readUploadsIndex()

	// Write index
	f, err := os.OpenFile(s.uploadIndexFileName(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0611)
	panicOnError(err)
	defer f.Close()

	ix = append(ix, item)
	panicOnError(json.NewEncoder(f).Encode(ix))
}

// Appends to the tar archive
func (s *clStorage) saveUploadFile(r io.Reader, fileName string) {
	f, err := os.OpenFile(s.uploadArchiveTarFileName(), os.O_CREATE|os.O_RDWR, 0644)
	panicOnError(err)
	defer f.Close()

	// gzip first
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	io.Copy(gw, r)
	gw.Close()

	stat, err := f.Stat()
	panicOnError(err)
	if stat.Size() > 1024 { // Skip the tar tail
		f.Seek(-1024, os.SEEK_END)
	}
	w := tar.NewWriter(f)
	defer w.Close()

	size := len(buf.Bytes())
	w.WriteHeader(&tar.Header{
		Name:    fileName + ".gz",
		Size:    int64(size),
		Mode:    0644,
		ModTime: time.Now(),
	})
	_, err = io.Copy(w, bytes.NewReader(buf.Bytes()))
	panicOnError(err)
	log.Println(alog.DEBUG, "Saved uploaded file:", fileName)
}

func (s *clStorage) bookFileName(bookId string) string {
	return filepath.Join(s.dir, "book"+bookId+".json")
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
