# embed-encrypt 
This program wraps the newly in Go 1.16 `embed` package, to save files aes encrypted in the binary. 

The goal is to be as easy to use as the `embed` package, while providing 

**Note:** You will need the go 1.16 beta1 or newer to use this.

### If you have an older go version

To install the 1.16. beta just run:
```
go get golang.org/dl/go1.16beta1
go1.16beta1 download
```

## How to Install
```bash
go1.16beta1 get github.com/lu4p/embed-encrypt
```

## Usage

Replace your `//go:embed file.txt` directive  with `//encrypted:embed file.txt`,
the syntax is the same as the embed directives read the embed docs [here](https://pkg.go.dev/embed?utm_source=gopls#hdr-Directives). 
If you find that the syntax is not 100% compatible to `go:embed` please open an issue.

After you added the comment directives to some variables run:
```bash
embed-encrypt
```
or if you haven't added `GOBIN` to your `PATH`

```bash
go run github.com/lu4p/embed-encrypt
```

This generates an aes encrypted version for all embedded files, 
`<filename>.enc` and generates a `encrypted_fs.go` which includes these files via `//go:embed`, these files are automatically decrypted at runtime.

## Supported
Multiple directives, for a single variable:
```go
//encrypted:embed gopher.png
//encrypted:embed hello.txt
var multipleDirectives encryptedfs.FS
```

Multiple files, for a single variable:
```go
//encrypted:embed gopher.png
//encrypted:embed hello.txt "another.txt" "with spaces .txt"
var multipleFiles encryptedfs.FS
```

Glob patterns:
```go
//encrypted:embed *.txt
var glob encryptedfs.FS
```

Strings variables:
```go
//encrypted:embed hello.txt
var hello string
```

Byte slices:
```go
//encrypted:embed gopher.png
var gopher []byte
```


## Caveats
- locally scoped embed variables are not supported, because a seperate file is generated, which needs to manipulate the variables

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
- no tests (will be added shortly)
