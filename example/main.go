package main

import (
	"io/fs"
	"log"

	"github.com/abakum/embed-encrypt/encryptedfs"
)

// go:generate go run public/main.go key tool
//go:generate go run public/main.go

// go:generate go run .. key tool Priv
//go:generate go run ..

//encrypted:embed hello.txt
var hello string

//encrypted:embed bin/gopher.png
var gopher []byte

//encrypted:embed hello.txt "another.txt" "with spaces .txt" bin/gopher.png
var multiplefiles encryptedfs.FS

//encrypted:embed bin/* *.txt
var glob encryptedfs.FS

func main() {
	g, err := fs.ReadFile(glob, "bin/gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(hello, len(gopher), len(g))
	h, err := multiplefiles.ReadFile("hello.txt")
	log.Println(string(h), err)
}
