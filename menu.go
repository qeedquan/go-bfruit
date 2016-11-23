package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

type bgSlide struct {
	selected int
	alpha    int
	counter  int
}

type Menu struct {
	bg       *bgSlide
	selected int
	showHS   bool

	bsound *sdlmixer.Chunk

	sav             *Image
	highScore       *Image
	background      *Image
	backgroundAdded *Image
	menuBG          []*Image

	smallFont *sdlttf.Font
	font      *sdlttf.Font

	mid       []int
	choices   []string
	allChoice string
	selector  Selector
}

type Selector interface {
	Choices() []string
	Select(choice int) bool
}

type menuSelector struct{}

func (menuSelector) Choices() []string {
	return []string{
		"  New Game  ",
		"  Settings  ",
		"  High score  ",
		"  Exit  ",
	}
}

func (menuSelector) Select(choice int) bool {
	switch choice {
	case 0:
		state = game.Run
		return true
	case 1:
		state = settings.Run
		return true
	case 2:
	default:
		os.Exit(0)
	}

	return false
}

type settingsSelector struct{}

func (settingsSelector) Choices() []string {
	return []string{
		"  Fullscreen  ",
		"  Exit  ",
	}
}

func (settingsSelector) Select(choice int) bool {
	switch choice {
	case 0:
		conf.fullscreen = !conf.fullscreen
		flags := sdl.WindowFlags(0)
		if conf.fullscreen {
			flags |= sdl.WINDOW_FULLSCREEN_DESKTOP
		}
		screen.SetFullscreen(flags)
		return false
	case 1:
		state = menu.Run
		return true
	}
	return false
}

var (
	bgSlider = &bgSlide{}
)

func newMenu(selector Selector) *Menu {
	m := &Menu{
		bg:        bgSlider,
		smallFont: loadFont("LiberationSans-Regular.ttf", 15),
		font:      loadFont("LiberationSans-Regular.ttf", 25),
		bsound:    loadSound("sounds/CLICK10A.WAV"),
		menuBG: []*Image{
			loadImage("menubg/al.png"),
			loadImage("menubg/ci.png"),
			loadImage("menubg/he.png"),
			loadImage("menubg/na.png"),
			loadImage("menubg/di.png"),
		},
		sav:             loadImage("menubg/sav.png"),
		highScore:       loadImage("menubg/highscore.png"),
		background:      loadImage("menubg/menubg.png"),
		backgroundAdded: loadImage("menubg/added.png"),
		selector:        selector,
		choices:         selector.Choices(),
	}

	for _, s := range m.choices {
		w, _, err := m.font.SizeUTF8(s)
		ck(err)
		m.mid = append(m.mid, w)
	}
	m.allChoice = strings.Join(m.choices, "")

	return m
}

func (m *Menu) Reset() {
	*m.bg = bgSlide{}
}

func (m *Menu) Run() {
	for {
		if m.event() {
			return
		}
		m.draw()
		sdl.Delay(1000 / 60)
	}
}

func (m *Menu) event() bool {
	for {
		m.showHS = false
		if m.selected == 2 {
			m.showHS = true
		}

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
				os.Exit(0)
			case sdl.K_LEFT:
				if m.selected--; m.selected < 0 {
					m.selected = len(m.choices) - 1
				}
			case sdl.K_RIGHT:
				if m.selected++; m.selected >= len(m.choices) {
					m.selected = 0
				}
			case sdl.K_SPACE, sdl.K_RETURN:
				if m.selector.Select(m.selected) {
					return true
				}
			}
		}
	}
	return false
}

func (m *Menu) draw() {
	screen.SetDrawColor(sdlcolor.Black)
	screen.Clear()

	bg := m.menuBG[m.bg.selected]
	bg.SetAlphaMod(uint8(m.bg.alpha))
	bg.Blit(0, 0)
	m.backgroundAdded.Blit(0, 0)
	m.drawSelection()
	m.background.Blit(0, 0)

	blitText(m.smallFont, 3, 460, sdlcolor.White, fmt.Sprint("Balazs Nagy - BFruit -", version))

	if m.showHS {
		m.sav.Blit(0, 60)
		m.sav.Blit(0, 120)
		m.highScore.Blit(50, 60)
		blitText(m.font, 295, 110, sdlcolor.White, fmt.Sprint(score))
	}

	m.bg.counter += 4
	if m.bg.counter < 245 {
		m.bg.alpha += 4
	}
	if m.bg.counter > 244 {
		m.bg.alpha -= 4
	}

	if m.bg.counter > 490 {
		m.bg.alpha = 0
		m.bg.counter = 0
		m.bg.selected = (m.bg.selected + 1) % len(m.menuBG)
	}

	screen.Present()
}

func (m *Menu) drawSelection() {
	x := 0
	for i := 0; i <= m.selected; i++ {
		x -= m.mid[i]
	}
	x += m.mid[m.selected] / 2
	blitText(m.font, 320+x, 15, sdlcolor.White, fmt.Sprint(m.allChoice))
}
