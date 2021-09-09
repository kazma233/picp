package main

import (
	"image"
	"net/http"
	"picp/entity"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func main() {
	logx := NewSugaredLogger()
	ic := NewImageCache()

	g := gin.Default()
	g.HandleMethodNotAllowed = true

	g.GET("/", func(context *gin.Context) {
		context.File("templates/index.html")
	})

	fg := g.Group("/file")
	fg.POST("", func(context *gin.Context) {
		fh, err := context.FormFile("file")
		if err != nil {
			logx.Errorf("解析文件错误: %v", err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		f, err := fh.Open()
		if err != nil {
			logx.Errorf("打开文件错误: %v", err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		img, err := imaging.Decode(f)
		if err != nil {
			logx.Errorf("打开图片文件错误: %v", err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// save
		fn := xid.New().String()
		ic.SetImageMust(fn, img)

		context.String(http.StatusOK, "%s", fn)
	})

	fg.GET("/:fileId", func(context *gin.Context) {
		fId := context.Param("fileId")

		bf := ic.GetBufferMust(fId)
		context.Data(http.StatusOK, "image/png", bf.Bytes())
	})

	g.POST("/mark", func(context *gin.Context) {
		var addMask entity.AddMask
		if err := context.BindJSON(&addMask); err != nil {
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}
		logx.Infof("mark info: %v", addMask)

		origin := ic.GetImageMust(addMask.Origin)

		for _, mask := range addMask.Infos {
			if mask.Type == "text" {
				face, err := PreWordMask(WordMaskPreInfo{
					Word: mask.Word,
					Font: getFont(FontType(mask.Font)),
					Size: mask.Size,
					Dpi:  mask.Dpi,
				})
				if err != nil {
					logx.Errorf("处理字体失败: %v", err)
					break
				}

				origin, err = WriteWordMask(face, WordMaskInfo{
					BgImg:      origin,
					Word:       mask.Word,
					ColorPoint: ColorPoint{X: mask.X, Y: mask.Y, C: parseHexColorFast(mask.Color)},
				})
				if err != nil {
					logx.Errorf("写入文字失败: %v", err)
					break
				}
			}

			if mask.Type == "img" {
				waterMask := ic.GetImageMust(mask.WaterMask)
				origin = imaging.Overlay(origin, waterMask, image.Pt(mask.X, mask.Y), mask.Opacity)
			}
		}

		fn := xid.New().String()
		ic.SetImageMust(fn, origin)

		context.String(http.StatusOK, "%v", fn)
	})

	logG := g.Group("/log")
	logG.GET("/high", func(context *gin.Context) {
		context.File(GetHighLogPath())
	})

	logG.GET("/low", func(context *gin.Context) {
		context.File(GetLowLogPath())
	})

	g.Static("/js/", "./templates/js/")
	g.Static("/css/", "./templates/css/")

	logx.Error(g.Run(":9000"))
}
