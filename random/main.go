package main

import (
	"io/fs"
	"log"

	"github.com/abakum/embed-encrypt/encryptedfs"
)

//go:generate go run github.com/abakum/embed-encrypt random

// try pass unexists file
//
// encrypted:embed unexists.txt
// var _ string

// https://pkg.go.dev/embed#hdr-Directives
// The difference is that ‘image/*’ embeds ‘image/.tempfile’ while ‘image’ does not.
//
//encrypted:embed image
var glob encryptedfs.FS

//encrypted:embed image/* image/dir
var glob2 encryptedfs.FS

//encrypted:embed image/* image/dir/*
var glob3 encryptedfs.FS

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
