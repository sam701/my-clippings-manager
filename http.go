package main

import (
	"a4world/util/alog"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

type httpHandler struct{}

func StartHttpServer() {
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
	h := &httpHandler{}
	http.HandleFunc("/books", h.books)
	http.HandleFunc("/books/", h.clippings)
	http.HandleFunc("/upload", h.fileUpload)
	http.ListenAndServe("127.0.0.1:3333", nil)
}

func httpInternalError(w http.ResponseWriter, err error) {
	log.Println(alog.ERROR, err)
	w.WriteHeader(500)
	w.Write([]byte(fmt.Sprintln("Internal server error:", err)))
}

func writeJson(w http.ResponseWriter, obj interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	en := json.NewEncoder(w)
	return en.Encode(obj)
}

func (h *httpHandler) books(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("select id, title, authors from book order by title")
	if err != nil {
		httpInternalError(w, err)
		return
	}

	type book struct {
		Id, Title, Authors string
	}
	out := make([]*book, 0, 100)

	for rows.Next() {
		var b book
		err = rows.Scan(&b.Id, &b.Title, &b.Authors)
		if err != nil {
			httpInternalError(w, err)
			return
		}
		out = append(out, &b)
	}

	err = writeJson(w, out)
	if err != nil {
		httpInternalError(w, err)
	}
}

func (h *httpHandler) clippings(w http.ResponseWriter, r *http.Request) {
	bookId := strings.Split(r.URL.Path, "/")[2]

	rows, err := db.Query(`select loc_start, loc_end, creation_time, content
		from clipping
		where book = $1
		order by loc_start, creation_time
		`, bookId)
	if err != nil {
		httpInternalError(w, err)
		return
	}

	type clip struct {
		LocStart, LocEnd int
		CreationTime     int64
		Content          string
	}

	out := make([]*clip, 0, 500)
	for rows.Next() {
		var c clip
		err = rows.Scan(&c.LocStart, &c.LocEnd, &c.CreationTime, &c.Content)
		if err != nil {
			httpInternalError(w, err)
			return
		}
		out = append(out, &c)
	}

	writeJson(w, out)
}

func (h *httpHandler) fileUpload(w http.ResponseWriter, r *http.Request) {
	file, headers, err := r.FormFile("file")
	if err != nil {
		httpInternalError(w, err)
		return
	}
	defer file.Close()
	log.Println("HEADERS:", headers.Filename)
}
