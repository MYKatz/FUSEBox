package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type fusebox struct {
	path       string
	done       chan bool
	watcher    *fsnotify.Watcher
	filedigest *digest
	netlify    *netlifySite
}

func newFusebox(path string, siteID string, netlifyToken string) *fusebox {
	fb := fusebox{path: path}
	fb.done = make(chan bool)
	fb.filedigest = newDigest()

	var err error
	fb.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	setupDirectory(path)
	fb.filedigest.resetWithPath(path)

	fb.netlify = &netlifySite{siteID: siteID, accessKey: netlifyToken}

	err = fb.watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	go fb.watch() // run watcher

	return &fb
}

func setupDirectory(path string) {
	log.Printf("Setting up fusebox directory")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0777) // make dir as a directory (perm 1<<31)
		log.Printf("Creating directory %s \n", path)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Printf("Using existant directory %s \n", path)
	}
}

func (fb fusebox) watch() {
	for {
		select {
		case <-fb.done:
			return
		case event, ok := <-fb.watcher.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
			}
		case err, ok := <-fb.watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func (fb fusebox) Stop() {
	fb.watcher.Close()
	fb.done <- true
}

func (fb fusebox) Start() {
	var err error
	fb.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
}

func (fb fusebox) debug() {
	fmt.Printf("{path: %s} \n", fb.path)
}

func (fb fusebox) update() {
	fb.filedigest.resetWithPath(fb.path)
	res, _ := fb.netlify.sendDigest(fb.filedigest.json())
	for _, fingerprint := range res.Required {
		if val, ok := fb.filedigest.inverted[fingerprint]; ok {
			err := fb.netlify.putFile(res.ID, val)
			if err != nil {
				log.Fatal(err) // TODO: fail gracefully
			}
		}
	}
}
