package encryptedfs

import (
	"crypto/aes"
	"crypto/cipher"
	"embed"
	"io"
	"io/fs"
	"strings"
	"time"
)

func remSuffix(name string) string {
	return strings.TrimSuffix(name, ENC)
}

var (
	_ fs.ReadDirFS  = FS{}
	_ fs.ReadFileFS = FS{}
)

// FS a wrapper for an encrypted embed.FS
type FS struct {
	underlying embed.FS
	key        []byte
}

func DecByte(encData []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	plaintext, err := aesgcm.Open(nil, encData[:12], encData[12:], nil)
	if err != nil {
		panic(err)
	}

	return plaintext
}

func DecString(encString string, key []byte) string {
	return string(DecByte([]byte(encString), key))
}

func InitFS(eFS embed.FS, key []byte) FS {
	return FS{
		underlying: eFS,
		key:        key,
	}
}

// Open opens the named file or dir for reading and returns it as an fs.File.
func (f FS) Open(name string) (fs.File, error) {
	file, err := f.underlying.Open(name) //dir
	if err != nil {
		file, err = f.underlying.Open(name + ENC) //file
	}
	if err != nil {
		return nil, err
	}

	info, _ := file.Stat()

	if info.IsDir() {
		return &openFile{underlying: file}, nil
	}

	block, err := aes.NewCipher(f.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	encData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, encData[:12], encData[12:], nil)
	if err != nil {
		return nil, err
	}

	return &openFile{
		underlying: file,
		decrypted:  plaintext,
	}, nil
}

func (f FS) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := f.underlying.ReadDir(name)
	if err != nil {
		return nil, err
	}

	newEntries := make([]fs.DirEntry, len(entries))

	for i, entry := range entries {
		file, err := f.Open(remSuffix(entry.Name()))
		if err != nil {
			return nil, err
		}
		defer file.Close()

		info, _ := file.Stat()

		newEntries[i] = info.(*fileInfo)
	}

	return newEntries, nil
}

func (f FS) ReadFile(name string) ([]byte, error) {
	file, err := f.Open(name)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	return io.ReadAll(file)
}

type openFile struct {
	underlying fs.File
	decrypted  []byte
	offset     int64
}

func (f *openFile) Read(b []byte) (int, error) {
	if f.offset >= int64(len(f.decrypted)) {
		return 0, io.EOF
	}

	if f.offset < 0 {
		info, _ := f.Stat() // cannot error

		return 0, &fs.PathError{Op: "read", Path: info.Name(), Err: fs.ErrInvalid}
	}
	n := copy(b, f.decrypted[f.offset:])
	f.offset += int64(n)
	return n, nil
}

func (f *openFile) Stat() (fs.FileInfo, error) {
	info, _ := f.underlying.Stat() // cannot error

	return &fileInfo{
		underlying:    info,
		decryptedSize: int64(len(f.decrypted)),
	}, nil
}

func (f *openFile) Close() error {
	f.decrypted = nil
	return f.underlying.Close()
}

var (
	_ fs.FileInfo = (*fileInfo)(nil)
	_ fs.DirEntry = (*fileInfo)(nil)
)

type fileInfo struct {
	underlying    fs.FileInfo
	decryptedSize int64
}

func (f *fileInfo) Name() string               { return remSuffix(f.underlying.Name()) }
func (f *fileInfo) Size() int64                { return f.decryptedSize }
func (f *fileInfo) Mode() fs.FileMode          { return f.underlying.Mode() }
func (f *fileInfo) ModTime() time.Time         { return f.underlying.ModTime() }
func (f *fileInfo) IsDir() bool                { return f.underlying.IsDir() }
func (f *fileInfo) Sys() interface{}           { return nil }
func (f *fileInfo) Type() fs.FileMode          { return f.underlying.Mode().Type() }
func (f *fileInfo) Info() (fs.FileInfo, error) { return f, nil }
