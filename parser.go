package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type clippingHandler func(*rawClipping)

type parser struct {
	handler clippingHandler
}

func (p *parser) parseClippingFile(clippingFile string) {
	f, err := os.Open(clippingFile)
	if err != nil {
		log.Fatalln("Cannot open file", err.Error())
	}
	defer f.Close()

	p.parse(f)
}

type clippingType int

const (
	undefined clippingType = iota
	highlight
	note
	bookmark
)

type rawClipping struct {
	baseClipping
	cType clippingType
}

func (i *parser) parse(r io.Reader) {
	c := new(rawClipping)

	lineNo := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		lineNo++

		if line == "==========" {
			lineNo = 0
			if c != nil {
				i.handler(c)
			}
			c = new(rawClipping)
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
			if c.Content != "" {
				c.Content += "\n"
			}
			c.Content += line
		}

	}
}

func (i *parser) extractBook(str string, c *rawClipping) {
	ix := strings.LastIndex(str, " (")
	if ix < 0 {
		c.Book.Title = str
	} else {
		c.Book.Title = str[:ix]
		c.Book.Authors = str[ix+2 : len(str)-1]
	}
}

func (i *parser) extractLocationAndDate(str string, c *rawClipping) {
	ix := strings.LastIndex(str, " | ")

	i.extractLocation(str[:ix], c)
	i.extractAddDate(str[ix+3:], c)
	c.cType = extractType(str)
}

func extractType(str string) clippingType {
	if strings.Contains(str, "Your Highlight") {
		return highlight
	} else if strings.Contains(str, "Your Note") {
		return note
	} else if strings.Contains(str, "Your Bookmark") {
		return bookmark
	} else {
		log.Fatalln("Cannot deternime clipping type:", str)
		return undefined
	}
}

func (i *parser) extractAddDate(str string, c *rawClipping) {
	ix := strings.Index(str, ",")
	dateStr := str[ix+2:]
	t, err := time.Parse("January 2, 2006 3:04:05 PM", dateStr)
	if err != nil {
		log.Fatalln("Cannot parse date:", dateStr)
	}
	c.CreationTime = t.Unix()
}

func (i *parser) extractLocation(str string, c *rawClipping) {
	ix := strings.LastIndex(str, " ")
	pageStr := strings.Split(str[ix+1:], "-")
	ii, err := strconv.Atoi(pageStr[0])
	if err != nil {
		log.Fatalln("Cannot parse start page", pageStr[0], "in", str)
	}
	c.Loc.Start = ii

	if len(pageStr) == 2 {
		ii, err = strconv.Atoi(pageStr[1])
		if err != nil {
			log.Fatalln("Cannot pares end page", pageStr[1], "in", str)
		}
		c.Loc.End = ii
	} else {
		c.Loc.End = ii
	}
}
