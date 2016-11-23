package main

import (
	"log"
	"path/filepath"

	"github.com/qeedquan/go-media/sdl/sdlmixer"
)

var (
	musics = make(map[string]*sdlmixer.Music)
	sounds = make(map[string]*sdlmixer.Chunk)
)

func loadMusic(name string) *sdlmixer.Music {
	log.SetPrefix("sound: ")
	filename := filepath.Join(conf.assets, name)
	if music, found := musics[filename]; found {
		return music
	}

	music, err := sdlmixer.LoadMUS(filename)
	if err != nil {
		log.Println(err)
		return nil
	}

	musics[filename] = music
	return music
}

func loadSound(name string) *sdlmixer.Chunk {
	log.SetPrefix("sound: ")
	filename := filepath.Join(conf.assets, name)
	if chunk, found := sounds[filename]; found {
		return chunk
	}

	chunk, err := sdlmixer.LoadWAV(filename)
	if err != nil {
		log.Print(err)
		return nil
	}

	sounds[filename] = chunk
	return chunk
}

func playSound(chunk *sdlmixer.Chunk) int {
	if !conf.sound || chunk == nil {
		return 0
	}
	return chunk.PlayChannel(-1, 0)
}

func playMusic(mus *sdlmixer.Music) {
	if !conf.music || mus == nil {
		return
	}

	mus.Play(-1)
}

func stopMusic() {
	sdlmixer.HaltMusic()
}
