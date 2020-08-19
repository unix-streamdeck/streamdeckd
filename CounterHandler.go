package main

import (
	"github.com/fogleman/gg"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/driver"
	"golang.org/x/image/font/inconsolata"
	"strconv"
)

type CounterIconHandler struct {
	count int
	running bool
}

func (c *CounterIconHandler) Icon(page int, index int, key *api.Key, dev streamdeck.Device) {
	if c.running {
		img := gg.NewContext(72, 72)
		img.SetRGB(0, 0, 0)
		img.Clear()
		img.SetRGB(1, 1, 1)
		img.SetFontFace(inconsolata.Regular8x16)
		count := strconv.Itoa(c.count)
		img.DrawStringAnchored(count, 72/2, 72/2, 0.5, 0.5)
		img.Clip()
		SetImage(img.Image(), index, page, dev)
		key.Buff = img.Image()
	}
}

func (c CounterIconHandler) Stop()  {
	c.running = false
}

type CounterKeyHandler struct{}

func (CounterKeyHandler) Key(page int, index int, key *api.Key, dev streamdeck.Device) {
	if key.IconHandler != "Counter" {
		return
	}
	handler := key.IconHandlerStruct.(*CounterIconHandler)
	handler.count += 1
	handler.Icon(page, index, key, dev)
}
