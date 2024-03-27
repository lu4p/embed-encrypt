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

//encrypted:embed bin *.txt
var glob encryptedfs.FS

func main() {
	log.SetFlags(log.Lshortfile)
	g, err := fs.ReadFile(glob, "bin/gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(hello, len(gopher), len(g))
	h, err := multiplefiles.ReadFile("hello.txt")
	log.Println(string(h), err)

	// all files in embed
	glob.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fi, _ := d.Info()
		if !fi.IsDir() {
			fmt.Println(fi.Mode(), fi.ModTime(), fi.Size(), path)
		}
		return nil
	})

	// files in root of embed
	log.Println(fs.ReadDir(glob, "."))

	// glob of go not like globstar
	log.Println(fs.Glob(glob, "**"))

	//GlobStar test
	log.Println(encryptedfs.GlobStar(glob, "."))

	wd, err := os.Getwd()
	if err == nil {
		// https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/xcopy
		// xcopy embed:\bin %CD%\copy\ /syd
		log.Println(encryptedfs.Xcopy(glob, "bin", wd, "copy"))
		// bin\gopher.png -> %CD%\u\bin\gopher.png

		// xcopy embed:\. %CD%\copy\ /syd
		log.Println(encryptedfs.Xcopy(glob, ".", wd, "copy"))
		// another.txt -> %CD%\u\another.txt
		// hello.txt -> %CD%\u\hello.txt
		// "with spaces .txt" -> %CD%\u\"with spaces .txt"
	}
}
