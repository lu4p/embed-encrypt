package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/abakum/embed-encrypt/encryptedfs"
)

// try pass unexists file
//
// go:embed unexists.txt
// var _ string

// https://pkg.go.dev/embed#hdr-Directives
// The difference is that ‘image/*’ embeds ‘image/.tempfile’ while ‘image’ does not
// and ‘image’ is recursive  but ‘image/*’ does not.
//
//go:embed image
var glob embed.FS

//go:embed image/* image/dir
var glob2 embed.FS

//go:embed image/* image/dir/*
var glob3 embed.FS

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println(glob.ReadDir("image"))
	log.Println(glob.ReadDir("image/dir"))
	log.Println(encryptedfs.GlobStar(glob, "."))
	log.Println()
	log.Println(fs.ReadDir(glob2, "image"))
	log.Println(fs.ReadDir(glob2, "image/dir"))
	log.Println(encryptedfs.GlobStar(glob2, "."))
	log.Println()
	log.Println(fs.ReadDir(glob3, "image"))
	log.Println(fs.ReadDir(glob3, "image/dir"))
	log.Println(encryptedfs.GlobStar(glob3, "."))
}
