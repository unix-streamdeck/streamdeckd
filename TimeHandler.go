package main

import (
	"github.com/fogleman/gg"
	"github.com/unix-streamdeck/streamdeck"
	"golang.org/x/image/font/inconsolata"
	"strconv"
	"time"
)

type TimeIconHandler struct{
}

func (TimeIconHandler) Icon(page int, index int, key *Key, dev streamdeck.Device) {

	go timeLoop(page, index, dev, key)
}

func timeLoop(page int, index int, dev streamdeck.Device, key *Key) {
	for true {
		img := gg.NewContext(72, 72)
		img.SetRGB(0, 0, 0)
		img.Clear()
		img.SetRGB(1, 1, 1)
		img.SetFontFace(inconsolata.Regular8x16)
		t := time.Now()
		tString := zeroes(t.Hour()) + ":" + zeroes(t.Minute()) + ":" + zeroes(t.Second())
		img.DrawStringAnchored(tString, 72/2, 72/2, 0.5, 0.5)
		img.Clip()
		SetImage(img.Image(), index, page, dev)
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