package main

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/abakum/embed-encrypt/encryptedfs"
)

// go:generate go run public/main.go key tool
//go:generate go run public/main.go

// go:generate go run .. key tool Priv
//go:generate go run ..

//encrypted:embed hello.txt
var hello string

//encrypted:embed image/gopher.png
var gopher []byte

//encrypted:embed hello.txt "another.txt" "with spaces .txt" image/gopher.png
var multiplefiles encryptedfs.FS

//encrypted:embed image *.txt
var glob encryptedfs.FS

func main() {
	log.SetFlags(log.Lshortfile)
	g, err := fs.ReadFile(glob, "image/gopher.png")
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
			log.Println(fi.Mode(), fi.ModTime(), fi.Size(), path)
		}
		return nil
	})

	// files in root of embed
	log.Println(fs.ReadDir(glob, "."))

	// glob of go not like globstar
	log.Println(fs.Glob(glob, "**"))

	//GlobStar test
	log.Println(encryptedfs.GlobStar(glob, "."))

	wd, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	wd = filepath.Dir(wd)
	if err == nil {
		// https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/xcopy
		// xcopy embed:\image %CD%\copy\ /syd
		log.Println(encryptedfs.Xcopy(glob, "image", wd, "copy"))
		// image\gopher.png -> %CD%\copy\image\gopher.png

		// xcopy embed:\. %CD%\copy\ /syd
		log.Println(encryptedfs.Xcopy(glob, ".", wd, "copy"))
		// another.txt -> %CD%\copy\another.txt
		// hello.txt -> %CD%\copy\hello.txt
		// "with spaces .txt" -> %CD%\copy\"with spaces .txt"
	}
}
