package main

type location struct {
	Start, End int
}

type book struct {
	Title   string
	Authors string
}

type baseClipping struct {
	Book         book
	Loc          location
	CreationTime int64
	Content      string
}
