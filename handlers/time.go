package handlers

import (
	"github.com/fogleman/gg"
	"github.com/unix-streamdeck/api"
	"golang.org/x/image/font/inconsolata"
	"image"
	"time"
)


func (t *TimeIconHandler) Icon(key *api.Key, info api.StreamDeckInfo, callback func(image image.Image)) {
	t.Running = true
	go timeLoop(key, info, callback, t)
}

func (t *TimeIconHandler) Stop() {
	t.Running = false
}

func timeLoop(key *api.Key, info api.StreamDeckInfo, callback func(image image.Image), handler *TimeIconHandler) {
	for handler.Running {
		img := gg.NewContext(info.IconSize, info.IconSize)
		img.SetRGB(0, 0, 0)
		img.Clear()
		img.SetRGB(1, 1, 1)
		img.SetFontFace(inconsolata.Regular8x16)
		t := time.Now()
		tString := t.Format("15:04:05")
		img.DrawStringAnchored(tString, float64(info.IconSize)/2, float64(info.IconSize)/2, 0.5, 0.5)
		img.Clip()
		callback(img.Image())
		key.Buff = img.Image()
		time.Sleep(time.Second)
	}
}