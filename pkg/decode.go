package pkg

import (
	"fmt"
	"image/color"
	"io"
	"math"

	"github.com/OpenDiablo2/bitstream"
)

func (pl2 *PL2) Decode(rs io.ReadSeeker) (*PL2, error) {
	stream := bitstream.NewReader(rs)

	if err := pl2.decodeBasePalette(stream); err != nil {
		return nil, err
	}

	if err := pl2.decodeTransforms(stream); err != nil {
		return nil, err
	}

	if err := pl2.decodeTextColors(stream); err != nil {
		return nil, err
	}

	if err := pl2.decodeTextColorTransforms(stream); err != nil {
		return nil, err
	}

	return pl2, nil
}

func (pl2 *PL2) decodeColors(stream *bitstream.Reader, dst color.Palette, colorBytes int) error {
	const (
		rOff, gOff, bOff = 0, 1, 2 // rgb offsets in the returned bytes
	)

	for idx := range dst {
		rgba, err := stream.Next(colorBytes).Bytes().AsBytes()

		if err != nil {
			return fmt.Errorf("could not decode color, %w", err)
		}

		r, g, b, a := rgba[rOff], rgba[gOff], rgba[bOff], uint8(math.MaxUint8)

		dst[idx] = &color.RGBA{R: r, G: g, B: b, A: a}
	}

	return nil
}

func (pl2 *PL2) decodeBasePalette(stream *bitstream.Reader) error {
	pl2.BasePalette = make(color.Palette, numPaletteColors)

	return pl2.decodeColors(stream, pl2.BasePalette, 4)
}

func (pl2 *PL2) decodeTextColors(stream *bitstream.Reader) error {
	pl2.TextColors = make(color.Palette, numTextColors)

	return pl2.decodeColors(stream, pl2.TextColors, 3)
}

func (pl2 *PL2) decodeTransforms(stream *bitstream.Reader) (err error) {
	if err = pl2.decodeLightingTransforms(stream); err != nil {
		return err
	}

	if err = pl2.decodeBlendModeTransforms(stream); err != nil {
		return err
	}

	if err = pl2.decodeColorVariationTransforms(stream); err != nil {
		return err
	}

	if err = pl2.decodeOtherTransforms(stream); err != nil {
		return err
	}

	return nil
}

func (pl2 *PL2) decodeLightingTransforms(stream *bitstream.Reader) (err error) {
	pl2.LightLevelVariations = make([]Transform, lightLevelVariations)
	if err = pl2.decodeTransformMulti(stream, &pl2.LightLevelVariations); err != nil {
		return err
	}

	pl2.InvColorVariations = make([]Transform, invColorVariations)
	if err = pl2.decodeTransformMulti(stream, &pl2.InvColorVariations); err != nil {
		return err
	}

	if err = pl2.decodeTransformSingle(stream, &pl2.SelectedUnitShift); err != nil {
		return err
	}

	return nil
}

func (pl2 *PL2) decodeBlendModeTransforms(stream *bitstream.Reader) (err error) {
	pl2.AlphaBlend = make([][]Transform, alphaBlendCoarse)
	for blendIdx := range pl2.AlphaBlend {
		pl2.AlphaBlend[blendIdx] = make([]Transform, alphaBlendFine)
		if err = pl2.decodeTransformMulti(stream, &pl2.AlphaBlend[blendIdx]); err != nil {
			return err
		}
	}

	pl2.AdditiveBlend = make([]Transform, additiveBlends)
	if err = pl2.decodeTransformMulti(stream, &pl2.AdditiveBlend); err != nil {
		return err
	}

	pl2.MultiplicativeBlend = make([]Transform, multiplyBlends)

	return pl2.decodeTransformMulti(stream, &pl2.MultiplicativeBlend)
}

func (pl2 *PL2) decodeColorVariationTransforms(stream *bitstream.Reader) (err error) {
	pl2.HueVariations = make([]Transform, hueVariations)
	if err = pl2.decodeTransformMulti(stream, &pl2.HueVariations); err != nil {
		return err
	}

	if err = pl2.decodeTransformSingle(stream, &pl2.RedTones); err != nil {
		return err
	}

	if err = pl2.decodeTransformSingle(stream, &pl2.GreenTones); err != nil {
		return err
	}

	return pl2.decodeTransformSingle(stream, &pl2.BlueTones)
}

func (pl2 *PL2) decodeOtherTransforms(stream *bitstream.Reader) (err error) {
	pl2.UnknownVariations = make([]Transform, unknownVariations)
	if err = pl2.decodeTransformMulti(stream, &pl2.UnknownVariations); err != nil {
		return err
	}

	pl2.MaxComponentBlend = make([]Transform, maxComponentBlends)
	if err = pl2.decodeTransformMulti(stream, &pl2.MaxComponentBlend); err != nil {
		return err
	}

	return pl2.decodeTransformSingle(stream, &pl2.DarkenedColorShift)
}

func (pl2 *PL2) decodeTransformMulti(stream *bitstream.Reader, dst *[]Transform) error {
	for idx := range *dst {
		err := pl2.decodeTransformSingle(stream, &((*dst)[idx]))

		if err != nil {
			return err
		}
	}

	return nil
}

func (pl2 *PL2) decodeTransformSingle(stream *bitstream.Reader, dst *Transform) error {
	indices, err := stream.Next(numPaletteColors).Bytes().AsBytes()
	if err != nil {
		return fmt.Errorf("could not decode transform, %w", err)
	}

	for idx, paletteIndex := range indices {
		dst[idx] = paletteIndex
	}

	return nil
}

func (pl2 *PL2) decodeTextColorTransforms(stream *bitstream.Reader) (err error) {
	pl2.TextColorShifts = make([]Transform, textShifts)

	return pl2.decodeTransformMulti(stream, &pl2.TextColorShifts)
}