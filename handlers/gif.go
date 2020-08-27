package handlers

import (
	"github.com/nfnt/resize"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/driver"
	"image"
	"image/gif"
	"log"
	"os"
	"time"
)

func (s *GifIconHandler) Icon(page int, index int, key *api.Key, dev streamdeck.Device) {
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
		frames[i] = resize.Resize(dev.Pixels, dev.Pixels, frame, resize.Lanczos3)
	}
	go loop(frames, timeDelay, page, index, dev, key, s)
}

func (s *GifIconHandler) Stop() {
	s.Running = false
}

func loop(frames []image.Image, timeDelay int, page int, index int, dev streamdeck.Device, key *api.Key, s *GifIconHandler) {
	gifIndex := 0
	for s.Running {
		img := frames[gifIndex]
		s.OnSetImage(img, index, page, dev)
		key.Buff = img
		gifIndex++
		if gifIndex >= len(frames) {
			gifIndex = 0
		}
		time.Sleep(time.Duration(timeDelay * 10000000))
	}
}
