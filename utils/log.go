package utils

import (
	"io"
	"log"
	"os"
)

var TmpDir string
var LogPath string
var LogX *log.Logger

func init() {
	name, err := os.MkdirTemp("", "kazma")
	if err != nil {
		panic(err)
	}
	log.Printf("tmp path is: %v", name)
	TmpDir = name + "/"
	LogPath = TmpDir + "log.log"

	logFile, err := os.OpenFile(LogPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)

	LogX = log.New(mw, "", log.LstdFlags)
}
