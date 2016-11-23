package main

import (
	"fmt"
	"os"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlgfx"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

const (
	maxScore = 999999
)

type Game struct {
	bsound    *sdlmixer.Chunk
	rsound    *sdlmixer.Chunk
	beepsound *sdlmixer.Chunk
	bgsound   *sdlmixer.Music

	background  *Image
	rlayer      *Image
	windowLayer *Image
	images      [8]*Image

	digiFont   *sdlttf.Font
	creditFont *sdlttf.Font
	font       *sdlttf.Font

	menu    string
	wins    [5]int
	show    [9]int
	showOld [9]int
	keys    bool
	mut     bool
	lastwin int
	credit  int
	bet     int
}

func newGame() *Game {
	g := &Game{
		bsound:    loadSound("sounds/CLICK10A.WAV"),
		rsound:    loadSound("sounds/film_projector.wav"),
		beepsound: loadSound("sounds/beep.wav"),
		bgsound:   loadMusic("sounds/background001.wav"),

		digiFont:   loadFont("DIGITAL2.ttf", 24),
		font:       loadFont("LiberationSans-Regular.ttf", 15),
		creditFont: loadFont("LiberationSans-Regular.ttf", 55),

		background:  loadImage("img/bg.png"),
		rlayer:      loadImage("img/rlayer.png"),
		windowLayer: loadImage("img/windowlayer.png"),
	}

	for i := range g.images {
		g.images[i] = loadImage(fmt.Sprintf("img/%d.png", i+1))
	}

	return g
}

func (g *Game) reset() {
	g.mut = false
	g.keys = true
	g.credit = 20
	g.bet = 1
	g.lastwin = 0
	for i := range g.wins {
		g.wins[i] = 0
	}
	for i := range g.show {
		g.show[i] = 8
	}
	playMusic(g.bgsound)
}

func (g *Game) Run() {
	defer func() {
		recover()
	}()

	g.reset()
	for {
		screen.SetDrawColor(sdlcolor.Black)
		screen.Clear()
		g.background.Blit(0, 0)

		if g.event() {
			break
		}

		if g.draw() {
			break
		}
	}
}

func (g *Game) event() bool {
	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}
		switch ev := ev.(type) {
		case sdl.QuitEvent:
			os.Exit(0)
		case sdl.KeyDownEvent:
			playSound(g.bsound)
			if (ev.Sym == sdl.K_LEFT || ev.Sym == sdl.K_RIGHT) && g.keys {
				if g.credit > 0 {
					if g.credit-g.bet < 0 {
						g.bet = g.credit
					}
					if !conf.invincible {
						g.credit -= g.bet
					}
					g.randi()
					g.check()
					g.roll()
					g.background.Blit(0, 0)
					g.drawl()
					g.winner()
				} else if g.credit == 0 && g.bet == 0 {
					stopMusic()
					menu.Reset()
					state = menu.Run
					return true
				}
			}

			if g.credit > 0 {
				if ev.Sym == sdl.K_UP && g.keys {
					if g.credit-g.bet-1 >= 0 {
						g.bet++
					} else {
						g.bet = 1
					}

					if g.bet >= 11 {
						g.bet = 1
					}
				} else if ev.Sym == sdl.K_DOWN && g.keys {
					if g.bet--; g.bet <= 0 {
						g.bet = 10
					}
				}
			} else {
				g.bet = 0
			}

			if ev.Sym == sdl.K_F1 {
				if g.keys {
					g.menu = "h"
				} else {
					g.menu = "n"
				}
				g.keys = !g.keys
			}

			if ev.Sym == sdl.K_RETURN {
				g.keys = false
				g.menu = "e"
			}

			if ev.Sym == sdl.K_ESCAPE && g.keys {
				stopMusic()
				menu.Reset()
				state = menu.Run
				return true
			}
		}
	}

	return false
}

func (g *Game) draw() bool {
	g.drawSide()

	if g.mut {
		g.drawl()
		g.check()
		for i := range g.wins {
			g.wins[i] = 0
		}
	}

	if g.credit == 0 && g.bet == 0 {
		blitText(g.creditFont, 70, 190, sdlcolor.Red, "Game Over")
	}

	g.rlayer.Blit(37, 48)
	g.windowLayer.Blit(0, 0)

	if !g.keys {
		switch g.menu {
		case "h":
			g.helpMenu()
		case "e":
			if g.endGame() {
				return true
			}
		}
	}

	screen.Present()
	return false
}

