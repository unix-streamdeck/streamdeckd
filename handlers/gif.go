package handlers

import (
	"github.com/nfnt/resize"
	"github.com/unix-streamdeck/api"
	"image"
	"image/gif"
	"log"
	"os"
	"time"
)

func (s *GifIconHandler) Icon(key *api.Key, info api.StreamDeckInfo, callback func(image image.Image)) {
	s.Running = true
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
		frames[i] = resize.Resize(uint(info.IconSize), uint(info.IconSize), frame, resize.Lanczos3)
	}
	go loop(frames, timeDelay, callback, key, s)
}

func (s *GifIconHandler) Stop() {
	s.Running = false
}

func loop(frames []image.Image, timeDelay int, callback func(image image.Image), key *api.Key, s *GifIconHandler) {
	gifIndex := 0
	for s.Running {
		img := frames[gifIndex]
		callback(img)
		key.Buff = img
		gifIndex++
		if gifIndex >= len(frames) {
			gifIndex = 0
		}
		time.Sleep(time.Duration(timeDelay * 10000000))
	}
}
