package main

import (
	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
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

type (
	WordMaskPreInfo struct {
		Word string         // 文字水印的文字
		Font *opentype.Font // 字体
		Size float64        // 文字大小
		Dpi  float64        // DPI
	}

	WordMaskInfo struct {
		BgImg      image.Image // 背景图
		Word       string      // 文字
		ColorPoint             // 颜色 位置
	}

	WordMaskCenterInfo struct {
		BgImg image.Image // 背景图
		Word  string      // 文字
		C     color.Color
		Width int
		Y     int
		Font  *opentype.Font
		Size  float64
		Dpi   float64
	}
)

type (
	WriteColorInfo struct {
		BgImage      image.Image // 底图
		ColorBoxInfo ColorBox    // 长宽 位置 颜色
	}
)

// WriteColorMask 在图片上填满色块
func WriteColorMask(info WriteColorInfo) image.Image {
	return imaging.Paste(
		info.BgImage,
		imaging.New(info.ColorBoxInfo.Width, info.ColorBoxInfo.High, info.ColorBoxInfo.Color),
		info.ColorBoxInfo.Point,
	)
}

// PreWordMask 获取文字水印
func PreWordMask(info WordMaskPreInfo) (font.Face, error) {
	return opentype.NewFace(info.Font, &opentype.FaceOptions{
		Size:    info.Size,
		DPI:     info.Dpi,
		Hinting: font.HintingNone,
	})
}

// WriteWordMask 写文字水印
func WriteWordMask(face font.Face, info WordMaskInfo) (image.Image, error) {
	drawer := font.Drawer{Face: face}

	bgImg := info.BgImg
	dstImg := image.NewRGBA(bgImg.Bounds())
	drawer.Dst = dstImg
	drawer.Dot = fixed.P(info.X, info.Y)
	drawer.Src = image.NewUniform(info.C)

	drawer.DrawString(info.Word)

	return imaging.OverlayCenter(bgImg, dstImg, 1), nil
}

// WriteFontCenter 写出文字（中间）
func WriteFontCenter(info WordMaskCenterInfo) (image.Image, error) {
	fontFace, err := PreWordMask(
		WordMaskPreInfo{
			Word: info.Word,
			Font: info.Font,
			Size: info.Size,
			Dpi:  info.Dpi,
		},
	)

	if err != nil {
		return nil, err
	}

	drawer := font.Drawer{Face: fontFace}
	fSize := drawer.MeasureString(info.Word)

	w := (info.Width - fSize.Floor()) / 2
	wordMaskImage, err := WriteWordMask(
		fontFace,
		WordMaskInfo{
			BgImg: info.BgImg,
			Word:  info.Word,
			ColorPoint: ColorPoint{
				C: info.C,
				X: w,
				Y: info.Y,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return wordMaskImage, nil
}
