package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/gravestench/pl2/pkg"
)

type options struct {
	pl2     *string
	pngPath *string
}

func parseOptions(o *options) (terminate bool) {
	o.pl2 = flag.String("pl2", "", "input pkg file (required)")
	o.pngPath = flag.String("png", "", "path to output png file (optional)")

	flag.Parse()

	return *o.pl2 == ""
}

func main() {
	o := &options{}

	if parseOptions(o) {
		flag.Usage()
	}

	data, err := ioutil.ReadFile(*o.pl2)
	if err != nil {
		const fmtErr = "could not read file, %v"
		fmt.Print(fmt.Errorf(fmtErr, err))

		return
	}

	pl2, err := pkg.FromBytes(data)
	if err != nil {
		return
	}

	if *o.pngPath != "" {
		img := makeImage(pl2)

		err = writeImage(*o.pngPath, img)
		if err != nil {
			fmt.Errorf("problem writing image, %w", err)
		}
	}
}

func makeImage(pl2 *pkg.PL2) image.Image {
	const (
		width = 256
	)

	mainTransforms := getMainTransforms(pl2)
	textTransforms := getTextTransforms(pl2)
	allTransforms := append(mainTransforms, textTransforms...)

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: len(allTransforms)},
	})

	for row := range allTransforms {
		writeTransformAsRowInImage(row, &allTransforms[row], pl2.BasePalette, img)
	}

	return img
}

func getMainTransforms(p *pkg.PL2) []pkg.Transform {
	transforms := make([]pkg.Transform, 0)

	baseTransform := pkg.Transform{}
	for idx := range baseTransform {
		baseTransform[idx] = uint8(idx)
	}

	addTransform := func(t pkg.Transform) {
		transforms = append(transforms, t)
	}

	addTransforms := func(transforms []pkg.Transform) {
		for idx := range transforms {
			addTransform(transforms[idx])
		}
	}

	addTransform(baseTransform)
	addTransforms(p.LightLevelVariations)
	addTransforms(p.InvColorVariations)

	addTransform(p.SelectedUnitShift)

	for idx := range p.AlphaBlend {
		addTransforms(p.AlphaBlend[idx])
	}

	addTransforms(p.AdditiveBlend)
	addTransforms(p.MultiplicativeBlend)
	addTransforms(p.HueVariations)

	addTransform(p.RedTones)
	addTransform(p.GreenTones)
	addTransform(p.BlueTones)

	addTransforms(p.MaxComponentBlend)
	addTransform(p.DarkenedColorShift)

	return transforms
}

func getTextTransforms(p *pkg.PL2) []pkg.Transform {
	transforms := make([]pkg.Transform, 0)

	baseTransform := pkg.Transform{}
	for idx := range baseTransform {
		baseTransform[idx] = uint8(idx)
	}

	addTransform := func(t pkg.Transform) {
		transforms = append(transforms, t)
	}

	addTransforms := func(transforms []pkg.Transform) {
		for idx := range transforms {
			addTransform(transforms[idx])
		}
	}

	addTransform(baseTransform)
	addTransforms(p.TextColorShifts)

	return transforms
}

func writeTransformAsRowInImage(y int, transform *pkg.Transform, pal color.Palette, img *image.RGBA) {
	for x, palIdx := range transform {
		if x > img.Bounds().Dx() {
			continue
		}

		if y > img.Bounds().Dy() {
			continue
		}

		var c color.Color
		c = color.Black

		if int(palIdx) < len(pal) {
			c = pal[palIdx]
		}

		r, g, b, a := c.RGBA()

		rgba := color.RGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: uint8(a),
		}

		img.Set(x, y, rgba)
	}
}

func writeImage(outPath string, img image.Image) error {
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}

	if err = png.Encode(f, img); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	return nil
}