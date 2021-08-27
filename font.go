package main

import (
	_ "embed"
	"log"

	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type FontType int

const (
	Bold FontType = iota // 0
	ExtraLight
	Heavy
	Light
	Medium
	Normal
	Regular
)

var (
	//go:embed `resource/font/SourceHanSansCN-Bold.otf`
	boldFontBs []byte
	//go:embed `resource/font/SourceHanSansCN-ExtraLight.otf`
	extraLightFontBs []byte
	//go:embed `resource/font/SourceHanSansCN-Heavy.otf`
	heavyFontBs []byte
	//go:embed `resource/font/SourceHanSansCN-Light.otf`
	lightFontBs []byte
	//go:embed `resource/font/SourceHanSansCN-Medium.otf`
	mediumFontBs []byte
	//go:embed `resource/font/SourceHanSansCN-Normal.otf`
	normalFontBs []byte
	//go:embed `resource/font/SourceHanSansCN-Regular.otf`
	regularFontBs []byte

	// 字体映射
	fontMap map[FontType]*sfnt.Font
)

func init() {
	fontMap = make(map[FontType]*sfnt.Font)

	// Bold
	_boldFont, err := opentype.Parse(boldFontBs)
	if err != nil {
		log.Printf("读取Bold字体失败: %v", err)
	} else {
		fontMap[Bold] = _boldFont
	}

	// ExtraLight
	_extraLight, err := opentype.Parse(extraLightFontBs)
	if err != nil {
		log.Printf("读取ExtraLight字体失败: %v", err)
	} else {
		fontMap[ExtraLight] = _extraLight
	}

	// Heavy
	_heavyFont, err := opentype.Parse(heavyFontBs)
	if err != nil {
		log.Printf("读取Heavy字体失败: %v", err)
	} else {
		fontMap[Heavy] = _heavyFont
	}

	// Light
	_lightFont, err := opentype.Parse(lightFontBs)
	if err != nil {
		log.Printf("读取Light字体失败: %v", err)
	} else {
		fontMap[Light] = _lightFont
	}

	// Medium
	_mediumFont, err := opentype.Parse(mediumFontBs)
	if err != nil {
		log.Printf("读取Medium字体失败: %v", err)
	} else {
		fontMap[Medium] = _mediumFont
	}

	// Normal
	_normalFont, err := opentype.Parse(normalFontBs)
	if err != nil {
		log.Printf("读取Normal字体失败: %v", err)
	} else {
		fontMap[Normal] = _normalFont
	}

	// Regular
	_regularFont, err := opentype.Parse(regularFontBs)
	if err != nil {
		log.Printf("读取Regular字体失败: %v", err)
	} else {
		fontMap[Regular] = _regularFont
	}
}

func getFont(ft FontType) *opentype.Font {
	return fontMap[ft]
}
