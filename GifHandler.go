package main

import (
	"github.com/unix-streamdeck/streamdeck"
	"image/gif"
	"log"
	"os"
	"time"
)

type GifIconHandler struct{
}

func (GifIconHandler) Icon(page int, index int, key *Key, dev streamdeck.Device) {
	f, err := os.Open(key.Icon)
	if err != nil {
		log.Fatal(err)
		return
	}
	gifs, err := gif.DecodeAll(f)
	timeDelay := gifs.Delay[0]
	gifIndex := 0
	go loop(gifs, gifIndex, timeDelay, page, index, dev, key)
}

func loop(gifs *gif.GIF, gifIndex int, timeDelay int, page int, index int, dev streamdeck.Device, key *Key) {
	for true {
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