package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type httpHandler struct {
	lastCallTime time.Time
}

// Returns the URL
func StartHttpServer() string {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		if path == "" {
			path = "index.html"
		}
		_, name := filepath.Split(path)
		b, err := Asset(path)
		if err == nil {
			info, _ := AssetInfo(path)
			http.ServeContent(w, r, name, info.ModTime(), bytes.NewReader(b))
		} else {
			log.Println("ERROR:", err)
			http.NotFound(w, r)
		}
	})

	h := &httpHandler{time.Now()}
	http.HandleFunc("/books", h.books)
	http.HandleFunc("/books/", h.clippings)
	http.HandleFunc("/upload", h.fileUpload)
	http.HandleFunc("/uploadIndex", h.fileUploadIndex)
	base := "127.0.0.1:3333"
	go http.ListenAndServe(base, h.callWrapper())
	go h.checkShutdown()
	return "http://" + base
}

func (h *httpHandler) callWrapper() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.lastCallTime = time.Now()
		http.DefaultServeMux.ServeHTTP(w, r)
	})
}

func (h *httpHandler) checkShutdown() {
	for t := range time.Tick(5 * time.Minute) {
		if t.Sub(h.lastCallTime) > 10*time.Minute {
			log.Println("Not used for more than 10min. Shutting down...")
			time.Sleep(time.Second)
			os.Exit(0)
		}
	}
}

func httpInternalError(w http.ResponseWriter, err error) {
	log.Println("ERROR:", err)
	w.WriteHeader(500)
	w.Write([]byte(fmt.Sprintln("Internal server error:", err)))
}

func writeJson(w http.ResponseWriter, obj interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	en := json.NewEncoder(w)
	return en.Encode(obj)
}

func (h *httpHandler) fileUploadIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, storage.uploadIndexFileName())
}

func (h *httpHandler) books(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, storage.bookIndexFileName())
}

func (h *httpHandler) clippings(w http.ResponseWriter, r *http.Request) {
	defer handleInternalError(w)
	bookId := strings.Split(r.URL.Path, "/")[2]
	http.ServeFile(w, r, storage.bookFileName(bookId))
}

func handleInternalError(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Println("ERROR: Recovered:", r)
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintln(r)))
	}
}

func (h *httpHandler) fileUpload(w http.ResponseWriter, r *http.Request) {
	defer handleInternalError(w)

	file, headers, err := r.FormFile("file")
	panicOnError(err)
	defer file.Close()

	uploadedFileData := importClippings(file, headers.Filename)
	writeJson(w, uploadedFileData)
}
