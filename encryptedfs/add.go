package encryptedfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ENC                 = ".enc"
	MODTIME             = /*version*/ 1 + /*sec*/ 8 + /*nsec*/ 4 + /*zone offset*/ 2
	NONCE               = 12
	timeBinaryVersionV2 = 2 // For LMT only
)

func (f FS) WalkDir(root string, fn fs.WalkDirFunc) error {
	return fs.WalkDir(f, root, fn)
}

func (f *fileInfo) String() string { return fs.FormatFileInfo(f) }

/*
Like xcopy embed:\src root\trg\ /syd

src - name of dir was embed. Root as "."

root - root dir for target

trg - target dir as `root/trg/“. If `a` and `b/c` was embed, and  root=`/tmp`

src="." trg="" then will be `/tmp/a` and `/tmp/b/c`

src="b" trg="" then will be `/tmp/b/c`

src="b" trg="d" then will be `/tmp/d/b/c`
*/
func Xcopy(bin fs.FS, src, root, trg string) (fns map[string]string, report string, err error) {
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
	dirs := append([]string{strings.ReplaceAll(root, `\`, "/")}, strings.Split(strings.ReplaceAll(trg, `\`, "/"), "/")...)
	write := func(unix string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		path := filepath.Join(append(dirs, strings.Split(unix, "/")[srcLen:]...)...)
		fns[strings.TrimPrefix(unix, src+"/")] = path
		eInfo, _ := fs.Stat(bin, unix)
		fInfo, err := os.Stat(path)
		ts := ""
		if err == nil {
			if fInfo.ModTime().Compare(eInfo.ModTime()) >= 0 { // xcopy /d fInfo.ModTime().After(eInfo.ModTime())
				return nil
			}
			ts = fmt.Sprint(fInfo.ModTime(), " ", fInfo.Size())
		}
		if d.IsDir() {
			if os.IsNotExist(err) {
				err = os.MkdirAll(path, DIRMODE)
				if err != nil {
					return err
				}
			}
			return err
		}
		bytes, err := fs.ReadFile(bin, unix)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, bytes, FiLEMODE)
		if err != nil {
			return err
		}
		fi, err := os.Stat(path)
		if err != nil {
			return err
		}
		s := fi.Size()
		l := int64(len(bytes))
		if l != s {
			err = fmt.Errorf("writing error to %s, expected %d, was recorded %d", path, l, s)
			return err
		}
		os.Chtimes(path, time.Now().Local(), eInfo.ModTime())
		report += fmt.Sprintln(eInfo.ModTime(), eInfo.Size(), unix, "->", ts, path)
		return nil
	}
	err = fs.WalkDir(bin, src, write)
	return
}

// like src/** case shopt -s globstar
func GlobStar(bin fs.FS, src string) (paths []string, err error) {
	list := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		paths = append(paths, path)
		return nil
	}
	err = fs.WalkDir(bin, src, list)
	return
}
