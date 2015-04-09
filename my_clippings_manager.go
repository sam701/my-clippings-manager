package main

//go:generate $GOPATH/bin/go-bindata -prefix web web
import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

var storage *clStorage

func main() {
	storage = NewStorage()
	err := openUrl(StartHttpServer())
	if err != nil {
		log.Fatalln(err)
	}

	waitForShutdown()
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

func waitForShutdown() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-c
	log.Println("INFO: Shutting down...")
}
