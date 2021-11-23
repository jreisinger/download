package main

import (
	"net/url"
	"os"
	"path"
)

func getFilename(URL string) (filename string, err error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	filename = path.Base(u.Path)
	return filename, nil
}

func fileExists(file *os.File) bool {
	_, err := file.Stat()
	return !os.IsNotExist(err)
}
