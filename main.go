package main

import (
	"fmt"
	"log"
	"os"
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

	if err := download(URL, file, 10); err != nil {
		log.Fatal(err)
	}

	if fi, err := file.Stat(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s %d bytes\n", filename, fi.Size())
	}
}
