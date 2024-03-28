package main

import (
	"flag"
	"public/tool"
)

func main() {
	flag.Parse()
	tool.Priv(flag.Arg(0), flag.Arg(1))
}
