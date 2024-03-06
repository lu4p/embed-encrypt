package main

import (
	"io/fs"
	"log"

	"github.com/lu4p/embed-encrypt/encryptedfs"
)

//go:generate go run ..

//encrypted:embed hello.txt
var hello string

//encrypted:embed gopher.png
var gopher []byte

//encrypted:embed gopher.png
//encrypted:embed hello.txt "another.txt" "with spaces .txt"
var multiplefiles encryptedfs.FS

//encrypted:embed *.txt
var glob encryptedfs.FS

func main() {
	c, err := fs.ReadFile(glob, "hello.txt")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(hello, len(gopher), multiplefiles)

	log.Println(string(c))
}
