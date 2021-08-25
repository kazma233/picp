package entity

import (
	"image"
	"image/color"
)

type (
	Shape struct {
		Width int
		High  int
	}

	Box struct {
		Shape
		image.Point
	}

	ColorPoint struct {
		X int
		Y int
		C color.Color
	}

	ColorBox struct {
		Width int
		High  int
		image.Point
		color.Color
	}
)
