// Download file at URL resuming when interrupted.
// Technique 50 from Go in Practice.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")

	if len(os.Args[1:]) != 1 {
		fmt.Printf("usage: %s URL\n", os.Args[0])
		os.Exit(1)
	}
	URL := os.Args[1]

	filename, err := getFilename(URL)
	if err != nil {
		log.Fatal(err)
	}

	// Append existing file or create new one.
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = download(URL, file, 10)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %d bytes\n", filename, fi.Size())
}

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

func download(URL string, file *os.File, retries int) error {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	fi, err := file.Stat()
	if err != nil {
		return err
	}
	current := fi.Size()
	if current > 0 {
		start := strconv.FormatInt(current, 10)
		// Range HTTP header specifies a range of bytes to download. This allows
		// you to request a file, starting where you left off.
		req.Header.Set("Range", "bytes="+start+"-")
	}

	cc := &http.Client{Timeout: 5 * time.Minute}
	res, err := cc.Do(req)
	if err != nil && hasTimedOut(err) {
		if retries > 0 {
			return download(URL, file, retries-1)
		}
		return err
	} else if err != nil {
		return err
	}

	if res.StatusCode/100 != 2 {
		if res.Header.Get("Content-Length") == "0" { // file already fully donwloaded
			return nil
		} else {
			return fmt.Errorf("GET %s: %s", URL, res.Status)
		}
	}

	if res.Header.Get("Accept-Ranges") != "bytes" { // server doesn't support serving partial files
		retries = 0
	}

	_, err = io.Copy(file, res.Body)
	if err != nil && hasTimedOut(err) {
		if retries > 0 {
			return download(URL, file, retries-1)
		}
		return err
	} else if err != nil {
		return err
	}

	return nil
}

func hasTimedOut(err error) bool {
	switch err := err.(type) {
	case *url.Error:
		if err, ok := err.Err.(net.Error); ok && err.Timeout() {
			return true
		}
	case net.Error:
		if err.Timeout() {
			return true
		}
	case *net.OpError:
		if err.Timeout() {
			return true
		}
	}

	errTxt := "use of closed network connection"
	if err != nil && strings.Contains(err.Error(), errTxt) {
		return true
	}

	return false
}
