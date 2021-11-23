package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

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
		if res.StatusCode == 416 && res.Header.Get("Content-Length") == "0" { // file already fully donwloaded
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
