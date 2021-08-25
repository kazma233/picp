package main

import (
	_ "embed"
	"errors"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"picp/entity"
	"picp/utils"
)

var (
	//go:embed `resource/font/SourceHanSansCN-Heavy.otf`
	heavyFontBs []byte
	//go:embed `resource/font/SourceHanSansCN-Bold.otf`
	boldFontBs []byte
	//go:embed `resource/font/SourceHanSansCN-Regular.otf`
	regularFontBs []byte
	//go:embed `resource/font/SourceHanSansCN-Medium.otf`
	mediumFontBs []byte

	// 字体
	heavyFont   *sfnt.Font
	boldFont    *sfnt.Font
	regularFont *sfnt.Font
	mediumFont  *sfnt.Font
)

var (
	errInvalidFormat = errors.New("invalid color format")
)

func init() {
	// 字体
	// heavy
	_heavyFont, err := opentype.Parse(heavyFontBs)
	if err != nil {
		panic(err)
	}

	heavyFont = _heavyFont

	// bold
	_boldFont, err := opentype.Parse(boldFontBs)
	if err != nil {
		panic(err)
	}

	boldFont = _boldFont

	// regular
	_regularFont, err := opentype.Parse(regularFontBs)
	if err != nil {
		panic(err)
	}

	regularFont = _regularFont

	// medium
	_mediumFont, err := opentype.Parse(mediumFontBs)
	if err != nil {
		panic(err)
	}

	mediumFont = _mediumFont
}

func main() {
	g := gin.New()
	g.Use(gin.ErrorLogger())

	g.GET("/", func(context *gin.Context) {
		context.File("templates/index.html")
	})

	fg := g.Group("/file")
	fg.POST("", func(context *gin.Context) {
		f, err := context.FormFile("file")
		if err != nil {
			utils.LogX.Printf("解析文件错误: %v", err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		fn := xid.New().String()
		err = context.SaveUploadedFile(f, utils.TmpDir+fn)
		if err != nil {
			utils.LogX.Printf("保存文件错误: %v", err)
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		context.String(http.StatusOK, "%s", fn)
	})

	fg.GET("/:fileId", func(context *gin.Context) {
		fId := context.Param("fileId")
		f, err := os.OpenFile(utils.TmpDir+"/"+fId, os.O_RDONLY, 0755)
		if err != nil {
			utils.LogX.Printf("文件打开错误: %v", err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		fs, err := ioutil.ReadAll(f)
		if err != nil {
			utils.LogX.Printf("文件读取错误: %v", err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		context.Data(http.StatusOK, "image/jpeg", fs)
	})

	g.POST("/mark", func(context *gin.Context) {
		var addMask entity.AddMask
		if err := context.BindJSON(&addMask); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		origin, err := imaging.Open(utils.TmpDir + addMask.Origin)
		if err != nil {
			utils.LogX.Printf("打开原文件异常: %v", err)
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		for _, mask := range addMask.Infos {
			if mask.Type == "text" {
				face, err := utils.PreWordMask(entity.WordMaskPreInfo{
					Word: mask.Word,
					Font: getFont(mask.Font),
					Size: mask.Size,
					Dpi:  mask.Dpi,
				})
				if err != nil {
					utils.LogX.Printf("处理字体失败: %v", err)
					break
				}

				origin, err = utils.WriteWordMask(face, entity.WordMaskInfo{
					BgImg:      origin,
					Word:       mask.Word,
					ColorPoint: entity.ColorPoint{X: mask.X, Y: mask.Y, C: parseHexColorFast(mask.Color)},
				})
				if err != nil {
					utils.LogX.Printf("写入文字失败: %v", err)
					break
				}
			}

			if mask.Type == "img" {
				waterMask, err := imaging.Open(utils.TmpDir + mask.WaterMask)
				if err != nil {
					utils.LogX.Printf("打开水印文件异常：%v", err)
					break
				}

				origin = imaging.Overlay(origin, waterMask, image.Pt(mask.X, mask.Y), mask.Opacity)
			}

		}

		fn := xid.New().String() + ".jpg"
		err = imaging.Save(origin, utils.TmpDir+fn)
		if err != nil {
			utils.LogX.Printf("保存文件异常：%v", err)
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		context.String(http.StatusOK, "%v", fn)
	})

	g.Static("/js/", "./templates/js/")
	g.Static("/css/", "./templates/css/")

	log.Println(g.Run(":9000"))
}

func getFont(ff string) *opentype.Font {
	if ff == "Heavy" {
		return heavyFont
	}
	if ff == "Bold" {
		return boldFont
	}
	if ff == "Regular" {
		return regularFont
	}
	if ff == "Medium" {
		return mediumFont
	}

	return regularFont
}

func parseHexColorFast(s string) (c color.RGBA) {
	c.A = 0xff

	if s[0] != '#' {
		panic(parseHexColorFast)
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}

		panic(parseHexColorFast)
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		panic(parseHexColorFast)
	}

	return
}
