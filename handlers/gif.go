package handlers

import (
	"context"
	"github.com/unix-streamdeck/api"
	"image"
	"image/gif"
	"log"
	"os"
	"time"
)

func (s *GifIconHandler) Start(key api.Key, info api.StreamDeckInfo, callback func(image image.Image)) {
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
		img := api.ResizeImage(frame, info.IconSize)
		if key.Text != "" {
			img, err = api.DrawText(img, key.Text)
			if err != nil {
				log.Println(err)
			}
		}
		frames[i] = img
	}
	go loop(frames, timeDelay, callback, s)
}

func (s *GifIconHandler) IsRunning() bool {
	return s.Running
}

func (s *GifIconHandler) SetRunning(running bool)  {
	s.Running = running
}

func (s *GifIconHandler) Stop() {
	s.Running = false
}

func loop(frames []image.Image, timeDelay int, callback func(image image.Image), s *GifIconHandler) {
	ctx := context.Background()
	err := s.Lock.Acquire(ctx, 1)
	if err != nil {
		return
	}
	defer s.Lock.Release(1)
	gifIndex := 0
	for s.Running {
		img := frames[gifIndex]
		callback(img)
		gifIndex++
		if gifIndex >= len(frames) {
			gifIndex = 0
		}
		time.Sleep(time.Duration(timeDelay * 10000000))
	}
}
