package main

import (
	"bytes"
	"github.com/allegro/bigcache"
	"golang.org/x/image/bmp"
	"image"
	"time"
)

type ImageCache struct {
	bc *bigcache.BigCache
}

func NewImageCache() ImageCache {
	bc, err := bigcache.NewBigCache(bigcache.Config{
		Shards:             1024,
		LifeWindow:         24 * time.Hour,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1024 * 60,
		MaxEntrySize:       100,
		Verbose:            true,
		Hasher:             nil,
		HardMaxCacheSize:   0,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
		Logger:             nil,
	})

	if err != nil {
		panic(err)
	}

	return ImageCache{bc: bc}
}

func (ic ImageCache) SetImageMust(key string, img image.Image) {
	bf := bytes.NewBuffer(nil)
	err := bmp.Encode(bf, img)
	if err != nil {
		panic(err)
	}

	ic.SetMust(key, bf.Bytes())
}

func (ic ImageCache) GetImageMust(key string) image.Image {
	bs, err := ic.bc.Get(key)
	if err != nil {
		panic(err)
	}

	bf := bytes.NewBuffer(bs)
	img, err := bmp.Decode(bf)
	if err != nil {
		panic(err)
	}

	return img
}

func (ic ImageCache) SetMust(key string, bs []byte) {
	err := ic.bc.Set(key, bs)
	if err != nil {
		panic(err)
	}
}

func (ic ImageCache) GetMust(key string) []byte {
	bs, err := ic.bc.Get(key)
	if err != nil {
		panic(err)
	}

	return bs
}

func (ic ImageCache) GetBufferMust(key string) *bytes.Buffer {
	return bytes.NewBuffer(ic.GetMust(key))
}
