package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	name, directives, err := findDirectives()
	if err != nil {
		log.Fatal(err)
	}

	err = directives2Files(directives)
	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range directives {
		err := encryptFiles(dir.files)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = generateCode(name, directives)
	if err != nil {
		log.Fatal(err)
	}
}

func encryptFile(name string, key []byte) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}

	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	encData := append(nonce, gcm.Seal(nil, nonce, content, nil)...)

	info, err := f.Stat()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(name+".enc", encData, info.Mode())
}

const keyname = ".fskey"

func encryptFiles(filenames []string) error {
	key, err := ioutil.ReadFile(keyname)
	if err != nil {
		key = make([]byte, 16)
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			return fmt.Errorf("key couldn't be generated: %w", err)
		}

		if err := ioutil.WriteFile(keyname, key, 0666); err != nil {
			return fmt.Errorf("key couldn't be written: %w", err)
		}
	}

	for _, name := range filenames {
		if err := encryptFile(name, key); err != nil {
			return err
		}
	}

	return nil
}

type directive struct {
	identifier string
	typ        string
	patterns   []string
	files      []string
}

func findDirectives() (string, []directive, error) {
	fset := token.NewFileSet()

	filter := func(info os.FileInfo) bool {
		if strings.HasSuffix(info.Name(), "_test.go") {
			return false
		}
		return !strings.HasPrefix(info.Name(), "encrypted")
	}

	pkgs, err := parser.ParseDir(fset, ".", filter, parser.ParseComments)
	if err != nil {
		return "", nil, err
	}

	var directives []directive

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			visit := func(node ast.Node) bool {
				if node == nil {
					return true
				}

				gendecl, ok1 := node.(*ast.GenDecl)
				if !(ok1) {
					return true
				}

				if gendecl.Tok != token.VAR {
					return true
				}

				if len(gendecl.Specs) != 1 {
					return true
				}

				if gendecl.Doc == nil {
					return true
				}

				var dir directive

				hasComment := false
				for _, c := range gendecl.Doc.List {
					if !strings.HasPrefix(c.Text, "//encrypted:embed") {
						continue
					}

					args := strings.TrimPrefix(c.Text, "//encrypted:embed")

					fields, err := parseGoEmbed(args, token.Position{})
					if err != nil {
						panic(err)
					}

					dir.patterns = append(dir.patterns, fields...)

					hasComment = true
				}

				if !hasComment {
					return true
				}

				spec, ok := gendecl.Specs[0].(*ast.ValueSpec)
				if !ok {
					return true
				}

				dir.typ = fmt.Sprint(spec.Type)
				if strings.Contains(dir.typ, "encryptedfs") {
					dir.typ = "embed.FS"
				} else if strings.Contains(dir.typ, "byte") {
					dir.typ = "[]byte"
				}

				if len(spec.Names) != 1 {
					return true
				}
				dir.identifier = spec.Names[0].Name

				directives = append(directives, dir)

				return true
			}
			ast.Inspect(file, visit)
		}

		return pkg.Name, directives, nil
	}

	return "nil", nil, nil
}

func directives2Files(directives []directive) error {
	for i, d := range directives {
		for _, p := range d.patterns {
			if _, err := os.Stat(p); err == nil {
				directives[i].files = append(directives[i].files, p)
				continue
			}

			matches, err := filepath.Glob(p)
			if err != nil {
				return err
			}

			directives[i].files = append(directives[i].files, matches...)
		}
	}

	return nil
}
