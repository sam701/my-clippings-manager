package main

import (
	"crypto/sha1"
	"encoding/hex"
)

type location struct {
	Start, End int
}

type book struct {
	Title   string
	Authors string
}

func (b book) getId() string {
	return getHash(b.Title + b.Authors)
}

func getHash(str string) string {
	bb := sha1.Sum([]byte(str))
	return hex.EncodeToString(bb[:])
}

type clipping struct {
	Book         book
	Loc          location
	CreationTime int64
	Content      string
}

func (c *clipping) getId() string {
	return getHash(c.Book.getId() + c.Content)
}
