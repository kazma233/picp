package utils

import (
	"image"
	"picp/entity"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)



// WriteColorMask 在图片上填满色块
func WriteColorMask(info entity.WriteColorInfo) image.Image {
	return imaging.Paste(
		info.BgImage,
		imaging.New(info.ColorBoxInfo.Width, info.ColorBoxInfo.High, info.ColorBoxInfo.Color),
		info.ColorBoxInfo.Point,
	)
}

// PreWordMask 获取文字水印
func PreWordMask(info entity.WordMaskPreInfo) (font.Face, error) {
	return opentype.NewFace(info.Font, &opentype.FaceOptions{
		Size:    info.Size,
		DPI:     info.Dpi,
		Hinting: font.HintingNone,
	})
}

// WriteWordMask 写文字水印
func WriteWordMask(face font.Face, info entity.WordMaskInfo) (image.Image, error) {
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
func WriteFontCenter(info entity.WordMaskCenterInfo) (image.Image, error) {
	fontFace, err := PreWordMask(
		entity.WordMaskPreInfo{
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
		entity.WordMaskInfo{
			BgImg: info.BgImg,
			Word:  info.Word,
			ColorPoint: entity.ColorPoint{
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
