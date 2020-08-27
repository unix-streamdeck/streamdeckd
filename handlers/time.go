package handlers

import (
	"github.com/fogleman/gg"
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/driver"
	"golang.org/x/image/font/inconsolata"
	"strconv"
	"time"
)


func (t *TimeIconHandler) Icon(page int, index int, key *api.Key, dev streamdeck.Device) {
	t.Running = true
	go timeLoop(page, index, dev, key, t)
}

func (t *TimeIconHandler) Stop() {
	t.Running = false
}

func timeLoop(page int, index int, dev streamdeck.Device, key *api.Key, handler *TimeIconHandler) {
	for handler.Running {
		img := gg.NewContext(72, 72)
		img.SetRGB(0, 0, 0)
		img.Clear()
		img.SetRGB(1, 1, 1)
		img.SetFontFace(inconsolata.Regular8x16)
		t := time.Now()
		tString := t.Format("15:04:05")
		img.DrawStringAnchored(tString, 72/2, 72/2, 0.5, 0.5)
		img.Clip()
		handler.OnSetImage(img.Image(), index, page, dev)
		key.Buff = img.Image()
		time.Sleep(time.Second)
	}
}

func zeroes(i int) (string) {
	if i < 10 {
		return "0" + strconv.Itoa(i)
	} else {
		return strconv.Itoa(i)
	}
}