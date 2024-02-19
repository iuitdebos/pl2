package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"

	pl2 "github.com/muitdebos/pl2/pkg"

	gpl "github.com/gravestench/gpl/pkg"
)


type options struct {
	gpl       *string
	out    *string
	outPrefix *string
}

func parseOptions(o *options) (terminate bool) {
	o.gpl = flag.String("gpl", "", "input dcc file (required)")
	o.out = flag.String("pl2", "./Pal.pl2", "the output directory (required)")

	flag.Parse()

	return *o.gpl == "" || *o.out == ""
}


func main() {
	o := &options{}

	if parseOptions(o) {
		flag.Usage()
	}

	data, err := ioutil.ReadFile(*o.gpl)
	if err != nil {
		const fmtErr = "could not read file, %v"
		fmt.Print(fmt.Errorf(fmtErr, err))

		return
	}

	gplPalette, err := gpl.Decode(bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return
	}

	pl2Bytes, err := pl2.EncodePalette(color.Palette(*gplPalette))
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Create(*o.out)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write(pl2Bytes); err != nil {
		_ = f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}