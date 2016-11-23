package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	prand "math/rand"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ek(err error) {
	if err != nil {
		log.Print(err)
	}
}

func randn(a, b int) int {
	z := new(big.Int)
	z.SetString(fmt.Sprint(b-a), 10)
	r, err := rand.Int(rand.Reader, z)

	c := prand.Intn(b - a)
	if err == nil {
		c = int(r.Uint64())
	}

	return a + c
}

func blitText(font *sdlttf.Font, x, y int, c sdl.Color, text string) {
	r, err := font.RenderUTF8BlendedEx(surface, text, c)
	ck(err)

	p, err := texture.Lock(nil)
	ck(err)

	err = surface.Lock()
	ck(err)

	s := surface.Pixels()
	for i := 0; i < len(p); i += 4 {
		p[i] = s[i+2]
		p[i+1] = s[i]
		p[i+2] = s[i+1]
		p[i+3] = s[i+3]
	}

	surface.Unlock()
	texture.Unlock()

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	screen.Copy(texture, &sdl.Rect{0, 0, r.W, r.H}, &sdl.Rect{int32(x), int32(y), r.W, r.H})
}