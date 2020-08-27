package handlers

import (
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/driver"
	"github.com/nfnt/resize"
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
	timeDelay := gifs.Delay[0]
	gifIndex := 0
	go loop(gifs, gifIndex, timeDelay, page, index, dev, key, s)
}

func (s *GifIconHandler) Stop() {
	s.Running = false
}

func loop(gifs *gif.GIF, gifIndex int, timeDelay int, page int, index int, dev streamdeck.Device, key *api.Key, s *GifIconHandler) {
	for s.Running {
		img := resize.Resize(dev.Pixels, dev.Pixels, gifs.Image[gifIndex], resize.Lanczos3)
		s.OnSetImage(img, index, page, dev)
		key.Buff = img
		gifIndex++
		if gifIndex >= len(gifs.Image) {
			gifIndex = 0
		}
		time.Sleep(time.Duration(timeDelay * 10000000))
	}
}