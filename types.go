package main

import (
	"crypto/sha1"
	"encoding/hex"
)

type location struct {
	start, end int
}

type book struct {
	title   string
	authors string
}

func (b book) getId() string {
	return getHash(b.title + b.authors)
}

func getHash(str string) string {
	bb := sha1.Sum([]byte(str))
	return hex.EncodeToString(bb[:])
}

type clipping struct {
	book         book
	loc          location
	creationTime int64
	content      string
}

func (c *clipping) getId() string {
	return getHash(c.book.getId() + c.content)
}
