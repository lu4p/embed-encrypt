# embed-encrypt 
This program wraps the newly in Go 1.16 `embed` package, to save files aes encrypted in the binary. 

The goal is to be as easy to use as the `embed` package, while providing 

**Note:** You will need the go 1.19 or newer to use this.

### If you have an older go version

To install the Go 1.19 just run:
```
go get golang.org/dl/go1.19
go1.19 download
```

## How to Install
```bash
go1.19 get github.com/abakum/embed-encrypt
```

## Usage

Replace your `//go:embed file.txt` directive  with `//encrypted:embed file.txt`,
the syntax is the same as the embed directives read the embed docs [here](https://pkg.go.dev/embed?utm_source=gopls#hdr-Directives). 
Replace your `embed.FS` type  with `encryptedfs.FS`.

After you added the comment directives to some variables run:
```bash
embed-encrypt
```
or if you haven't added `GOBIN` to your `PATH`

```bash
go run github.com/abakum/embed-encrypt
```

This generates an aes encrypted version for all embedded files, 
`<filename>.enc` and generates a `encrypted_fs.go` which includes these files via `//go:embed`, these files are automatically decrypted at runtime.

## Supported
Multiple directives, for a single variable:
```go
//encrypted:embed bin/gopher.png
//encrypted:embed hello.txt
var multipleDirectives encryptedfs.FS
```

Multiple files, for a single variable:
```go
//encrypted:embed bin/gopher.png
//encrypted:embed hello.txt "another.txt" "with spaces .txt"
var multipleFiles encryptedfs.FS
```

Glob patterns:
```go
//encrypted:embed bin/* *.txt
var glob encryptedfs.FS
```

Strings variables:
```go
//encrypted:embed hello.txt
var hello string
```

Byte slices:
```go
//encrypted:embed bin/gopher.png
var gopher []byte
```


## Caveats
- locally scoped embed variables are not supported, in go1.19

```go
func main() {
	//go:embed test.txt
    var glob encryptedfs.FS
}
```

- Block variables are currently not supported, they would be relatively easy to implement:
```go
var (
	//go:embed hello.txt
	hello string
	//go:embed gopher.png
	gopher []byte
)
```

Пять `потому, что` в ответ на вопрос `почему я использую этот форк encrypted:embed вместо go:embed`
 - даже без секретного ключа шифрования полезно шифровать исполняемые файлы в `embed` чтоб не беспокоить антивирусы и площадки где запрещена дистрибуция чужих бинарников
 - даже без секретного ключа шифрования этот форк сохраняет и восстанавливает даты `embed` файлов
 - этот форк может не включать ключ шифрования в `embed`
 - в этом форке можно использовать `WalkDir`, `Xcopy`, `GlobStar`
 - в этом форке ключ шифрования 
 ```go
 var key []byte
 ```
 может иметь имя отличное от `key`