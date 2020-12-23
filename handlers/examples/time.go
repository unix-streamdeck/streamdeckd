package examples

import (
	"github.com/unix-streamdeck/api"
	"github.com/unix-streamdeck/streamdeckd/handlers"
	"image"
	"image/draw"
	"log"
	"time"
)

type TimeIconHandler struct {
	Running bool
}

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

func RegisterTime() handlers.Module {
	return handlers.Module{NewIcon: func() api.IconHandler {
		return &TimeIconHandler{Running: true}
	}, Name: "Time"}
}