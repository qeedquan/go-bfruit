package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func loadScore() int {
	var err error

	defer func() {
		if err != nil {
			log.SetPrefix("score: ")
			log.Println(err)
		}
	}()

	filename := filepath.Join(conf.pref, "score")
	f, err := os.Open(filename)
	if err != nil {
		return 1
	}
	defer f.Close()

	score := 1
	fmt.Fscan(f, &score)
	return score
}

func saveScore(score int) {
	var err error

	defer func() {
		if err != nil {
			log.SetPrefix("score: ")
			log.Println(err)
		}
	}()

	filename := filepath.Join(conf.pref, "score")
	f, err := os.Create(filename)
	if err != nil {
		return
	}

	_, err = fmt.Fprint(f, score)
	errClose := f.Close()

	if err == nil {
		err = errClose
	}
}
