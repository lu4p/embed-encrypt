package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/abakum/embed-encrypt/encryptedfs"
)

//go:generate go run github.com/abakum/embed-encrypt random

// try pass unexists file
//
//encrypted:embed unexists.txt
// var _ string

// bin like bin/* if bin is dir
//
//encrypted:embed bin hello.txt
var glob encryptedfs.FS

func main() {
	log.SetFlags(log.Lshortfile)
	g, err := fs.ReadFile(glob, "bin/gopher.png")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(g))

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

	wd, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	wd = filepath.Dir(wd)
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
