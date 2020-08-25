package main

import (
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/driver"
	"image"
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
	if err != nil {
		log.Println(err)
		return
	}
	timeDelay := gifs.Delay[0]
	frames := make([]image.Image, len(gifs.Image))
	for i, frame := range gifs.Image {
		frames[i] = ResizeImage(frame)
	}
	go loop(frames, timeDelay, page, index, dev, key, s)
}

func (s *GifIconHandler) Stop() {
	s.running = false
}

func loop(frames []image.Image, timeDelay int, page int, index int, dev streamdeck.Device, key *api.Key, s *GifIconHandler) {
	gifIndex := 0
	for s.running {
		img := frames[gifIndex]
		SetImage(img, index, page, dev)
		key.Buff = img
		gifIndex++
		if gifIndex >= len(frames) {
			gifIndex = 0
		}
		time.Sleep(time.Duration(timeDelay * 10000000))
	}
}
