package handlers

import (
	"github.com/fogleman/gg"
	"github.com/unix-streamdeck/api"
	"golang.org/x/image/font/inconsolata"
	"image"
	"strconv"
	"time"
)

func (c *CounterIconHandler) Icon(key *api.Key, _ api.StreamDeckInfo, callback func(image image.Image)) {
	if c.Callback == nil {
		c.Callback = callback
	}
	if c.Running {
		img := gg.NewContext(72, 72)
		img.SetRGB(0, 0, 0)
		img.Clear()
		img.SetRGB(1, 1, 1)
		img.SetFontFace(inconsolata.Regular8x16)
		Count := strconv.Itoa(c.Count)
		img.DrawStringAnchored(Count, 72/2, 72/2, 0.5, 0.5)
		img.Clip()
		callback(img.Image())
		key.Buff = img.Image()
		time.Sleep(250 * time.Millisecond)
	}
}

func (c CounterIconHandler) Stop()  {
	c.Running = false
}

type CounterKeyHandler struct{}

func (CounterKeyHandler) Key(key *api.Key, info api.StreamDeckInfo) {
	if key.IconHandler != "Counter" {
		return
	}
	handler := key.IconHandlerStruct.(*CounterIconHandler)
	handler.Count += 1
	if handler.Callback != nil {
		handler.Icon(key, info, handler.Callback)
	}
}
