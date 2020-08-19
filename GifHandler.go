package main

import (
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/driver"
	"image/gif"
	"log"
	"os"
	"time"
)

type GifIconHandler struct {
	running bool
}

func (s *GifIconHandler) Icon(page int, index int, key *api.Key, dev streamdeck.Device) {
	s.running = true
	f, err := os.Open(key.Icon)
	if err != nil {
		log.Println(err)
		return
	}
	gifs, err := gif.DecodeAll(f)
	timeDelay := gifs.Delay[0]
	gifIndex := 0
	go loop(gifs, gifIndex, timeDelay, page, index, dev, key, s)
}

func (s *GifIconHandler) Stop() {
	s.running = false
}

func loop(gifs *gif.GIF, gifIndex int, timeDelay int, page int, index int, dev streamdeck.Device, key *api.Key, s *GifIconHandler) {
	for s.running {
		img := ResizeImage(gifs.Image[gifIndex])
		SetImage(img, index, page, dev)
		key.Buff = img
		gifIndex++
		if gifIndex >= len(gifs.Image) {
			gifIndex = 0
		}
		time.Sleep(time.Duration(timeDelay * 10000000))
	}
}