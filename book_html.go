package main

import (
	"html/template"
	"log"
	"os"
)

const (
	bookTemplate = `<!DOCTYPE html><html>
	<head>
		<meta charset="utf-8"/>
		<style>
			body {
				font-family: "Helvetica Neue",Helvetica,Arial,sans-serif;
				font-size: 14px;
				line-height: 1.42857143;
				color: #333;
			}
			.title {
				font-weight: bold;
				font-size: 24px;
			}
			.authors {
				font-weight: bold;
				font-size: 18px;
				color: #959595;
			}
			.clippings {
				margin-top: 15px;
			}
			.clipping {
				padding-bottom: 10px;
				border-top: 1px solid #ddd;
			}
			.note {
				margin-left: 1cm;
				color: #959595;
			}
		</style>
	</head>
	<body>
		<div class="title">{{.Book.Title}}</div>
		<div class="authors">{{.Book.Authors}}</div>
		<div class="clippings">
		{{range .Clippings}}
			<div class="clipping">
				<div class="content">{{.Content}}</div>
				{{ if .Note}}
					<div class="note">{{.Note}}</div>
				{{end}}
			</div>
		{{end}}
		</div>
	</body>
</html>`
)

type bookHtmlData struct {
	Book      book
	Clippings []*dbClipping
}

func createBookHtml(bookId string, outputFile string) {
	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	var data bookHtmlData

	// read book
	row := db.QueryRow("select id, title, authors from book where id like $1", bookId+"%")
	var bookHash string
	err = row.Scan(&bookHash, &data.Book.Title, &data.Book.Authors)
	failOnError(err)

	rows, err := db.Query(`select loc_start, loc_end, creation_time, content, note
		from clipping
		where book = $1
		order by loc_start, creation_time
		`, bookHash)
	failOnError(err)

	data.Clippings = make([]*dbClipping, 0, 500)
	for rows.Next() {
		c := new(dbClipping)
		c.baseClipping = new(baseClipping)
		err = rows.Scan(&c.Loc.Start, &c.Loc.End, &c.CreationTime, &c.Content, &c.Note)
		failOnError(err)
		data.Clippings = append(data.Clippings, c)
	}

	t := template.Must(template.New("book-html").Parse(bookTemplate))
	t.Execute(f, data)
}

func failOnError(err error) {
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
}
