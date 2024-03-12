package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"

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

	glob.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fi, _ := d.Info()
		fmt.Println(fi.Mode(), fi.Size(), path)
		return nil
	})
	log.Println(glob.ReadDir("."))
	log.Println(fs.ReadDir(glob, "."))
	log.Println(encryptedfs.EmbedList(glob, "."))
	log.Println(fs.Stat(glob, "bin/gopher.png"))

	cwd, err := os.Getwd()
	if err == nil {
		log.Println(encryptedfs.UnloadEmbed(glob, "", cwd, "u", true))
		log.Println(encryptedfs.UnloadEmbed(glob, "bin", cwd, "u", true))
	}
}
