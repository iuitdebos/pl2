package pkg

import (
	"bytes"
	"image/color"
	"io"
	"math"

	color2 "github.com/lucasb-eyer/go-colorful"
)

const (
	numPaletteColors = 256
	numTextColors    = 13
)

const (
	lightLevelVariations = 32
	invColorVariations   = 16
	alphaBlendCoarse     = 3
	alphaBlendFine       = 256
	additiveBlends       = 256
	multiplyBlends       = 256
	hueVariations        = 111
	unknownVariations    = 14
	maxComponentBlends   = 256
	textShifts           = 13
	numOtherTransforms 	 = 5 // selected/darkened/r/g/b tones
	NumTransforms = 1 + // lightLevelVariations +
		invColorVariations +
		(alphaBlendCoarse * alphaBlendFine) +
		additiveBlends +
		multiplyBlends +
		hueVariations +
		unknownVariations +
		maxComponentBlends +
		textShifts +
		numOtherTransforms
)

// Transform represents a PL2 palette transform.
type Transform [numPaletteColors]uint8

func (t Transform) MakePaletteFromPalette(src color.Palette) color.Palette {
	dst := make(color.Palette, numPaletteColors)

	for idx, transformedIdx := range t {
		c := src[transformedIdx]
		cr, cg, cb, _ := c.RGBA()
		dst[idx] = color.RGBA{
			R: uint8(cr),
			G: uint8(cg),
			B: uint8(cb),
			A: math.MaxUint8,
		}
	}

	return dst
}

// PL2 represents a base palette, and different categories of "transforms" into that palette.
type PL2 struct {
	BasePalette color.Palette

	LightLevelVariations []Transform
	InvColorVariations   []Transform
	SelectedUnitShift   Transform
	AlphaBlend          [][]Transform
	AdditiveBlend       []Transform
	MultiplicativeBlend []Transform
	HueVariations       []Transform
	RedTones            Transform
	GreenTones          Transform
	BlueTones           Transform
	UnknownVariations   []Transform
	MaxComponentBlend   []Transform
	DarkenedColorShift  Transform

	TextColors      color.Palette
	TextColorShifts []Transform

	hslColorsBuffer []color2.Color
}

// FromBytes reads the bytes into a struct
func FromBytes(data []byte) (*PL2, error) {
	return (&PL2{}).Decode(bytes.NewReader(data))
}

func ToBytes(pl2 *PL2) ([]byte, error) {
	return EncodePalette(pl2.BasePalette)
}

// Decode the stream into a PL2
func Decode(rs io.ReadSeeker) (*PL2, error) {
	return (&PL2{}).Decode(rs)
}

// EncodePalette encodes the given palette as a PL2
func EncodePalette(p color.Palette) ([]byte, error) {
	pl2 := &PL2{}

	pl2.SetMainPalette(p)
	pl2.regenerate()

	b := bytes.NewBuffer(nil)
	err := pl2.Encode(b)

	return b.Bytes(), err
}

func (pl2 *PL2) SetMainPalette(src color.Palette) {
	dst := make(color.Palette, numPaletteColors)
	pl2.hslColorsBuffer = nil

	// ensure grayscale palette as default
	if len(src) < numPaletteColors {
		for idx := range dst {
			dst[idx] = color.RGBA{
				R: uint8(idx),
				G: uint8(idx),
				B: uint8(idx),
				A: math.MaxUint8,
			}
		}
	}

	copy(dst, src)

	pl2.BasePalette = dst
}

func (pl2 *PL2) SetTextPalette(src color.Palette) {
	dst := defaultTextColors()

	copy(dst, src)

	pl2.TextColors = dst
}