func (g *Game) drawSide() {
	// animation
	blitText(g.digiFont, 470, 50, sdl.Color{60, 0, 0, 255}, "88888888888")

	blitText(g.digiFont, 470, 50, sdlcolor.White, "F1 FOR HELP")

	blitText(g.font, 500, 185, sdl.Color{230, 255, 255, 255}, "Bet:")

	// multip
	blitText(g.digiFont, 500, 210, sdl.Color{60, 0, 0, 255}, "88")

	blitText(g.digiFont, 500, 210, sdl.Color{255, 0, 0, 255}, fmt.Sprintf("%02d", g.bet))

	blitText(g.font, 500, 255, sdl.Color{230, 255, 255, 255}, "Winner Paid:")

	// last win
	blitText(g.digiFont, 500, 280, sdl.Color{60, 0, 0, 255}, "888")

	blitText(g.digiFont, 500, 280, sdl.Color{255, 0, 0, 255}, fmt.Sprintf("%03d", g.lastwin))

	blitText(g.font, 500, 325, sdl.Color{230, 255, 255, 255}, "Credit:")

	// startsum
	blitText(g.digiFont, 500, 350, sdl.Color{60, 0, 0, 255}, "888888")

	blitText(g.digiFont, 500, 350, sdl.Color{255, 0, 0, 255}, fmt.Sprintf("%06d", g.credit))

}

func (g *Game) drawl() {
	xs := [3]int{36, 165, 295}
	ys := [3]int{46, 174, 302}

	var i int
	for _, x := range xs {
		for _, y := range ys {
			g.images[g.show[i]-1].Blit(x, y)
			i++
		}
	}
}

func (g *Game) randi() {
	copy(g.showOld[:], g.show[:])
	g.mut = true

	for i := range g.show {
		r := randn(1, 335)
		n := 0
		switch {
		case 1 <= r && r <= 5:
			n = 8
		case 6 <= r && r <= 15:
			n = 7
		case 16 <= r && r <= 30:
			n = 6
		case 31 <= r && r <= 50:
			n = 5
		case 51 <= r && r <= 120:
			n = 4
		case 121 <= r && r <= 180:
			n = 3
		case 181 <= r && r <= 253:
			n = 2
		case 254 <= r && r <= 334:
			n = 1
		}

		g.show[i] = n
	}
}

func (g *Game) check() {
	s := g.show
	if s[0] == s[3] && s[3] == s[6] {
		sdlgfx.ThickLine(screen.Renderer, 36, 111, 423, 111, 8, sdl.Color{246, 226, 0, 255})
		g.wins[0] = s[0]
	}
	if s[1] == s[4] && s[4] == s[7] {
		sdlgfx.ThickLine(screen.Renderer, 36, 239, 423, 239, 8, sdl.Color{246, 226, 0, 255})
		g.wins[1] = s[1]
	}
	if s[2] == s[5] && s[5] == s[8] {
		sdlgfx.ThickLine(screen.Renderer, 36, 367, 423, 367, 8, sdl.Color{246, 226, 0, 255})
		g.wins[2] = s[2]
	}
	if s[0] == s[4] && s[4] == s[8] {
		sdlgfx.ThickLine(screen.Renderer, 37, 47, 422, 433, 8, sdl.Color{246, 226, 0, 255})
		g.wins[3] = s[0]
	}
	if s[2] == s[4] && s[4] == s[6] {
		sdlgfx.ThickLine(screen.Renderer, 37, 432, 422, 47, 8, sdl.Color{246, 226, 0, 255})
		g.wins[4] = s[2]
	}
}

func (g *Game) genRollColumn(col, n int) []*Image {
	var m []*Image

	img := g.images
	s := g.show
	m = append(m, img[s[n]-1])
	m = append(m, img[s[n+1]-1])
	m = append(m, img[s[n+2]-1])
	for i := 0; i <= col-3; i++ {
		m = append(m, img[randn(0, 8)])
	}

	s = g.showOld
	m = append(m, img[s[n]-1])
	m = append(m, img[s[n+1]-1])
	m = append(m, img[s[n+2]-1])

	return m
}

