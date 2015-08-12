package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jteeuwen/go-bindata"
)

const (
	Out     = "web/bindata.go"
	Prefix  = "web/assets"
	Package = "web"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "In debug mode, the files are read from disk")

	flag.Parse()

	config := &bindata.Config{
		Package:    Package,
		Output:     Out,
		Prefix:     Prefix,
		NoMemCopy:  false,
		NoCompress: false,
		Debug:      debug,
		Input: []bindata.InputConfig{
			bindata.InputConfig{
				Path: Prefix, Recursive: true,
			},
		},
	}
	err := bindata.Translate(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
