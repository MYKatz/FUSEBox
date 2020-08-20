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
	siteID := flag.String("siteid", "", "Netlify site ID")
	netlifyKey := flag.String("netlifykey", "", "Netlify API key")
	flag.Parse()

	if *siteID == "" {
		log.Fatal("Bad siteID")
	}
	if *netlifyKey == "" {
		log.Fatal("Bad netlifyKey")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	fb := newFusebox(*fuseboxPath, *siteID, *netlifyKey)
	fb.debug()

	<-done // blockhing channel read, program runs indefinitely
}
