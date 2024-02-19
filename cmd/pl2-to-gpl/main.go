package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	pl2 "github.com/muitdebos/pl2/pkg"
	gpl "github.com/gravestench/gpl/pkg"
)

func main() {
	o := &options{}

	if parseOptions(o) {
		flag.Usage()
	}

	data, err := ioutil.ReadFile(*o.pl2Path)
	if err != nil {
		const fmtErr = "could not read file, %v"
		fmt.Print(fmt.Errorf(fmtErr, err))

		return
	}

	p, err := pl2.FromBytes(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	gplPath := *o.gplPath

	var f io.Writer

	f = log.Writer()

	if gplPath != "" {
		ff, err := os.Create(gplPath)
		if err != nil {
			log.Fatal(err)
		}

		f = ff

		close := func() {
			_ = ff.Close()
		}

		defer close()
	}

	gplPalette := gpl.FromPalette(p.BasePalette)

	if err := gplPalette.Encode("", f); err != nil {
		log.Fatal(err)
	}
}

type options struct {
	pl2Path   *string
	gplPath    *string
}

func parseOptions(o *options) (terminate bool) {
	o.pl2Path = flag.String("pl2", "", "input dcc file (required)")
	o.gplPath = flag.String("gpl", "", "the output gpl file")

	flag.Parse()

	return *o.pl2Path == ""
}

func fileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}
