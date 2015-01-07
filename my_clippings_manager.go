package main

import (
	"a4world/util/asignal"
	"flag"
)

func main() {
	dbFile := flag.String("db", "", "File path to the clipping DB (will be created if not exists)")
	flag.Parse()

	if *dbFile == "" {
		flag.PrintDefaults()
		return
	}
	db := openDb(*dbFile)
	defer db.Close()

	StartHttpServer()
	asignal.WaitForShutdown(nil)
}