func (g *Game) rollColumn(r []*Image, l, x int) ([]*Image, int) {
	if l > 2 {
		r[len(r)-3].Blit(x, 46)
		r[len(r)-2].Blit(x, 174)
		r[len(r)-1].Blit(x, 302)
		l--
		r = r[:len(r)-1]
	} else {
		r[len(r)-3].Blit(x, 46)
		r[len(r)-2].Blit(x, 174)
		r[len(r)-1].Blit(x, 302)
	}

	return r, l
}

func (g *Game) roll() {
	// toll time
	a := randn(5, 14)
	b := randn(a+1, a+5)
	c := randn(b+1, b+5)

	ra := g.genRollColumn(a, 0)
	ca := playSound(g.rsound)

	rb := g.genRollColumn(b, 3)
	cb := playSound(g.rsound)

	rc := g.genRollColumn(c, 6)
	cc := playSound(g.rsound)

	la := len(ra) - 1
	lb := len(rb) - 1
	lc := len(rc) - 1

	for lc > 2 {
		g.qevent(ca, cb, cc)

		screen.SetDrawColor(sdlcolor.Black)
		g.background.Blit(0, 0)

		ra, la = g.rollColumn(ra, la, 36)
		rb, lb = g.rollColumn(rb, lb, 165)
		rc, lc = g.rollColumn(rc, lc, 295)

		if la <= 2 {
			sdlmixer.HaltChannel(ca)
		}
		if lb <= 2 {
			sdlmixer.HaltChannel(cb)
		}
		if lc <= 2 {
			sdlmixer.HaltChannel(cc)
		}

		g.drawSide()
		g.rlayer.Blit(37, 48)
		g.windowLayer.Blit(0, 0)
		screen.Present()
	}
}

func (g *Game) qevent(ca, cb, cc int) {
	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}
		switch ev := ev.(type) {
		case sdl.QuitEvent:
			os.Exit(0)
		case sdl.KeyDownEvent:
			switch ev.Sym {
			case sdl.K_ESCAPE:
				sdlmixer.HaltChannel(ca)
				sdlmixer.HaltChannel(cb)
				sdlmixer.HaltChannel(cc)
				stopMusic()
				menu.Reset()
				state = menu.Run
				panic(nil)
			}
		}
	}
}

func (g *Game) winner() {
	g.lastwin = 0
	for _, n := range g.wins {
		winsu := g.bet * n
		winsum := winsu + g.bet
		if winsum > g.bet {
			g.credit += winsum
			g.lastwin += winsum
			playSound(g.beepsound)
		}
		if g.credit > maxScore {
			g.credit = maxScore
		}
	}
}

func (g *Game) helpMenu() {
	sdlgfx.ThickLine(screen.Renderer, 50, 250, 590, 250, 400, sdl.Color{176, 176, 176, 255})

	y := 250 - 120
	blitText(g.font, 60, y+60, sdlcolor.Red, "How to play:")
	blitText(g.font, 60, y+80, sdlcolor.Red, "New spin: left or right arrow")
	blitText(g.font, 60, y+100, sdlcolor.Red, "Raise bet: up arrow")
	blitText(g.font, 60, y+120, sdlcolor.Red, "To end game to high score press Enter")
	blitText(g.font, 60, y+160, sdlcolor.Red, "To close this as game over help press F1")
}

func (g *Game) endGame() bool {
	sdlgfx.ThickLine(screen.Renderer, 50, 250, 590, 250, 400, sdl.Color{176, 176, 176, 255})

	if g.credit > score {
		y := 250 - 110
		blitText(g.font, 60, y+60, sdlcolor.Red, "You have a new high score!!!")
		blitText(g.font, 60, y+80, sdlcolor.Red, fmt.Sprint("Old high score: ", score))
		blitText(g.font, 60, y+100, sdlcolor.Red, fmt.Sprint("New high score: ", g.credit))
		saveScore(score)
	} else {
		y := 180
		blitText(g.font, 100, y+60, sdlcolor.Red, "You ended the game, but you don't have a new high score...")
	}

	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}
		switch ev.(type) {
		case sdl.KeyDownEvent:
			stopMusic()
			state = menu.Run
			if g.credit > score {
				score = g.credit
			}
			return true
		}
	}

	return false
}
