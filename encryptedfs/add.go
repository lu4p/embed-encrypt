package encryptedfs

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const ENC = ".enc"

// `fs.WalkDir(encryptedfs` do not work replace it with `encryptedfs.WalkDir(`
func (f FS) WalkDir(root string, fn fs.WalkDirFunc) error {
	nfn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			path = remSuffix(path)
		}
		file, err := f.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		info, err := file.Stat()
		return fn(path, info.(*fileInfo), err)
	}
	return fs.WalkDir(f.underlying, root, nfn)
}

func (f *fileInfo) String() string { return fs.FormatFileInfo(f) }

/*
copy from embed

src - name of dir was embed. Root as "."

root - root dir for target

trg - target dir as `root/trg“. If `a` and `b/c` was ebbed, and  root=`/tmp`

src="." trg="" then will be `/tmp/a` `/tmp/b/c`

src="b" trg="" then will be `/tmp/b/c`

src="b" trg="d" then will be `/tmp/d/b/c`

keep == true if not exist then write

keep == false it will be replaced if it differs from the embed
*/
func UnloadEmbed(bin any, src, root, trg string, keep bool) (fns map[string]string, err error) {
	var fsys fs.FS
	const (
		FiLEMODE = 0644
		DIRMODE  = 0755
	)
	fns = make(map[string]string)
	if src == "" {
		src = "."
	}
	src = strings.ReplaceAll(src, `\`, "/")
	srcLen := strings.Count(src, "/")
	trg = strings.ReplaceAll(trg, `\`, "/")
	dirs := append([]string{root}, strings.Split(trg, "/")...)
	write := func(unix string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		win := filepath.Join(append(dirs, strings.Split(unix, "/")[srcLen:]...)...)
		fns[strings.TrimPrefix(unix, src+"/")] = win
		if d.IsDir() {
			_, err = os.Stat(win)
			if os.IsNotExist(err) {
				err = os.MkdirAll(win, DIRMODE)
			}
			return err
		}
		bytes, err := fs.ReadFile(fsys, unix)
		if err != nil {
			return err
		}
		var size int64
		fi, err := os.Stat(win)
		if err == nil {
			size = fi.Size()
			if int64(len(bytes)) == size || keep {
				return nil
			}
		}
		log.Println(win, len(bytes), "->", size)
		return os.WriteFile(win, bytes, FiLEMODE)
	}
	switch efs := bin.(type) {
	case embed.FS:
		fsys = efs
		err = fs.WalkDir(efs, src, write)
	case FS:
		fsys = efs
		err = efs.WalkDir(src, write)
	}
	return
}

// template usage WalkDir for UnloadEmbed
func EmbedList(bin any, src string) (paths []string, err error) {
	list := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		paths = append(paths, path)
		return nil
	}
	switch efs := bin.(type) {
	case embed.FS:
		err = fs.WalkDir(efs, src, list)
	case FS:
		err = efs.WalkDir(src, list)
	}
	return
}
