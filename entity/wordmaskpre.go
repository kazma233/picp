package entity

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type WordMaskPreInfo struct {
	Word string         // 文字水印的文字
	Font *opentype.Font // 字体
	Size float64        // 文字大小
	Dpi  float64        // DPI
}

type WordMaskPreResult struct {
	Draw font.Drawer // 绘制文字的对象
	X    int
	Y    int
}
