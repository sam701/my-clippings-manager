package main

import (
	"a4world/util/asignal"
	"flag"
	"log"
	"os/exec"
	"runtime"
)

var storage *clStorage

func main() {
	dbFile := flag.String("db", "", "File path to the clipping DB (will be created if not exists)")
	flag.Parse()

	if *dbFile == "" {
		flag.PrintDefaults()
		return
	}

	storage = NewStorage()
	err := openUrl(StartHttpServer())
	if err != nil {
		log.Fatalln(err)
	}

	asignal.WaitForShutdown(nil)
}

func openUrl(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd.exe", "/C", "start", url).Run()

	case "darwin":
		return exec.Command("open", url).Run()

	default:
		return exec.Command("xdg-open", url).Run()
	}
}
