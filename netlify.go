package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type fpath struct {
	absolute string
	relative string // relative to fusebox folder
}

type digest struct {
	files    map[fpath]string
	inverted map[string]fpath
}

func newDigest() *digest {
	d := digest{}
	d.files = make(map[fpath]string)
	d.inverted = make(map[string]fpath)
	return &d
}

func (d digest) resetWithPath(path string) {
	paths := filesInFolder(path)
	fmt.Println(paths)
}

func filesInFolder(folder string) []fpath {
	var paths []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	var files []fpath
	for _, path := range paths {
		f, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}

		if f.Mode().IsRegular() {
			relative, err := filepath.Rel(folder, path)
			if err != nil {
				log.Fatal(err)
			}
			files = append(files, fpath{absolute: path, relative: relative})
		}
	}

	return files
}
