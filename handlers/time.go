package handlers

import (
	"github.com/unix-streamdeck/api"
	"image"
	"image/draw"
	"log"
	"time"
)


func (t *TimeIconHandler) Start(k api.Key, info api.StreamDeckInfo, callback func(image image.Image)) {
	t.Running = true
	go timeLoop(k, info, callback, t)
}

func (t *TimeIconHandler) IsRunning() bool {
	return t.Running
}

func (t *TimeIconHandler) SetRunning(running bool)  {
	t.Running = running
}

func (t *TimeIconHandler) Stop() {
	t.Running = false
}

func timeLoop(k api.Key, info api.StreamDeckInfo, callback func(image image.Image), handler *TimeIconHandler) {
	for handler.Running {
		img := image.NewRGBA(image.Rect(0, 0, info.IconSize, info.IconSize))
		draw.Draw(img, img.Bounds(), image.Black, image.ZP, draw.Src)
		t := time.Now()
		tString := t.Format("15:04:05")
		imgParsed, err := api.DrawText(img, tString, k.TextSize, k.TextAlignment)
		if err != nil {
			log.Println(err)
		} else {
			callback(imgParsed)
		}
		time.Sleep(time.Second)
	}
}