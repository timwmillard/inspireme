package inspireme

import "image/color"

// Alignment positions
const (
	LeftAlign = iota
	CenterAlign
	RightAlign
)

// Styles is a way to style your quote
type Styles struct {
	TextColor       color.RGBA
	Padding         int
	HorizontalAlign int
	VertialAlign    int
}

// DefaultStyles return the default styles
func DefaultStyles() Styles {
	return Styles{
		TextColor:       color.RGBA{255, 255, 255, 1},
		Padding:         0,
		HorizontalAlign: CenterAlign,
		VertialAlign:    CenterAlign,
	}
}
