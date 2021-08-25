package entity

import (
	"image"
	"image/color"

	"golang.org/x/image/font/opentype"
)

type WordMaskInfo struct {
	BgImg      image.Image // 背景图
	Word       string      // 文字
	ColorPoint             // 颜色 位置
}

type WordMaskCenterInfo struct {
	BgImg image.Image // 背景图
	Word  string      // 文字
	C     color.Color
	Width int
	Y     int
	Font  *opentype.Font
	Size  float64
	Dpi   float64
}
