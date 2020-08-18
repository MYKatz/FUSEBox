package main

import (
	"flag"
	"log"
	"os"
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
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	_, err = os.Stat(*fuseboxPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(*fuseboxPath, 0777) // make dir as a directory (perm 1<<31)
		log.Printf("Creating directory %s \n", *fuseboxPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("Using existant directory %s \n", *fuseboxPath)
	}

	err = watcher.Add(*fuseboxPath)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
