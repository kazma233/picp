package entity

import "image"

type WriteColorInfo struct {
	BgImage      image.Image // 底图
	ColorBoxInfo ColorBox    // 长宽 位置 颜色
}
