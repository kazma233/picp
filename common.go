package main

import (
	"log"
	"os"
	"path/filepath"
)

var TmpDir string

func init() {
	name, err := os.MkdirTemp("", "picp_")
	if err != nil {
		panic(err)
	}

	log.Printf("tmp path is: %v", name)

	TmpDir = name + string(filepath.Separator)
}
