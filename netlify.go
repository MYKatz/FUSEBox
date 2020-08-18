package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
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
	for _, path := range paths {
		hashstring, err := hashFile(path.absolute)
		if err != nil {
			log.Fatal(err)
		}

		d.files[path] = hashstring
		d.inverted[hashstring] = path
	}

	fmt.Println(d)
}

func hashFile(filepath string) (string, error) {

	var hash string

	file, err := os.Open(filepath)
	if err != nil {
		return hash, err
	}

	defer file.Close()

	s1 := sha1.New()

	_, err = io.Copy(s1, file)
	if err != nil {
		return hash, err
	}

	bytes := s1.Sum(nil)[:20]

	return hex.EncodeToString(bytes), nil
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
