package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/qeedquan/go-media/sdl"
	"github.com/qeedquan/go-media/sdl/sdlgfx"
	"github.com/qeedquan/go-media/sdl/sdlimage/sdlcolor"
	"github.com/qeedquan/go-media/sdl/sdlmixer"
	"github.com/qeedquan/go-media/sdl/sdlttf"
)

const version = "0.1.2"

type Display struct {
	*sdl.Window
	*sdl.Renderer
}

var (
	conf struct {
		assets     string
		pref       string
		fullscreen bool
		music      bool
		sound      bool
		invincible bool
	}

	screen *Display

	menu     *Menu
	settings *Menu
	game     *Game
	state    func()
	score    int
	fps      sdlgfx.FPSManager
	texture  *sdl.Texture
	surface  *sdl.Surface
)

func main() {
	runtime.LockOSThread()
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)
	parseFlags()
	initSDL()
	load()
	loop()
}

func parseFlags() {
	conf.assets = filepath.Join(sdl.GetBasePath(), "assets")
	conf.pref = sdl.GetPrefPath("", "bfruit")
	flag.StringVar(&conf.assets, "assets", conf.assets, "assets directory")
	flag.StringVar(&conf.pref, "pref", conf.pref, "pref directory")
	flag.BoolVar(&conf.fullscreen, "fullscreen", false, "fullscreen")
	flag.BoolVar(&conf.music, "music", true, "enable music")
	flag.BoolVar(&conf.sound, "sound", true, "enable sound")
	flag.BoolVar(&conf.invincible, "invincible", false, "don't lose credit")
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "BFruit %v: [options]\n", version)
	flag.PrintDefaults()
	os.Exit(2)
}

func newDisplay(w, h int, wflag sdl.WindowFlags) (*Display, error) {
	window, renderer, err := sdl.CreateWindowAndRenderer(w, h, wflag)
	if err != nil {
		return nil, err
	}
	return &Display{window, renderer}, nil
}

func initSDL() {
	err := sdl.Init(sdl.INIT_EVERYTHING &^ sdl.INIT_AUDIO)
	ck(err)

	err = sdl.InitSubSystem(sdl.INIT_AUDIO)
	ek(err)

	err = sdlmixer.OpenAudio(44100, sdl.AUDIO_S16, 2, 8192)
	ek(err)

	err = sdlttf.Init()
	ck(err)

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "best")

	w, h := 640, 480
	wflag := sdl.WINDOW_RESIZABLE
	if conf.fullscreen {
		wflag |= sdl.WINDOW_FULLSCREEN_DESKTOP
	}
	screen, err = newDisplay(w, h, wflag)
	ck(err)

	texture, err = screen.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, w, h)
	ck(err)

	surface, err = sdl.CreateRGBSurfaceWithFormat(sdl.SWSURFACE, w, h, 32, sdl.PIXELFORMAT_ABGR8888)
	ck(err)

	screen.SetTitle("BFruit")
	screen.SetLogicalSize(640, 480)
	screen.SetDrawColor(sdlcolor.Black)
	screen.Clear()
	screen.Present()
	sdl.ShowCursor(0)

	fps.Init()
	fps.SetRate(60)
}

func load() {
	score = loadScore()
	menu = newMenu(menuSelector{})
	settings = newMenu(settingsSelector{})
	game = newGame()
}

func loop() {
	state = intro
	for {
		state()
	}
}

func menuEvent() bool {
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
				os.Exit(0)
			case sdl.K_SPACE, sdl.K_RETURN:
				state = menu.Run
				return true
			}
		}
	}
	return false
}

func intro() {
	border := loadImage("intro/border.png")
	point := loadImage("intro/point.png")

	for szam := 0; szam < 256; szam += 4 {
		if menuEvent() {
			return
		}

		screen.SetDrawColor(sdlcolor.Black)
		screen.Clear()

		border.Blit(160, 120)
		for y := 150; y <= 270; y += 30 {
			for x := 185; x <= 425; x += 30 {
				point.Blit(x, y)
			}
		}

		screen.Present()
		fps.Delay()
	}

	state = title
}

func title() {
	border := loadImage("intro/border.png")
	point := loadImage("intro/point.png")
	sun := loadImage("intro/sun.png")
	font := loadFont("LiberationSans-Regular.ttf", 25)

	start := time.Now()
	for {
		if menuEvent() {
			return
		}
		dt := time.Now().Sub(start)

		screen.SetDrawColor(sdlcolor.Black)
		screen.Clear()

		border.Blit(160, 120)
		for y := 150; y <= 240; y += 30 {
			point.Blit(185, y)
		}

		if dt < 1000*time.Millisecond {
			point.Blit(185, 270)
		}

		if dt < 1200*time.Millisecond {
			point.Blit(215, 150)
		}

		point.Blit(215, 180)

		if dt < 1150*time.Millisecond {
			point.Blit(215, 210)
		}

		if dt < 1200*time.Millisecond {
			point.Blit(215, 240)
		}

		if dt < 1230*time.Millisecond {
			point.Blit(215, 270)
		}

		if dt < 1180*time.Millisecond {
			point.Blit(245, 150)
		}

		if dt < 1200*time.Millisecond {
			point.Blit(245, 180)
		}

		point.Blit(245, 210)

		if dt < 1430*time.Millisecond {
			point.Blit(245, 240)
		}

		if dt < 1500*time.Millisecond {
			point.Blit(245, 270)
		}

		for y := 150; y <= 240; y += 30 {
			point.Blit(275, y)
		}

		if dt < 1220*time.Millisecond {
			point.Blit(275, 270)
		}

		if dt < 1410*time.Millisecond {
			point.Blit(305, 150)
		}

		if dt < 1340*time.Millisecond {
			point.Blit(305, 180)
		}

		if dt < 1400*time.Millisecond {
			point.Blit(305, 210)
		}

		if dt < 1360*time.Millisecond {
			point.Blit(305, 240)
		}

		if dt < 1100*time.Millisecond {
			point.Blit(305, 270)
		}

		for y := 150; y <= 270; y += 30 {
			point.Blit(335, y)
		}

		point.Blit(365, 150)

		if dt < 1130*time.Millisecond {
			point.Blit(365, 180)
		}

		point.Blit(365, 210)

		if dt < 1400*time.Millisecond {
			point.Blit(365, 240)
		}

		point.Blit(365, 270)
		point.Blit(395, 150)

		if dt < 1310*time.Millisecond {
			point.Blit(395, 180)
		}

		point.Blit(395, 210)

		if dt < 1250*time.Millisecond {
			point.Blit(395, 240)
		}

		point.Blit(395, 270)

		if dt < 1430*time.Millisecond {
			point.Blit(425, 150)
		}

		point.Blit(425, 180)

		if dt < 2000*time.Millisecond {
			point.Blit(425, 210)
		}

		point.Blit(425, 240)

		if dt < 2400*time.Millisecond {
			point.Blit(425, 270)
		}

		if dt > 3000*time.Millisecond {
			blitText(font, 190, 273, sdlcolor.White, "nXBalazs")
		}

		if dt > 3500*time.Millisecond {
			blitText(font, 280, 310, sdlcolor.White, "games")
		}

		if 4000*time.Millisecond < dt && dt < 5000*time.Millisecond {
			for i := uint8(0); i < 100; i++ {
				menuEvent()
				sun.SetAlphaMod(i)
				sun.Blit(0, 0)
				screen.Present()
				sdl.Delay(1000 / 60)
			}
		}

		if dt > 5000*time.Millisecond {
			break
		}

		screen.Present()
		fps.Delay()
	}

	state = menu.Run
}
