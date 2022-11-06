package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/shogo82148/base45"
)

var flagDecode bool

func init() {
	flag.BoolVar(&flagDecode, "decode", false, "decode base45 instead of encoding")
}

func main() {
	flag.Parse()

	if flagDecode {
		if err := decode(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else {
		if err := encode(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

func decode() error {
	dec := base45.NewDecoder(os.Stdin)
	_, err := io.Copy(os.Stdout, dec)
	if err != nil {
		return err
	}
	return nil
}

func encode() error {
	enc := base45.NewEncoder(os.Stdout)
	_, err := io.Copy(enc, os.Stdin)
	if err != nil {
		return err
	}
	if err := enc.Close(); err != nil {
		return err
	}
	return nil
}