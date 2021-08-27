package main

import (
	"image"
	"io/ioutil"
	"net/http"
	"os"
	"picp/entity"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

func main() {
	logx := NewSugaredLogger()

	g := gin.Default()
	g.HandleMethodNotAllowed = true

	g.GET("/", func(context *gin.Context) {
		context.File("templates/index.html")
	})

	fg := g.Group("/file")
	fg.POST("", func(context *gin.Context) {
		f, err := context.FormFile("file")
		if err != nil {
			logx.Errorf("解析文件错误: %v", err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		fn := xid.New().String()
		err = context.SaveUploadedFile(f, TmpDir+fn)
		if err != nil {
			logx.Errorf("保存文件错误: %v", err)
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		context.String(http.StatusOK, "%s", fn)
	})

	fg.GET("/:fileId", func(context *gin.Context) {
		fId := context.Param("fileId")
		f, err := os.OpenFile(TmpDir+fId, os.O_RDONLY, 0755)
		if err != nil {
			logx.Errorf("文件打开错误: %v", err)
			context.AbortWithStatus(http.StatusBadRequest)
			return
		}

		fs, err := ioutil.ReadAll(f)
		if err != nil {
			logx.Errorf("文件读取错误: %v", err)
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
		logx.Infof("mark info: %v", addMask)

		origin, err := imaging.Open(TmpDir + addMask.Origin)
		if err != nil {
			logx.Errorf("打开原文件异常: %v", err)
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

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
				waterMask, err := imaging.Open(TmpDir + mask.WaterMask)
				if err != nil {
					logx.Errorf("打开水印文件异常：%v", err)
					break
				}

				origin = imaging.Overlay(origin, waterMask, image.Pt(mask.X, mask.Y), mask.Opacity)
			}

		}

		fn := xid.New().String() + ".jpg"
		err = imaging.Save(origin, TmpDir+fn)
		if err != nil {
			logx.Errorf("保存文件异常：%v", err)
			context.AbortWithStatus(http.StatusInternalServerError)
			return
		}

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
