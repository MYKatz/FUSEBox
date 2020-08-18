package main

import (
	"flag"
	"log"
	"os/user"
	"path"

	"github.com/fsnotify/fsnotify"
)

func main() {

	user, err := user.Current()
	fuseboxPath := flag.String("path", path.Join(user.HomeDir, "fusebox"), "Path to your FUSEBox directory")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	fb := newFusebox(*fuseboxPath)
	fb.debug()

	<-done // blockhing channel read, program runs indefinitely
}
