package tool

import (
	"crypto/sha512"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
)

func Priv(keyEnc, libEnc string) (key []byte) {
	const ENC = ".enc"

	hash := sha512.Sum512_256([]byte("Priv"))
	if _, generate := os.LookupEnv("GOLINE"); generate {
		if keyEnc == "" {
			keyEnc = "key"
		}
		keyEnc += ENC
		if err := os.WriteFile(keyEnc, hash[:], 0600); err != nil {
			fmt.Printf("key couldn't be written: %v", err)
			return
		}

		if lib := getPackage(); lib != "" {
			if libEnc == "" {
				libEnc = path.Base(lib)
			}
			libEnc += ENC
			os.WriteFile(libEnc, []byte(lib), 0600)
		}
	}
	return hash[:]
}

func getPackage() string {
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	if parts[pl-2][0] == '(' {
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return packageName
}
