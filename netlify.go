package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const netlifyAPIURL = "https://api.netlify.com/api/v1/"

type fpath struct {
	absolute string
	relative string // relative to fusebox folder
}

type digest struct {
	files    map[fpath]string
	inverted map[string]fpath
}

type netlifySite struct {
	accessKey string
	siteID    string
}

type netlifyResponse struct {
	ID       string   `json:"id"`
	Required []string `json:"required"`
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

func (d digest) json() string {
	out := make(map[string]map[string]string)
	out["files"] = make(map[string]string)
	for key, val := range d.files {
		out["files"][key.relative] = val
	}
	j, err := json.Marshal(out)
	if err != nil {
		log.Fatal(err)
	}

	return string(j)
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

func (ns netlifySite) sendDigest(digest string) (string, error) {
	var response string
	requestURL := fmt.Sprintf("%ssites/%s/deploys", netlifyAPIURL, ns.siteID)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer([]byte(digest)))
	if err != nil {
		return response, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ns.accessKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	response = string(body)
	return response, nil
}
