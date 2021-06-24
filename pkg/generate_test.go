package pkg

import (
	"bytes"
	"image/color"
	"testing"
)

func TestPL2_regenerate(t *testing.T) {
	type fields struct {
		BasePalette          color.Palette
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{"default grayscale", fields{nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl2 := &PL2{
				BasePalette:          tt.fields.BasePalette,
			}

			b := bytes.NewBuffer(nil)

			pl2.regenerate()

			err := pl2.Encode(b)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
