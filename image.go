package main

import (
	"log"
	"path/filepath"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

type Image struct {
	*sdl.Texture
	w, h int
}

var (
	images = make(map[string]*Image)
)

func loadImage(name string) *Image {
	log.SetPrefix("image: ")
	filename := filepath.Join(conf.assets, name)

	if m, found := images[filename]; found {
		return m
	}

	texture, err := sdlimage.LoadTextureFile(screen.Renderer, filename)
	ck(err)

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	_, _, w, h, _ := texture.Query()

	m := &Image{
		Texture: texture,
		w:       w,
		h:       h,
	}
	images[filename] = m
	return m
}

func (m *Image) Blit(x, y int) {
	screen.Copy(m.Texture, nil, &sdl.Rect{int32(x), int32(y), int32(m.w), int32(m.h)})
}

type fontKey struct {
	filename string
	ptsize   int
}

var (
	fonts = make(map[fontKey]*sdlttf.Font)
)

func loadFont(name string, ptsize int) *sdlttf.Font {
	log.SetPrefix("font: ")

	filename := filepath.Join(conf.assets, name)
	key := fontKey{filename, ptsize}
	if font, found := fonts[key]; found {
		return font
	}

	font, err := sdlttf.OpenFont(filename, ptsize)
	ck(err)

	fonts[key] = font
	return font
}
