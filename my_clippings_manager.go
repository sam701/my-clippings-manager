package main

//go:generate $GOPATH/bin/go-bindata -prefix web web
import (
	"a4world/util/asignal"
	"log"
	"os/exec"
	"runtime"
)

var storage *clStorage

func main() {
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
