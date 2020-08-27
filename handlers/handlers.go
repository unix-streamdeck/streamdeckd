package handlers

import (
	"image"
	"github.com/unix-streamdeck/driver"
)

type CounterIconHandler struct {
	Count int
	Running bool
	OnSetImage func(img image.Image, i int, page int, dev streamdeck.Device)
}

type GifIconHandler struct {
	Running bool
	OnSetImage func(img image.Image, i int, page int, dev streamdeck.Device)
}

type TimeIconHandler struct{
	Running bool
	OnSetImage func(img image.Image, i int, page int, dev streamdeck.Device)
}

